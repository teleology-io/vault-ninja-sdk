// Package zk is the client-side half of Vault Ninja's end-to-end encryption:
// everything vn needs to turn /api/sdk/v2 ciphertext into plaintext locally.
// Byte-for-byte compatible with the server's internal/auth (seal.go,
// encrypt.go, frame.go) and the web client's zk-crypto package.
package zk

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	ephemeralPubLen = 32
	nonceLen        = 12
	tagLen          = 16
	sealInfo        = "vaultninja:seal"
)

// UnsealDataKey opens a secret's wrapped data key (x25519-sealed-v1:
// base64(ephemeral pub32 | nonce12 | AES-256-GCM ct)) with the API key's raw
// X25519 private key.
func UnsealDataKey(priv []byte, wrappedB64 string) ([]byte, error) {
	return Unseal(priv, wrappedB64)
}

// OpenKeyringGrant opens a keyring grant (the keyring's raw 32-byte private
// key sealed to us) and checks it really is the private half of the
// keyring's public key, so a key substituted by the server is refused.
func OpenKeyringGrant(priv []byte, grantB64, keyringPubB64 string) ([]byte, error) {
	kpriv, err := Unseal(priv, grantB64)
	if err != nil {
		return nil, err
	}
	sk, err := ecdh.X25519().NewPrivateKey(kpriv)
	if err != nil {
		return nil, fmt.Errorf("invalid keyring key: %w", err)
	}
	want, err := base64.StdEncoding.DecodeString(keyringPubB64)
	if err != nil || !bytes.Equal(sk.PublicKey().Bytes(), want) {
		return nil, errors.New("keyring key doesn't match its public key — refusing it")
	}
	return kpriv, nil
}

// Unseal opens any x25519-sealed-v1 blob with a raw X25519 private key.
func Unseal(priv []byte, wrappedB64 string) ([]byte, error) {
	sk, err := ecdh.X25519().NewPrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	raw, err := base64.StdEncoding.DecodeString(wrappedB64)
	if err != nil {
		return nil, fmt.Errorf("invalid wrapped key: %w", err)
	}
	if len(raw) < ephemeralPubLen+nonceLen+tagLen {
		return nil, errors.New("wrapped key too short")
	}
	epk, err := ecdh.X25519().NewPublicKey(raw[:ephemeralPubLen])
	if err != nil {
		return nil, fmt.Errorf("invalid ephemeral key: %w", err)
	}
	shared, err := sk.ECDH(epk)
	if err != nil {
		return nil, fmt.Errorf("ecdh: %w", err)
	}
	key, err := hkdf.Key(sha256.New, shared, nil, sealInfo, 32)
	if err != nil {
		return nil, fmt.Errorf("hkdf: %w", err)
	}
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	out, err := gcm.Open(nil, raw[ephemeralPubLen:ephemeralPubLen+nonceLen], raw[ephemeralPubLen+nonceLen:], nil)
	if err != nil {
		return nil, errors.New("couldn't unseal key (wrong API key?)")
	}
	return out, nil
}

// DecryptString decrypts a text value (title, label, field value, file name):
// base64(nonce12 | AES-256-GCM ct). Empty input is an empty value.
func DecryptString(key []byte, encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}
	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}
	if len(data) < nonceLen+tagLen {
		return "", errors.New("ciphertext too short")
	}
	pt, err := gcm.Open(nil, data[:nonceLen], data[nonceLen:], nil)
	if err != nil {
		return "", errors.New("decrypt failed")
	}
	return string(pt), nil
}

// DecryptFrameStream decrypts the "VNF1" framed file format:
//
//	header : "VNF1"(4) | version u8 | frameSize u32-BE
//	frame* : ctLen u32-BE | nonce(12) | ciphertext+tag(ctLen)
//
// Each frame's AAD is (index u32-BE | final u8); one-frame lookahead
// authenticates the real last frame, so truncation fails loudly.
func DecryptFrameStream(w io.Writer, r io.Reader, key []byte) error {
	gcm, err := newGCM(key)
	if err != nil {
		return err
	}
	br := bufio.NewReaderSize(r, 64*1024)

	var hdr [9]byte
	if _, err := io.ReadFull(br, hdr[:]); err != nil {
		return errors.New("file: truncated header")
	}
	if string(hdr[:4]) != "VNF1" {
		return errors.New("file: bad magic")
	}
	if hdr[4] != 1 {
		return fmt.Errorf("file: unsupported version %d", hdr[4])
	}

	readFrame := func() (nonce, ct []byte, err error) {
		var lb [4]byte
		if _, err = io.ReadFull(br, lb[:]); err != nil {
			return nil, nil, err
		}
		n := binary.BigEndian.Uint32(lb[:])
		if n < tagLen {
			return nil, nil, errors.New("frame shorter than GCM tag")
		}
		nonce = make([]byte, nonceLen)
		if _, err = io.ReadFull(br, nonce); err != nil {
			return nil, nil, err
		}
		ct = make([]byte, n)
		if _, err = io.ReadFull(br, ct); err != nil {
			return nil, nil, err
		}
		return nonce, ct, nil
	}
	open := func(index uint32, final bool, nonce, ct []byte) ([]byte, error) {
		aad := make([]byte, 5)
		binary.BigEndian.PutUint32(aad, index)
		if final {
			aad[4] = 1
		}
		return gcm.Open(nil, nonce, ct, aad)
	}

	var index uint32
	var pNonce, pCt []byte
	have := false
	for {
		nonce, ct, err := readFrame()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("file: frame %d: %w", index, err)
		}
		if have {
			pt, err := open(index, false, pNonce, pCt)
			if err != nil {
				return fmt.Errorf("file: decrypt frame %d failed", index)
			}
			if _, err := w.Write(pt); err != nil {
				return err
			}
			index++
		}
		pNonce, pCt, have = nonce, ct, true
	}
	if !have {
		return errors.New("file: no frames (truncated)")
	}
	pt, err := open(index, true, pNonce, pCt)
	if err != nil {
		return fmt.Errorf("file: decrypt final frame %d failed", index)
	}
	_, err = w.Write(pt)
	return err
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cipher: %w", err)
	}
	return cipher.NewGCM(block)
}

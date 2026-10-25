// Package creds finds and stores the vn API key.
//
// A zero-knowledge key looks like vn_org_2_<kid>_<secret>.<privateKey>. Only
// the part before "." is ever sent to the server; the private key stays on
// this machine and unseals secrets' data keys locally.
package creds

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	keyPrefix       = "vn_org_2_"
	keyringService  = "vault-ninja"
	defaultProfile  = "default"
	envToken        = "VN_TOKEN"
	envLegacyAPIKey = "VN_API_KEY" // same variable the language SDKs read
	envProfile      = "VN_PROFILE"
)

// Key is a parsed zero-knowledge API key.
type Key struct {
	Auth string // credential sent as the Bearer token
	Priv []byte // X25519 private key; never sent
}

// Parse validates a full key. Older key formats can't be used by vn because
// they don't give the client its own private key.
func Parse(full string) (*Key, error) {
	full = strings.TrimSpace(full)
	if !strings.HasPrefix(full, keyPrefix) {
		if strings.HasPrefix(full, "vn_org_") {
			return nil, errors.New("this is an older-format API key — create a new key in Settings → API Keys (it starts with vn_org_2_)")
		}
		return nil, errors.New("not a Vault Ninja API key (expected vn_org_2_…)")
	}
	auth, priv, ok := strings.Cut(full, ".")
	if !ok || priv == "" {
		return nil, errors.New("API key is missing its private part (everything after the \".\") — copy the whole key")
	}
	raw, err := base64.RawURLEncoding.DecodeString(priv)
	if err != nil || len(raw) != 32 {
		return nil, errors.New("API key's private part is malformed")
	}
	return &Key{Auth: auth, Priv: raw}, nil
}

func profile() string {
	if p := os.Getenv(envProfile); p != "" {
		return p
	}
	return defaultProfile
}

// Resolve finds the key: VN_TOKEN, then VN_API_KEY, then tokenFile, then the
// OS keychain (vn login). source describes where it came from.
func Resolve(tokenFile string) (k *Key, source string, err error) {
	if v := os.Getenv(envToken); v != "" {
		k, err = Parse(v)
		return k, "$" + envToken, err
	}
	if v := os.Getenv(envLegacyAPIKey); v != "" {
		k, err = Parse(v)
		return k, "$" + envLegacyAPIKey, err
	}
	if tokenFile != "" {
		b, rerr := os.ReadFile(tokenFile)
		if rerr != nil {
			return nil, tokenFile, fmt.Errorf("read token file: %w", rerr)
		}
		k, err = Parse(string(b))
		return k, tokenFile, err
	}
	v, kerr := keyring.Get(keyringService, profile())
	if kerr == nil {
		k, err = Parse(v)
		return k, "keychain (" + profile() + ")", err
	}
	return nil, "", errors.New("no API key: set VN_TOKEN, pass --token-file, or run `vn login`")
}

// Save stores a key in the OS keychain after validating it.
func Save(full string) error {
	if _, err := Parse(full); err != nil {
		return err
	}
	return keyring.Set(keyringService, profile(), strings.TrimSpace(full))
}

// Delete removes the stored key, if any.
func Delete() error {
	err := keyring.Delete(keyringService, profile())
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	return err
}

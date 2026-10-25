// Package api talks to Vault Ninja's zero-knowledge SDK API (/api/sdk/v2)
// and decrypts everything locally. The server only ever sees the key's
// auth part and only ever returns ciphertext.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/teleology-io/vault-ninja-sdk/go/cmd/vn/internal/creds"
	"github.com/teleology-io/vault-ninja-sdk/go/cmd/vn/internal/zk"
)

const defaultHost = "https://api.vaultninja.org"

// Secret is a decrypted secret. Summaries (from List) have no fields/files.
type Secret struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags"`
	UpdatedAt   time.Time `json:"updated_at"`
	Fields      []Field   `json:"fields,omitempty"`
	Files       []File    `json:"files,omitempty"`

	dataKey []byte
}

type Field struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Label    string `json:"label"`
	Value    string `json:"value"`
	Position int    `json:"position"`
}

type File struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ContentType    string `json:"content_type,omitempty"`
	Size           int64  `json:"size"`
	BytesEncrypted bool   `json:"-"`
}

// Error is a non-2xx API response.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string { return fmt.Sprintf("vault ninja: %s (HTTP %d)", e.Message, e.Status) }

type keyringWire struct {
	KeyringID  string `json:"keyring_id"`
	Version    int    `json:"version"`
	WrappedKey string `json:"wrapped_key"`
	Grant      string `json:"grant"`
	PublicKey  string `json:"public_key"`
}

type Client struct {
	base string
	key  *creds.Key
	kind string // "cli" or "mcp" — shown in the org audit log
	http *http.Client

	// opened keyring private keys, by "id:version"
	keyrings map[string][]byte
}

// New builds a client. kind is "cli" or "mcp". VN_API_URL may be a host or
// any /api/sdk/... URL (the language SDKs use the v1 base); it's normalized.
func New(key *creds.Key, kind string) *Client {
	return &Client{base: baseURL(os.Getenv("VN_API_URL")), key: key, kind: kind, http: &http.Client{Timeout: 60 * time.Second}, keyrings: map[string][]byte{}}
}

var sdkPathRe = regexp.MustCompile(`/api/sdk(/v\d+)?/?$`)

func baseURL(env string) string {
	host := strings.TrimRight(env, "/")
	if host == "" {
		host = defaultHost
	}
	return sdkPathRe.ReplaceAllString(host, "") + "/api/sdk/v2"
}

func (c *Client) get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.key.Auth)
	req.Header.Set("X-VN-Client", c.kind)
	req.Header.Set("User-Agent", "vn/"+Version)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		var body struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&body)
		if body.Error == "" {
			body.Error = http.StatusText(resp.StatusCode)
		}
		return nil, &Error{Status: resp.StatusCode, Message: body.Error}
	}
	return resp, nil
}

func getJSON[T any](c *Client, ctx context.Context, path string) (T, error) {
	var out T
	resp, err := c.get(ctx, path)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	return out, json.NewDecoder(resp.Body).Decode(&out)
}

// Wire shapes (ciphertext).
type wireSummary struct {
	ID             string       `json:"id"`
	Title          string       `json:"title"`
	Tags           []string     `json:"tags"`
	UpdatedAt      time.Time    `json:"updated_at"`
	WrappedDataKey *string      `json:"wrapped_data_key"`
	KeyringWrap    *keyringWire `json:"keyring_wrap"`
}

// dataKey opens a secret's data key: keyring path first (preferred), then
// the legacy per-key wrap.
func (c *Client) dataKey(w wireSummary) ([]byte, error) {
	var kerr error
	if kw := w.KeyringWrap; kw != nil {
		id := fmt.Sprintf("%s:%d", kw.KeyringID, kw.Version)
		kpriv, ok := c.keyrings[id]
		if !ok {
			if kpriv, kerr = zk.OpenKeyringGrant(c.key.Priv, kw.Grant, kw.PublicKey); kerr == nil {
				c.keyrings[id] = kpriv
			}
		}
		if kerr == nil {
			dk, err := zk.Unseal(kpriv, kw.WrappedKey)
			if err == nil {
				return dk, nil
			}
			kerr = err
		}
	}
	if w.WrappedDataKey == nil {
		if kerr != nil {
			return nil, fmt.Errorf("keyring: %w", kerr)
		}
		return nil, errors.New("no key for this secret")
	}
	return zk.UnsealDataKey(c.key.Priv, *w.WrappedDataKey)
}

type wireSecret struct {
	wireSummary
	Description string `json:"description"`
	Fields      []struct {
		ID        string `json:"id"`
		FieldType string `json:"field_type"`
		Label     string `json:"label"`
		Value     string `json:"value"`
		Position  int    `json:"position"`
	} `json:"fields"`
	Files []struct {
		ID                 string `json:"id"`
		OriginalName       string `json:"original_name"`
		ContentType        string `json:"content_type"`
		SizeBytes          int64  `json:"size_bytes"`
		PlaintextSizeBytes *int64 `json:"plaintext_size_bytes"`
		BytesEncrypted     bool   `json:"bytes_encrypted"`
	} `json:"files"`
}

func (c *Client) openSummary(w wireSummary) (*Secret, error) {
	dk, err := c.dataKey(w)
	if err != nil {
		return nil, err
	}
	title, err := zk.DecryptString(dk, w.Title)
	if err != nil {
		return nil, fmt.Errorf("title: %w", err)
	}
	tags := make([]string, 0, len(w.Tags))
	for _, t := range w.Tags {
		if pt, err := zk.DecryptString(dk, t); err == nil {
			tags = append(tags, pt)
		}
	}
	return &Secret{ID: w.ID, Title: title, Tags: tags, UpdatedAt: w.UpdatedAt, dataKey: dk}, nil
}

// List returns every secret shared with this key, titles decrypted, sorted
// by title. Secrets that fail to decrypt are skipped.
func (c *Client) List(ctx context.Context) ([]Secret, error) {
	wire, err := getJSON[[]wireSummary](c, ctx, "/secrets")
	if err != nil {
		return nil, err
	}
	out := make([]Secret, 0, len(wire))
	for _, w := range wire {
		s, err := c.openSummary(w)
		if err != nil {
			continue
		}
		out = append(out, *s)
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Title) < strings.ToLower(out[j].Title) })
	return out, nil
}

// Get returns one fully decrypted secret by id.
func (c *Client) Get(ctx context.Context, id string) (*Secret, error) {
	w, err := getJSON[wireSecret](c, ctx, "/secrets/"+id)
	if err != nil {
		return nil, err
	}
	s, err := c.openSummary(w.wireSummary)
	if err != nil {
		return nil, err
	}
	if s.Description, err = zk.DecryptString(s.dataKey, w.Description); err != nil {
		return nil, fmt.Errorf("description: %w", err)
	}
	for _, f := range w.Fields {
		label, err1 := zk.DecryptString(s.dataKey, f.Label)
		value, err2 := zk.DecryptString(s.dataKey, f.Value)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("field %s: decrypt failed", f.ID)
		}
		s.Fields = append(s.Fields, Field{ID: f.ID, Type: f.FieldType, Label: label, Value: value, Position: f.Position})
	}
	sort.SliceStable(s.Fields, func(i, j int) bool { return s.Fields[i].Position < s.Fields[j].Position })
	for _, f := range w.Files {
		name, _ := zk.DecryptString(s.dataKey, f.OriginalName)
		ct, _ := zk.DecryptString(s.dataKey, f.ContentType)
		size := f.SizeBytes
		if f.PlaintextSizeBytes != nil {
			size = *f.PlaintextSizeBytes
		}
		s.Files = append(s.Files, File{ID: f.ID, Name: name, ContentType: ct, Size: size, BytesEncrypted: f.BytesEncrypted})
	}
	return s, nil
}

var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// Resolve finds a secret by id or title: exact (case-insensitive) title
// match first, then a unique title prefix.
func (c *Client) Resolve(ctx context.Context, ref string) (*Secret, error) {
	if uuidRe.MatchString(ref) {
		return c.Get(ctx, ref)
	}
	all, err := c.List(ctx)
	if err != nil {
		return nil, err
	}
	var exact, prefix []Secret
	for _, s := range all {
		switch {
		case strings.EqualFold(s.Title, ref):
			exact = append(exact, s)
		case strings.HasPrefix(strings.ToLower(s.Title), strings.ToLower(ref)):
			prefix = append(prefix, s)
		}
	}
	match := exact
	if len(match) == 0 {
		match = prefix
	}
	switch len(match) {
	case 0:
		return nil, fmt.Errorf("no secret matching %q is shared with this key", ref)
	case 1:
		return c.Get(ctx, match[0].ID)
	}
	names := make([]string, 0, len(match))
	for _, s := range match {
		names = append(names, fmt.Sprintf("%q (%s)", s.Title, s.ID))
	}
	return nil, fmt.Errorf("%q matches several secrets — use the id: %s", ref, strings.Join(names, ", "))
}

// Value picks a field value by selector: a field id, a label
// (case-insensitive), or a field type if only one field has it. Append
// "#login", "#password", "#number", "#expiry" or "#pin" to pick part of a
// login or card field. An empty selector works when the secret has exactly
// one field.
func (s *Secret) Value(selector string) (string, error) {
	sel, part, _ := strings.Cut(selector, "#")
	var f *Field
	switch {
	case sel == "" && len(s.Fields) == 1:
		f = &s.Fields[0]
	case sel == "":
		return "", fmt.Errorf("secret %q has %d fields — say which one", s.Title, len(s.Fields))
	default:
		var byType []*Field
		for i := range s.Fields {
			fi := &s.Fields[i]
			if fi.ID == sel || strings.EqualFold(fi.Label, sel) {
				f = fi
				break
			}
			if strings.EqualFold(fi.Type, sel) {
				byType = append(byType, fi)
			}
		}
		if f == nil && len(byType) == 1 {
			f = byType[0]
		}
	}
	if f == nil {
		return "", fmt.Errorf("secret %q has no field %q", s.Title, sel)
	}
	if part == "" {
		return f.Value, nil
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(f.Value), &obj); err != nil {
		return "", fmt.Errorf("field %q has no parts (it isn't a login or card field)", f.Label)
	}
	v, ok := obj[part]
	if !ok {
		return "", fmt.Errorf("field %q has no part %q", f.Label, part)
	}
	return fmt.Sprint(v), nil
}

// DownloadFile writes a file's plaintext to w.
func (c *Client) DownloadFile(ctx context.Context, s *Secret, fileID string, w io.Writer) error {
	var f *File
	for i := range s.Files {
		if s.Files[i].ID == fileID {
			f = &s.Files[i]
		}
	}
	if f == nil {
		return errors.New("file not found on this secret")
	}
	resp, err := c.get(ctx, "/secrets/"+s.ID+"/files/"+fileID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Header.Get("X-VN-Bytes-Encrypted") == "true" {
		return zk.DecryptFrameStream(w, resp.Body, s.dataKey)
	}
	_, err = io.Copy(w, resp.Body)
	return err
}

// FindFile picks a file by id or name (case-insensitive).
func (s *Secret) FindFile(ref string) (*File, error) {
	for i := range s.Files {
		if s.Files[i].ID == ref || strings.EqualFold(s.Files[i].Name, ref) {
			return &s.Files[i], nil
		}
	}
	return nil, fmt.Errorf("secret %q has no file %q", s.Title, ref)
}

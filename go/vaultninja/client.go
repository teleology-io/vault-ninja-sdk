// Package vaultninja provides a slim read-only client for the Vault Ninja API.
// stdlib only — no external dependencies.
//
// Usage:
//
//	client := vaultninja.New("vn_org_...")
//	secrets, err := client.ListSecrets(ctx)
//
// Env vars:
//
//	VN_API_KEY  — API key (used when apiKey arg is empty)
//	VN_API_URL  — override base URL
package vaultninja

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const defaultBaseURL = "https://api.vaultninja.org/api/sdk/v1"

// FieldType is the type of a secret field.
type FieldType string

const (
	FieldTypePassword FieldType = "PASSWORD"
	FieldTypeCard     FieldType = "CARD"
	FieldTypeTOTP     FieldType = "TOTP"
	FieldTypeURL      FieldType = "URL"
	FieldTypeNote     FieldType = "NOTE"
	FieldTypeText     FieldType = "TEXT"
)

// SecretField is a single decrypted field belonging to a secret.
type SecretField struct {
	ID        string    `json:"id"`
	SecretID  string    `json:"secret_id"`
	Position  int       `json:"position"`
	FieldType FieldType `json:"field_type"`
	Label     string    `json:"label"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// SecretFile is metadata for a file attached to a secret.
type SecretFile struct {
	ID           string    `json:"id"`
	SecretID     string    `json:"secret_id"`
	OriginalName string    `json:"original_name"`
	ContentType  string    `json:"content_type"`
	SizeBytes    int64     `json:"size_bytes"`
	Position     int       `json:"position"`
	CreatedAt    time.Time `json:"created_at"`
	URL          string    `json:"url"`
}

// Secret is a secret with all its decrypted fields and file metadata.
type Secret struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Tags        []string      `json:"tags"`
	OrgID       string        `json:"org_id"`
	CreatedBy   string        `json:"created_by"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Fields      []SecretField `json:"fields"`
	Files       []SecretFile  `json:"files"`
}

// APIError is returned for non-2xx HTTP responses.
type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("vault ninja API error %d: %s", e.Status, e.Message)
}

// Client is the Vault Ninja SDK client.
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// New creates a new Client.
// apiKey falls back to VN_API_KEY; an optional baseURL overrides VN_API_URL and the default.
func New(apiKey string, baseURL ...string) *Client {
	key := apiKey
	if key == "" {
		key = os.Getenv("VN_API_KEY")
	}
	url := ""
	if len(baseURL) > 0 {
		url = baseURL[0]
	}
	if url == "" {
		url = os.Getenv("VN_API_URL")
	}
	if url == "" {
		url = defaultBaseURL
	}
	return &Client{baseURL: url, apiKey: key, http: &http.Client{}}
}

func (c *Client) get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	return c.http.Do(req)
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	defer resp.Body.Close()
	var body struct {
		Error string `json:"error"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return &APIError{Status: resp.StatusCode, Message: body.Error}
}

// ListSecrets returns all secrets in the org.
func (c *Client) ListSecrets(ctx context.Context) ([]Secret, error) {
	resp, err := c.get(ctx, "/secrets")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var result []Secret
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// GetSecret returns a secret with all decrypted fields and file metadata.
func (c *Client) GetSecret(ctx context.Context, id string) (*Secret, error) {
	resp, err := c.get(ctx, "/secrets/"+id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var result Secret
	return &result, json.NewDecoder(resp.Body).Decode(&result)
}

// GetField returns a single decrypted field.
func (c *Client) GetField(ctx context.Context, id, fid string) (*SecretField, error) {
	resp, err := c.get(ctx, "/secrets/"+id+"/fields/"+fid)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var result SecretField
	return &result, json.NewDecoder(resp.Body).Decode(&result)
}

// GetFile returns the raw bytes of a file attachment.
func (c *Client) GetFile(ctx context.Context, id, fid string) ([]byte, error) {
	resp, err := c.get(ctx, "/secrets/"+id+"/files/"+fid)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

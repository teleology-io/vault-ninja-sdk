// Package vaultninja provides a client for the Vault Ninja SDK API.
// It wraps the generated openapi client with VAULT_API_URL env var support
// and a functional-options constructor.
package vaultninja

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	openapi "github.com/teleology-io/vault-ninja-sdk/go"
)

const defaultBaseURL = "https://api.vaultninja.org/api/sdk/v1"

// Client wraps the generated API client.
type Client struct {
	api     *openapi.APIClient
	baseURL string
	apiKey  string
	httpCli *http.Client
}

// Option is a functional option for Client.
type Option func(cfg *openapi.Configuration)

// WithBaseURL overrides the API base URL.
func WithBaseURL(url string) Option {
	return func(cfg *openapi.Configuration) {
		cfg.Servers = openapi.ServerConfigurations{
			{URL: url},
		}
	}
}

// New creates a new VaultNinja client.
// Base URL resolution order:
//  1. WithBaseURL option
//  2. VAULT_API_URL environment variable
//  3. https://api.vaultninja.org/api/sdk/v1
func New(apiKey string, opts ...Option) *Client {
	baseURL := os.Getenv("VAULT_API_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{{URL: baseURL}}
	cfg.DefaultHeader["Authorization"] = "Bearer " + apiKey

	for _, o := range opts {
		o(cfg)
	}

	// Read final URL after opts may have overridden it.
	if len(cfg.Servers) > 0 {
		baseURL = cfg.Servers[0].URL
	}

	return &Client{
		api:     openapi.NewAPIClient(cfg),
		baseURL: baseURL,
		apiKey:  apiKey,
		httpCli: &http.Client{},
	}
}

// ListSecrets returns all secrets in the org (no fields or files).
// Uses a direct HTTP call with lenient JSON decoding because the server
// returns full Secret objects on this endpoint, which the generated
// SecretSummary decoder rejects via DisallowUnknownFields.
func (c *Client) ListSecrets(ctx context.Context) ([]openapi.Secret, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/secrets", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var secrets []openapi.Secret
	if err := json.NewDecoder(resp.Body).Decode(&secrets); err != nil {
		return nil, err
	}
	return secrets, nil
}

// GetSecret returns a secret with all decrypted fields and file metadata.
func (c *Client) GetSecret(ctx context.Context, id string) (*openapi.Secret, error) {
	result, _, err := c.api.SecretsAPI.GetSecret(ctx, id).Execute()
	return result, err
}

// GetField returns a single decrypted field.
func (c *Client) GetField(ctx context.Context, secretID, fieldID string) (*openapi.SecretField, error) {
	result, _, err := c.api.SecretsAPI.GetField(ctx, secretID, fieldID).Execute()
	return result, err
}

// GetFile returns the downloaded file as an *os.File.
func (c *Client) GetFile(ctx context.Context, secretID, fileID string) (*os.File, error) {
	file, _, err := c.api.SecretsAPI.GetFile(ctx, secretID, fileID).Execute()
	return file, err
}

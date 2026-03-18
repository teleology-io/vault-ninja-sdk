// Package vaultninja provides a client for the Vault Ninja SDK API.
// It wraps the generated openapi client with VAULT_API_URL env var support
// and a functional-options constructor.
package vaultninja

import (
	"context"
	"os"

	openapi "github.com/teleology-io/vault-ninja-sdk/go"
)

const defaultBaseURL = "https://vaultninja.org/api/sdk/v1"

// Client wraps the generated API client.
type Client struct {
	api *openapi.APIClient
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
//  3. https://vaultninja.org/api/sdk/v1
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

	return &Client{api: openapi.NewAPIClient(cfg)}
}

// ListSecrets returns all secrets in the org (no fields or files).
func (c *Client) ListSecrets(ctx context.Context) ([]openapi.SecretSummary, error) {
	result, _, err := c.api.SecretsAPI.ListSecrets(ctx).Execute()
	return result, err
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

// GetFile returns the raw HTTP response for a file download.
// Use resp.Body to stream the file bytes.
func (c *Client) GetFile(ctx context.Context, secretID, fileID string) (*openapi.ApiResponse, error) {
	resp, err := c.api.SecretsAPI.GetFile(ctx, secretID, fileID).Execute()
	return resp, err
}

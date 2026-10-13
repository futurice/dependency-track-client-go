package dtrack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Secret is a v5 secret's metadata; Dependency-Track never returns the value.
type Secret struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   int64  `json:"created_at,omitempty"` // epoch milliseconds
	UpdatedAt   int64  `json:"updated_at,omitempty"` // epoch milliseconds
}

type CreateSecretRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value"`
}

// UpdateSecretRequest changes a secret. A nil field keeps its current value.
type UpdateSecretRequest struct {
	Description *string `json:"description,omitempty"`
	Value       *string `json:"value,omitempty"`
}

// SecretService manages v2 secrets, which x-secret-ref config fields reference by name.
type SecretService struct {
	client *Client
}

func secretPath(name string) string {
	return fmt.Sprintf("/api/v2/secrets/%s", url.PathEscape(name))
}

// Get returns the metadata of the secret with the given name.
func (ss SecretService) Get(ctx context.Context, name string) (s Secret, err error) {
	req, err := ss.client.newRequest(ctx, http.MethodGet, secretPath(name))
	if err != nil {
		return
	}

	_, err = ss.client.doRequest(req, &s)
	return
}

// Create creates a secret, or returns a 409 APIError when the name is taken.
func (ss SecretService) Create(ctx context.Context, secret CreateSecretRequest) (err error) {
	req, err := ss.client.newRequest(ctx, http.MethodPost, "/api/v2/secrets", withBody(secret))
	if err != nil {
		return
	}

	_, err = ss.client.doRequest(req, nil)
	return
}

// Update changes a secret's description or value; an unchanged update also succeeds.
func (ss SecretService) Update(ctx context.Context, name string, update UpdateSecretRequest) (err error) {
	req, err := ss.client.newRequest(ctx, http.MethodPatch, secretPath(name), withBody(update))
	if err != nil {
		return
	}

	_, err = ss.client.doRequest(req, nil)
	return
}

func (ss SecretService) Delete(ctx context.Context, name string) (err error) {
	req, err := ss.client.newRequest(ctx, http.MethodDelete, secretPath(name))
	if err != nil {
		return
	}

	_, err = ss.client.doRequest(req, nil)
	return
}

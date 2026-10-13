package dtrack

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestSecretService(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		var created []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPost, "http://localhost/api/v2/secrets",
			func(req *http.Request) (*http.Response, error) {
				require.Equal(t, "application/json, application/problem+json", req.Header.Get("Accept"))
				return captureJSON(t, &created, http.StatusCreated, "")(req)
			})

		err := client.Secret.Create(context.TODO(), CreateSecretRequest{Name: "trivy-token", Value: "s3cret"})
		require.NoError(t, err)
		require.Equal(t, map[string]interface{}{"name": "trivy-token", "value": "s3cret"}, created[0])
	})

	t.Run("Create conflict", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		httpmock.RegisterResponder(http.MethodPost, "http://localhost/api/v2/secrets",
			httpmock.NewStringResponder(http.StatusConflict, `{"title":"Resource already exists","status":409}`))

		err := client.Secret.Create(context.TODO(), CreateSecretRequest{Name: "trivy-token", Value: "s3cret"})
		var apiErr *APIError
		require.True(t, errors.As(err, &apiErr))
		require.Equal(t, http.StatusConflict, apiErr.StatusCode)
	})

	t.Run("Get", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		httpmock.RegisterResponder(http.MethodGet, "http://localhost/api/v2/secrets/trivy-token",
			httpmock.NewStringResponder(http.StatusOK, `{"name":"trivy-token","description":"d","created_at":1791874899788}`))

		secret, err := client.Secret.Get(context.TODO(), "trivy-token")
		require.NoError(t, err)
		require.Equal(t, Secret{Name: "trivy-token", Description: "d", CreatedAt: 1791874899788}, secret)
	})

	t.Run("Update only sends the fields that are set", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		var updated []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPatch, "http://localhost/api/v2/secrets/trivy-token",
			captureJSON(t, &updated, http.StatusNoContent, ""))

		value := "rotated"
		require.NoError(t, client.Secret.Update(context.TODO(), "trivy-token", UpdateSecretRequest{Value: &value}))
		require.Equal(t, map[string]interface{}{"value": "rotated"}, updated[0])
	})

	t.Run("Delete", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		httpmock.RegisterResponder(http.MethodDelete, "http://localhost/api/v2/secrets/trivy-token",
			httpmock.NewStringResponder(http.StatusNoContent, ""))

		require.NoError(t, client.Secret.Delete(context.TODO(), "trivy-token"))
	})
}

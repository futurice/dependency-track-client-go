package dtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestClient_doRequest_noBody(t *testing.T) {
	client := newTestClient(t, "5.2.0")

	httpmock.RegisterResponder(http.MethodPut, "http://localhost/api/v2/thing",
		httpmock.NewStringResponder(http.StatusNoContent, ""))

	req, err := client.newRequest(context.TODO(), http.MethodPut, "/api/v2/thing")
	require.NoError(t, err)

	var out map[string]interface{}
	_, err = client.doRequest(req, &out)
	require.NoError(t, err)
	require.Nil(t, out)
}

func TestClient_doRequest_notModified(t *testing.T) {
	t.Run("success without a body to decode", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		httpmock.RegisterResponder(http.MethodPut, "http://localhost/api/v2/thing",
			httpmock.NewStringResponder(http.StatusNotModified, ""))

		req, err := client.newRequest(context.TODO(), http.MethodPut, "/api/v2/thing")
		require.NoError(t, err)

		_, err = client.doRequest(req, nil)
		require.NoError(t, err)
	})

	t.Run("APIError when the caller decodes a body", func(t *testing.T) {
		client := newTestClient(t, "4.14.5")

		httpmock.RegisterResponder(http.MethodPost, "http://localhost/api/v1/notification/rule/x/project/y",
			httpmock.NewStringResponder(http.StatusNotModified, ""))

		req, err := client.newRequest(context.TODO(), http.MethodPost, "/api/v1/notification/rule/x/project/y")
		require.NoError(t, err)

		var out NotificationRule
		_, err = client.doRequest(req, &out)

		var apiErr *APIError
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotModified, apiErr.StatusCode)
	})
}

func TestClient_newRequest_accept(t *testing.T) {
	client := newTestClient(t, "5.2.0")

	v1, err := client.newRequest(context.TODO(), http.MethodGet, "/api/v1/team")
	require.NoError(t, err)
	require.Equal(t, "application/json", v1.Header.Get("Accept"))

	v2, err := client.newRequest(context.TODO(), http.MethodPost, "/api/v2/secrets")
	require.NoError(t, err)
	require.Equal(t, "application/json, application/problem+json", v2.Header.Get("Accept"))
}

package dtrack

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

// newTestClient returns a client that talks to httpmock and reports the given server version.
func newTestClient(t *testing.T, serverVersion string) *Client {
	t.Helper()

	client, err := NewClient("http://localhost", WithServerVersion(serverVersion))
	require.NoError(t, err)

	httpmock.ActivateNonDefault(client.httpClient)
	t.Cleanup(httpmock.DeactivateAndReset)

	return client
}

// captureJSON returns a responder that records each JSON request body and answers with status and responseBody.
func captureJSON(t *testing.T, bodies *[]map[string]interface{}, status int, responseBody string) httpmock.Responder {
	t.Helper()

	return func(req *http.Request) (*http.Response, error) {
		raw, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		var body map[string]interface{}
		require.NoError(t, json.Unmarshal(raw, &body))
		*bodies = append(*bodies, body)

		return httpmock.NewStringResponse(status, responseBody), nil
	}
}

func keys(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

package dtrack

import (
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

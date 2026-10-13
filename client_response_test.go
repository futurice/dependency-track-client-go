package dtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestClient_doRequest_noBody(t *testing.T) {
	for _, status := range []int{http.StatusNoContent, http.StatusNotModified} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			client := newTestClient(t, "5.2.0")

			httpmock.RegisterResponder(http.MethodPut, "http://localhost/api/v2/thing",
				httpmock.NewStringResponder(status, ""))

			req, err := client.newRequest(context.TODO(), http.MethodPut, "/api/v2/thing")
			require.NoError(t, err)

			var out map[string]interface{}
			_, err = client.doRequest(req, &out)
			require.NoError(t, err)
			require.Nil(t, out)
		})
	}
}

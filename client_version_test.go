package dtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestClient_ServerVersion(t *testing.T) {
	t.Run("fetched once and cached", func(t *testing.T) {
		client, err := NewClient("http://localhost")
		require.NoError(t, err)

		httpmock.ActivateNonDefault(client.httpClient)
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(http.MethodGet, "http://localhost/api/version",
			httpmock.NewStringResponder(http.StatusOK, `{"version":"5.2.0"}`))

		for i := 0; i < 3; i++ {
			version, err := client.ServerVersion(context.TODO())
			require.NoError(t, err)
			require.Equal(t, "5.2.0", version)
		}
		require.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("WithServerVersion skips the request", func(t *testing.T) {
		client := newTestClient(t, "4.14.5")

		version, err := client.ServerVersion(context.TODO())
		require.NoError(t, err)
		require.Equal(t, "4.14.5", version)
		require.Equal(t, 0, httpmock.GetTotalCallCount())
	})

	t.Run("a failed fetch is not cached", func(t *testing.T) {
		client, err := NewClient("http://localhost")
		require.NoError(t, err)

		httpmock.ActivateNonDefault(client.httpClient)
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(http.MethodGet, "http://localhost/api/version",
			httpmock.NewStringResponder(http.StatusServiceUnavailable, ""))
		_, err = client.ServerVersion(context.TODO())
		require.Error(t, err)

		httpmock.RegisterResponder(http.MethodGet, "http://localhost/api/version",
			httpmock.NewStringResponder(http.StatusOK, `{"version":"5.2.0"}`))
		version, err := client.ServerVersion(context.TODO())
		require.NoError(t, err)
		require.Equal(t, "5.2.0", version)
	})
}

func TestClient_isV5(t *testing.T) {
	for version, want := range map[string]bool{
		"4.14.5":          false,
		"4.13.0-SNAPSHOT": false,
		"5.0.0":           true,
		"5.2.0":           true,
		"5.3.0-SNAPSHOT":  true,
		"6.0.0":           true,
	} {
		t.Run(version, func(t *testing.T) {
			client := newTestClient(t, version)

			got, err := client.isV5(context.TODO())
			require.NoError(t, err)
			require.Equal(t, want, got)
		})
	}

	t.Run("unparseable", func(t *testing.T) {
		client := newTestClient(t, "unknown")

		_, err := client.isV5(context.TODO())
		require.Error(t, err)
	})
}

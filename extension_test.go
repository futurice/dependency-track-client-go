package dtrack

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const trivyConfigURL = "http://localhost/api/v2/extension-points/vuln-analyzer/extensions/trivy/config"

func TestExtensionService(t *testing.T) {
	t.Run("GetConfig", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		httpmock.RegisterResponder(http.MethodGet, trivyConfigURL,
			httpmock.NewStringResponder(http.StatusOK, `{"config":{"enabled":true,"apiToken":"trivy-token"}}`))

		config, err := client.Extension.GetConfig(context.TODO(), "vuln-analyzer", "trivy")
		require.NoError(t, err)
		require.JSONEq(t, `{"enabled":true,"apiToken":"trivy-token"}`, string(config))
	})

	t.Run("UpdateConfig wraps the config", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		var updated []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPut, trivyConfigURL, captureJSON(t, &updated, http.StatusNoContent, ""))

		err := client.Extension.UpdateConfig(context.TODO(), "vuln-analyzer", "trivy", json.RawMessage(`{"enabled":true}`))
		require.NoError(t, err)
		require.Equal(t, map[string]interface{}{"config": map[string]interface{}{"enabled": true}}, updated[0])
	})

	t.Run("UpdateConfig with no change", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		httpmock.RegisterResponder(http.MethodPut, trivyConfigURL, httpmock.NewStringResponder(http.StatusNotModified, ""))

		err := client.Extension.UpdateConfig(context.TODO(), "vuln-analyzer", "trivy", json.RawMessage(`{"enabled":true}`))
		require.NoError(t, err)
	})

	t.Run("GetConfigSchema", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		httpmock.RegisterResponder(http.MethodGet, trivyConfigURL+"-schema",
			httpmock.NewStringResponder(http.StatusOK, `{"$schema":"https://json-schema.org/draft/2020-12/schema","required":["enabled"]}`))

		schema, err := client.Extension.GetConfigSchema(context.TODO(), "vuln-analyzer", "trivy")
		require.NoError(t, err)
		require.Contains(t, string(schema), `"required":["enabled"]`)
	})

	t.Run("GetConfigSchema without a schema", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		httpmock.RegisterResponder(http.MethodGet, trivyConfigURL+"-schema", httpmock.NewStringResponder(http.StatusNoContent, ""))

		schema, err := client.Extension.GetConfigSchema(context.TODO(), "vuln-analyzer", "trivy")
		require.NoError(t, err)
		require.Nil(t, schema)
	})
}

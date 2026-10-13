package dtrack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ExtensionService manages v5 extension config, which replaces v4 config properties such as the Trivy analyzer's.
type ExtensionService struct {
	client *Client
}

type extensionConfig struct {
	Config json.RawMessage `json:"config"`
}

func extensionPath(extensionPoint, extension, suffix string) string {
	return fmt.Sprintf("/api/v2/extension-points/%s/extensions/%s/%s",
		url.PathEscape(extensionPoint), url.PathEscape(extension), suffix)
}

// GetConfig returns an extension's config as a JSON object, holding server defaults until updated.
func (es ExtensionService) GetConfig(ctx context.Context, extensionPoint, extension string) (config json.RawMessage, err error) {
	req, err := es.client.newRequest(ctx, http.MethodGet, extensionPath(extensionPoint, extension, "config"))
	if err != nil {
		return
	}

	var res extensionConfig
	_, err = es.client.doRequest(req, &res)
	return res.Config, err
}

// UpdateConfig replaces an extension's config, or returns a 400 APIError when it fails the config schema.
func (es ExtensionService) UpdateConfig(ctx context.Context, extensionPoint, extension string, config json.RawMessage) (err error) {
	req, err := es.client.newRequest(ctx, http.MethodPut, extensionPath(extensionPoint, extension, "config"),
		withBody(extensionConfig{Config: config}))
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, nil)
	return
}

// GetConfigSchema returns the JSON schema an extension's config must match, or nil if it has none.
func (es ExtensionService) GetConfigSchema(ctx context.Context, extensionPoint, extension string) (schema json.RawMessage, err error) {
	req, err := es.client.newRequest(ctx, http.MethodGet, extensionPath(extensionPoint, extension, "config-schema"))
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, &schema)
	return
}

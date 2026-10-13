package dtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const ruleURL = "http://localhost/api/v1/notification/rule"

var testRuleUUID = uuid.MustParse("6b8e8f3c-3c39-4e4a-9a4f-6d0f1e3f6a11")

func testRule() NotificationRule {
	return NotificationRule{
		Name:                 "alerts",
		Scope:                "PORTFOLIO",
		NotificationLevel:    "INFORMATIONAL",
		Publisher:            NotificationPublisher{UUID: testPublisherUUID},
		Enabled:              true,
		NotifyChildren:       true,
		LogSuccessfulPublish: true,
		NotifyOn:             []string{"NEW_VULNERABILITY"},
		PublisherConfig:      `{"destinationUrl":"https://example.com"}`,
	}
}

const ruleResponse = `{"uuid":"6b8e8f3c-3c39-4e4a-9a4f-6d0f1e3f6a11","name":"alerts","notificationLevel":"INFORMATIONAL"}`

func TestNotificationService_CreateRule(t *testing.T) {
	t.Run("v4 sends the whole rule, then updates", func(t *testing.T) {
		client := newTestClient(t, "4.14.5")

		var created, updated []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPut, ruleURL, captureJSON(t, &created, http.StatusCreated, ruleResponse))
		httpmock.RegisterResponder(http.MethodPost, ruleURL, captureJSON(t, &updated, http.StatusOK, ruleResponse))

		rule, err := client.Notification.CreateRule(context.TODO(), testRule())
		require.NoError(t, err)
		require.Equal(t, testRuleUUID, rule.UUID)

		require.Len(t, created, 1)
		require.Equal(t, "INFORMATIONAL", created[0]["notificationLevel"])

		// v4 create ignores notifyOn and publisherConfig, so the update carries them.
		require.Len(t, updated, 1)
		require.Equal(t, testRuleUUID.String(), updated[0]["uuid"])
		require.Equal(t, []interface{}{"NEW_VULNERABILITY"}, updated[0]["notifyOn"])
		require.Equal(t, `{"destinationUrl":"https://example.com"}`, updated[0]["publisherConfig"])
	})

	t.Run("v5 creates with the four allowed fields, then updates", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		var created, updated []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPut, ruleURL, captureJSON(t, &created, http.StatusCreated, ruleResponse))
		httpmock.RegisterResponder(http.MethodPost, ruleURL, captureJSON(t, &updated, http.StatusOK, ruleResponse))

		rule, err := client.Notification.CreateRule(context.TODO(), testRule())
		require.NoError(t, err)
		require.Equal(t, testRuleUUID, rule.UUID)

		require.Len(t, created, 1)
		require.ElementsMatch(t, []string{"name", "scope", "level", "publisher"}, keys(created[0]))
		require.Equal(t, "INFORMATIONAL", created[0]["level"])
		require.Equal(t, map[string]interface{}{"uuid": testPublisherUUID.String()}, created[0]["publisher"])

		require.Len(t, updated, 1)
		require.Equal(t, testRuleUUID.String(), updated[0]["uuid"])
		require.Equal(t, true, updated[0]["enabled"])
		require.Equal(t, []interface{}{"NEW_VULNERABILITY"}, updated[0]["notifyOn"])
		require.Equal(t, `{"destinationUrl":"https://example.com"}`, updated[0]["publisherConfig"])
	})
}

func TestNotificationService_CreateRuleFailedUpdate(t *testing.T) {
	client := newTestClient(t, "5.2.0")

	httpmock.RegisterResponder(http.MethodPut, ruleURL, httpmock.NewStringResponder(http.StatusCreated, ruleResponse))
	httpmock.RegisterResponder(http.MethodPost, ruleURL, httpmock.NewStringResponder(http.StatusBadRequest, `{"status":400}`))

	rule, err := client.Notification.CreateRule(context.TODO(), testRule())
	require.Error(t, err)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	require.Equal(t, testRuleUUID, rule.UUID)
}

func TestNotificationService_CreateRuleKeepsDefaultPublisherConfig(t *testing.T) {
	client := newTestClient(t, "5.2.0")

	createResponse := `{"uuid":"6b8e8f3c-3c39-4e4a-9a4f-6d0f1e3f6a11","name":"alerts","publisherConfig":"{\"destinationUrl\":\"https://example.com\"}"}`

	var created, updated []map[string]interface{}
	httpmock.RegisterResponder(http.MethodPut, ruleURL, captureJSON(t, &created, http.StatusCreated, createResponse))
	httpmock.RegisterResponder(http.MethodPost, ruleURL, captureJSON(t, &updated, http.StatusOK, ruleResponse))

	rule := testRule()
	rule.PublisherConfig = ""
	_, err := client.Notification.CreateRule(context.TODO(), rule)
	require.NoError(t, err)

	require.Len(t, updated, 1)
	require.Equal(t, `{"destinationUrl":"https://example.com"}`, updated[0]["publisherConfig"])
}

func TestNotificationService_UpdateRule(t *testing.T) {
	t.Run("v5 sends no publisher or projects", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		var updated []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPost, ruleURL, captureJSON(t, &updated, http.StatusOK, ruleResponse))

		rule := testRule()
		rule.UUID = testRuleUUID
		rule.Projects = []Project{{UUID: uuid.New()}}
		_, err := client.Notification.UpdateRule(context.TODO(), rule)
		require.NoError(t, err)

		require.Len(t, updated, 1)
		require.ElementsMatch(t, []string{
			"uuid", "name", "scope", "level", "enabled", "notifyChildren",
			"logSuccessfulPublish", "notifyOn", "publisherConfig", "tags",
		}, keys(updated[0]))
	})

	t.Run("v5 sends the filter expression and tags it was given", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		var updated []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPost, ruleURL, captureJSON(t, &updated, http.StatusOK, ruleResponse))

		rule := testRule()
		rule.FilterExpression = `subject.project.name == "x"`
		rule.Tags = []Tag{{Name: "prod"}}
		_, err := client.Notification.UpdateRule(context.TODO(), rule)
		require.NoError(t, err)
		require.Equal(t, `subject.project.name == "x"`, updated[0]["filterExpression"])
		require.Equal(t, []interface{}{map[string]interface{}{"name": "prod"}}, updated[0]["tags"])
	})

	t.Run("v5 sends an empty notifyOn list, not null", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		var updated []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPost, ruleURL, captureJSON(t, &updated, http.StatusOK, ruleResponse))

		rule := testRule()
		rule.NotifyOn = nil
		_, err := client.Notification.UpdateRule(context.TODO(), rule)
		require.NoError(t, err)
		require.Equal(t, []interface{}{}, updated[0]["notifyOn"])
	})

	t.Run("v4 sends the rule as is", func(t *testing.T) {
		client := newTestClient(t, "4.14.5")

		var updated []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPost, ruleURL, captureJSON(t, &updated, http.StatusOK, ruleResponse))

		_, err := client.Notification.UpdateRule(context.TODO(), testRule())
		require.NoError(t, err)
		require.Contains(t, updated[0], "notificationLevel")
		require.Contains(t, updated[0], "publisher")
	})
}

func TestNotificationService_GetAllRules(t *testing.T) {
	client := newTestClient(t, "5.2.0")

	httpmock.RegisterResponder(http.MethodGet, ruleURL,
		func(req *http.Request) (*http.Response, error) {
			require.Equal(t, "2", req.URL.Query().Get("pageNumber"))
			require.Equal(t, "50", req.URL.Query().Get("pageSize"))

			res := httpmock.NewStringResponse(http.StatusOK, `[`+ruleResponse+`]`)
			res.Header.Set("X-Total-Count", "51")
			return res, nil
		})

	page, err := client.Notification.GetAllRules(context.TODO(), PageOptions{PageNumber: 2, PageSize: 50})
	require.NoError(t, err)
	require.Equal(t, 51, page.TotalCount)
	require.Len(t, page.Items, 1)
	require.Equal(t, testRuleUUID, page.Items[0].UUID)
}

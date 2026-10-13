package dtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const publisherURL = "http://localhost/api/v1/notification/publisher"

var testPublisherUUID = uuid.MustParse("1c2b6a1e-8f0e-4a43-9a3e-2b8f3c9d0e22")

func TestNotificationService_CreatePublisher(t *testing.T) {
	publisher := NotificationPublisher{
		Name:             "JSM",
		PublisherClass:   "org.dependencytrack.notification.publisher.WebhookPublisher",
		ExtensionName:    "webhook",
		Template:         "{}",
		TemplateMimeType: "application/json",
	}
	response := `{"uuid":"1c2b6a1e-8f0e-4a43-9a3e-2b8f3c9d0e22","name":"JSM","extensionName":"webhook"}`

	t.Run("v5 sends extensionName and no v4 fields", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		var created []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPut, publisherURL, captureJSON(t, &created, http.StatusCreated, response))

		got, err := client.Notification.CreatePublisher(context.TODO(), publisher)
		require.NoError(t, err)
		require.Equal(t, "webhook", got.ExtensionName)

		require.ElementsMatch(t, []string{"name", "extensionName", "template", "templateMimeType"}, keys(created[0]))
	})

	t.Run("v5 update includes the uuid", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")

		var updated []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPost, publisherURL, captureJSON(t, &updated, http.StatusOK, response))

		p := publisher
		p.UUID = testPublisherUUID
		_, err := client.Notification.UpdatePublisher(context.TODO(), p)
		require.NoError(t, err)
		require.Equal(t, testPublisherUUID.String(), updated[0]["uuid"])
		require.NotContains(t, updated[0], "publisherClass")
	})

	t.Run("v4 sends publisherClass", func(t *testing.T) {
		client := newTestClient(t, "4.14.5")

		var created []map[string]interface{}
		httpmock.RegisterResponder(http.MethodPut, publisherURL, captureJSON(t, &created, http.StatusCreated, response))

		p := publisher
		p.ExtensionName = ""
		_, err := client.Notification.CreatePublisher(context.TODO(), p)
		require.NoError(t, err)
		require.Equal(t, publisher.PublisherClass, created[0]["publisherClass"])
		require.NotContains(t, created[0], "extensionName")
	})
}

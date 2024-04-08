package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type NotificationPublisher struct {
	UUID             uuid.UUID `json:"uuid"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	PublisherClass   string    `json:"publisherClass"`
	Template         string    `json:"template"`
	TemplateMimeType string    `json:"templateMimeType"`
	DefaultPublisher bool      `json:"defaultPublisher"`
}

type NotificationRule struct {
	UUID                 uuid.UUID             `json:"uuid"`
	Name                 string                `json:"name"`
	Publisher            NotificationPublisher `json:"publisher"`
	Scope                string                `json:"scope"`
	NotificationLevel    string                `json:"notificationLevel"`
	NotifyOn             []string              `json:"notifyOn"`
	Enabled              bool                  `json:"enabled"`
	NotifyChildren       bool                  `json:"notifyChildren"`
	LogSuccessfulPublish bool                  `json:"logSuccessfulPublish"`
	PublisherConfig      string                `json:"publisherConfig,omitempty"`
}

type NotificationService struct {
	client *Client
}

func (ns NotificationService) GetAllPublishers(ctx context.Context) (p []NotificationPublisher, err error) {
	req, err := ns.client.newRequest(ctx, http.MethodGet, "/api/v1/notification/publisher")
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &p)
	return
}

func (ns NotificationService) CreatePublisher(ctx context.Context, publisher NotificationPublisher) (r NotificationPublisher, err error) {
	req, err := ns.client.newRequest(ctx, http.MethodPut, "/api/v1/notification/publisher", withBody(publisher))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

func (ns NotificationService) UpdatePublisher(ctx context.Context, publisher NotificationPublisher) (r NotificationPublisher, err error) {
	req, err := ns.client.newRequest(ctx, http.MethodPost, "/api/v1/notification/publisher", withBody(publisher))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

func (ns NotificationService) DeletePublisher(ctx context.Context, ruleUuid uuid.UUID) (err error) {
	req, err := ns.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/notification/publisher/%s", ruleUuid.String()))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, nil)
	return
}

func (ns NotificationService) GetAllRules(ctx context.Context) (r []NotificationRule, err error) {
	req, err := ns.client.newRequest(ctx, http.MethodGet, "/api/v1/notification/rule")
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

func (ns NotificationService) CreateRule(ctx context.Context, rule NotificationRule) (r NotificationRule, err error) {
	req, err := ns.client.newRequest(ctx, http.MethodPut, "/api/v1/notification/rule", withBody(rule))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

func (ns NotificationService) UpdateRule(ctx context.Context, rule NotificationRule) (r NotificationRule, err error) {
	req, err := ns.client.newRequest(ctx, http.MethodPost, "/api/v1/notification/rule", withBody(rule))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

func (ns NotificationService) DeleteRule(ctx context.Context, ruleUuid uuid.UUID) (err error) {
	req, err := ns.client.newRequest(ctx, http.MethodDelete, "/api/v1/notification/rule", withBody(struct {
		UUID uuid.UUID `json:"uuid"`
	}{UUID: ruleUuid}))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, nil)
	return
}

func (ns NotificationService) AddProjectToRule(ctx context.Context, ruleUuid, projectUuid uuid.UUID) (r NotificationRule, err error) {
	req, err := ns.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/notification/rule/%s/project/%s", ruleUuid.String(), projectUuid.String()))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

func (ns NotificationService) DeleteProjectFromRule(ctx context.Context, ruleUuid, projectUuid uuid.UUID) (r NotificationRule, err error) {
	req, err := ns.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/notification/rule/%s/project/%s", ruleUuid.String(), projectUuid.String()))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

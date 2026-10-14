package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type NotificationPublisher struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	// PublisherClass is the publisher's Java class on Dependency-Track v4.
	PublisherClass string `json:"publisherClass,omitempty"`
	// ExtensionName replaces PublisherClass on Dependency-Track v5, e.g. "webhook" or "slack".
	ExtensionName    string `json:"extensionName,omitempty"`
	Template         string `json:"template"`
	TemplateMimeType string `json:"templateMimeType"`
	DefaultPublisher bool   `json:"defaultPublisher"`
}

// notificationPublisherRequestV5 is the v5 publisher create and update body, without publisherClass or defaultPublisher.
type notificationPublisherRequestV5 struct {
	UUID             *uuid.UUID `json:"uuid,omitempty"`
	Name             string     `json:"name"`
	Description      string     `json:"description,omitempty"`
	ExtensionName    string     `json:"extensionName"`
	Template         string     `json:"template,omitempty"`
	TemplateMimeType string     `json:"templateMimeType,omitempty"`
}

func newNotificationPublisherRequestV5(publisher NotificationPublisher, withUUID bool) notificationPublisherRequestV5 {
	req := notificationPublisherRequestV5{
		Name:             publisher.Name,
		Description:      publisher.Description,
		ExtensionName:    publisher.ExtensionName,
		Template:         publisher.Template,
		TemplateMimeType: publisher.TemplateMimeType,
	}
	if withUUID {
		req.UUID = &publisher.UUID
	}
	return req
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
	Projects             []Project             `json:"projects,omitempty"`
	// Both versions overwrite the rule's tags on update, and v5 its filter expression too. Copy them from
	// the current rule to keep them.
	FilterExpression string `json:"filterExpression,omitempty"`
	Tags             []Tag  `json:"tags"`
}

// createNotificationRuleRequestV5 is the v5 rule create body; v5 rejects any other field.
type createNotificationRuleRequestV5 struct {
	Name      string                `json:"name"`
	Scope     string                `json:"scope"`
	Level     string                `json:"level"`
	Publisher notificationRuleRefV5 `json:"publisher"`
}

type notificationRuleRefV5 struct {
	UUID uuid.UUID `json:"uuid"`
}

// updateNotificationRuleRequestV5 is the v5 rule update body, without publisher, projects or teams.
type updateNotificationRuleRequestV5 struct {
	UUID                 uuid.UUID `json:"uuid"`
	Name                 string    `json:"name"`
	Scope                string    `json:"scope"`
	Level                string    `json:"level"`
	Enabled              bool      `json:"enabled"`
	NotifyChildren       bool      `json:"notifyChildren"`
	LogSuccessfulPublish bool      `json:"logSuccessfulPublish"`
	NotifyOn             []string  `json:"notifyOn"`
	PublisherConfig      string    `json:"publisherConfig,omitempty"`
	FilterExpression     string    `json:"filterExpression,omitempty"`
	Tags                 []Tag     `json:"tags"`
}

func newUpdateNotificationRuleRequestV5(rule NotificationRule) updateNotificationRuleRequestV5 {
	notifyOn := rule.NotifyOn
	if notifyOn == nil {
		notifyOn = []string{}
	}
	tags := rule.Tags
	if tags == nil {
		tags = []Tag{}
	}

	return updateNotificationRuleRequestV5{
		UUID:                 rule.UUID,
		Name:                 rule.Name,
		Scope:                rule.Scope,
		Level:                rule.NotificationLevel,
		Enabled:              rule.Enabled,
		NotifyChildren:       rule.NotifyChildren,
		LogSuccessfulPublish: rule.LogSuccessfulPublish,
		NotifyOn:             notifyOn,
		PublisherConfig:      rule.PublisherConfig,
		FilterExpression:     rule.FilterExpression,
		Tags:                 tags,
	}
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
	v5, err := ns.client.isV5(ctx)
	if err != nil {
		return
	}

	var body interface{} = publisher
	if v5 {
		body = newNotificationPublisherRequestV5(publisher, false)
	}

	req, err := ns.client.newRequest(ctx, http.MethodPut, "/api/v1/notification/publisher", withBody(body))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

func (ns NotificationService) UpdatePublisher(ctx context.Context, publisher NotificationPublisher) (r NotificationPublisher, err error) {
	v5, err := ns.client.isV5(ctx)
	if err != nil {
		return
	}

	var body interface{} = publisher
	if v5 {
		body = newNotificationPublisherRequestV5(publisher, true)
	}

	req, err := ns.client.newRequest(ctx, http.MethodPost, "/api/v1/notification/publisher", withBody(body))
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

// GetAllRules returns one page of rules; use FetchAll, since v5 caps unpaged lists at 100.
func (ns NotificationService) GetAllRules(ctx context.Context, po PageOptions) (p Page[NotificationRule], err error) {
	return ns.getRules(ctx, "", po)
}

// SearchRules lists the rules whose name or publisher name contains searchText, ignoring case. The server
// treats a "." in searchText as any character, so filter the results yourself. Neither v4 nor v5 can fetch
// a single rule, so searching by name and then matching the UUID avoids listing every rule to find one.
func (ns NotificationService) SearchRules(ctx context.Context, searchText string, po PageOptions) (p Page[NotificationRule], err error) {
	return ns.getRules(ctx, searchText, po)
}

// getRules sends searchText only when it's set, so an empty search lists every rule like GetAllRules.
func (ns NotificationService) getRules(ctx context.Context, searchText string, po PageOptions) (p Page[NotificationRule], err error) {
	var params map[string]string
	if searchText != "" {
		params = map[string]string{"searchText": searchText}
	}

	req, err := ns.client.newRequest(ctx, http.MethodGet, "/api/v1/notification/rule", withParams(params), withPageOptions(po))
	if err != nil {
		return
	}

	res, err := ns.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}

// CreateRule creates a rule, then updates it, because create drops most settings on both v4 and v5.
// If the update fails, it returns the created rule with the error, so the caller can track or delete it.
func (ns NotificationService) CreateRule(ctx context.Context, rule NotificationRule) (r NotificationRule, err error) {
	v5, err := ns.client.isV5(ctx)
	if err != nil {
		return
	}

	var body interface{} = rule
	if v5 {
		body = createNotificationRuleRequestV5{
			Name:      rule.Name,
			Scope:     rule.Scope,
			Level:     rule.NotificationLevel,
			Publisher: notificationRuleRefV5{UUID: rule.Publisher.UUID},
		}
	}

	req, err := ns.client.newRequest(ctx, http.MethodPut, "/api/v1/notification/rule", withBody(body))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	if err != nil {
		return
	}

	rule.UUID = r.UUID
	if rule.PublisherConfig == "" {
		// v5 fills in a default config on create and rejects an update without one.
		rule.PublisherConfig = r.PublisherConfig
	}

	updated, err := ns.UpdateRule(ctx, rule)
	if err != nil {
		return r, fmt.Errorf("rule %s created, but its follow-up update failed: %w", r.UUID, err)
	}
	return updated, nil
}

func (ns NotificationService) UpdateRule(ctx context.Context, rule NotificationRule) (r NotificationRule, err error) {
	v5, err := ns.client.isV5(ctx)
	if err != nil {
		return
	}

	var body interface{} = rule
	if v5 {
		body = newUpdateNotificationRuleRequestV5(rule)
	}

	req, err := ns.client.newRequest(ctx, http.MethodPost, "/api/v1/notification/rule", withBody(body))
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

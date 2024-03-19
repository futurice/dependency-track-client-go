package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ACLMapping struct {
	Team    uuid.UUID `json:"team,omitempty"`
	Project uuid.UUID `json:"project,omitempty"`
}

type ACLMappingService struct {
	client *Client
}

func (as ACLMappingService) Get(ctx context.Context, teamUUID uuid.UUID) (mappings []Project, err error) {
	req, err := as.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/acl/team/%s", teamUUID.String()))
	if err != nil {
		return
	}

	_, err = as.client.doRequest(req, &mappings)
	return
}

func (as ACLMappingService) Create(ctx context.Context, mapping ACLMapping) (err error) {
	req, err := as.client.newRequest(ctx, http.MethodPut, "/api/v1/acl/mapping", withBody(mapping))
	if err != nil {
		return
	}

	_, err = as.client.doRequest(req, nil)
	return
}

func (as ACLMappingService) Delete(ctx context.Context, mapping ACLMapping) (err error) {
	req, err := as.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/acl/mapping/team/%s/project/%s", mapping.Team, mapping.Project))
	if err != nil {
		return
	}

	_, err = as.client.doRequest(req, nil)
	return
}

package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Team struct {
	UUID             uuid.UUID     `json:"uuid,omitempty"`
	Name             string        `json:"name,omitempty"`
	APIKeys          []APIKey      `json:"apiKeys,omitempty"`
	Permissions      []Permission  `json:"permissions,omitempty"`
	MappedOIDCGroups []OIDCMapping `json:"mappedOidcGroups,omitempty"`
}

type APIKey struct {
	PublicID  string `json:"publicId"`
	Key       string `json:"key"`
	MaskedKey string `json:"maskedKey"`
	Legacy    bool   `json:"legacy"`
	Created   uint   `json:"created"`
}

type TeamService struct {
	client *Client
}

func (ts TeamService) Get(ctx context.Context, teamUUID uuid.UUID) (t Team, err error) {
	req, err := ts.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/team/%s", teamUUID))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, &t)
	return
}

func (ts TeamService) GetAll(ctx context.Context, po PageOptions) (p Page[Team], err error) {
	req, err := ts.client.newRequest(ctx, http.MethodGet, "/api/v1/team", withPageOptions(po))
	if err != nil {
		return
	}

	res, err := ts.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}

// GetAPIKeys searches the team list for a team's keys (v4's single-team GET omits them), or returns a 404 APIError.
func (ts TeamService) GetAPIKeys(ctx context.Context, teamUUID uuid.UUID) (keys []APIKey, err error) {
	found := false
	err = ForEach(
		func(po PageOptions) (Page[Team], error) { return ts.GetAll(ctx, po) },
		func(team Team) error {
			if team.UUID == teamUUID {
				found = true
				keys = team.APIKeys
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, &APIError{StatusCode: http.StatusNotFound, Message: fmt.Sprintf("team %s not found", teamUUID)}
	}
	if keys == nil {
		keys = []APIKey{}
	}
	return keys, nil
}

func (ts TeamService) GenerateAPIKey(ctx context.Context, teamUUID uuid.UUID) (key APIKey, err error) {
	req, err := ts.client.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/team/%s/key", teamUUID))
	if err != nil {
		return
	}

	var apiKey APIKey
	_, err = ts.client.doRequest(req, &apiKey)
	key = apiKey
	return
}

func (ts TeamService) RegenerateAPIKey(ctx context.Context, oldAPIKey string) (APIKey, error) {
	req, err := ts.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/team/key/%s", oldAPIKey))
	if err != nil {
		return APIKey{}, err
	}

	var newAPIKey APIKey
	_, err = ts.client.doRequest(req, &newAPIKey)

	return newAPIKey, err
}

func (ts TeamService) DeleteAPIKey(ctx context.Context, apiKey string) (key string, err error) {
	req, err := ts.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/team/key/%s", apiKey))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TeamService) Create(ctx context.Context, team Team) (t Team, err error) {
	req, err := ts.client.newRequest(ctx, http.MethodPut, "/api/v1/team", withBody(team))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, &t)
	return
}

func (ts TeamService) Update(ctx context.Context, team Team) (t Team, err error) {
	req, err := ts.client.newRequest(ctx, http.MethodPost, "/api/v1/team", withBody(team))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, &t)
	return
}

func (ts TeamService) Delete(ctx context.Context, team Team) (err error) {
	req, err := ts.client.newRequest(ctx, http.MethodDelete, "/api/v1/team", withBody(team))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, nil)
	return
}

package dtrack

import (
	"context"
	"net/http"
)

type ConfigPropertyType string

const (
	ConfigPropertyTypeBoolean         ConfigPropertyType = "BOOLEAN"
	ConfigPropertyTypeInteger         ConfigPropertyType = "INTEGER"
	ConfigPropertyTypeNumber          ConfigPropertyType = "NUMBER"
	ConfigPropertyTypeString          ConfigPropertyType = "STRING"
	ConfigPropertyTypeEncryptedString ConfigPropertyType = "ENCRYPTEDSTRING"
	ConfigPropertyTypeTimestamp       ConfigPropertyType = "TIMESTAMP"
	ConfigPropertyTypeURL             ConfigPropertyType = "URL"
	ConfigPropertyTypeUUID            ConfigPropertyType = "UUID"
)

type ConfigProperty struct {
	GroupName   string				`json:"groupName"`
	Name        string 				`json:"propertyName"`
	Value       string 				`json:"propertyValue,omitempty"`
	Type        ConfigPropertyType	`json:"propertyType"`
	Description string 				`json:"description,omitempty"`
}

type ConfigService struct {
	client *Client
}

func (cs ConfigService) GetAll(ctx context.Context) (cps []ConfigProperty, err error) {
	req, err := cs.client.newRequest(ctx, http.MethodGet, "/api/v1/configProperty")
	if err != nil {
		return
	}
	_, err = cs.client.doRequest(req, &cps)
	return
}

func (cs ConfigService) Get(ctx context.Context, groupName, propertyName string) (cp ConfigProperty, err error) {
	cps, err := cs.GetAll(ctx)
	if err != nil {
		return
	}
	for _, cp := range cps {
		if cp.GroupName != groupName {
			continue
		}
		if cp.Name != propertyName {
			continue
		}
		return cp, nil
	}
	return
}

func (cs ConfigService) Update(ctx context.Context, config ConfigProperty) (cp ConfigProperty, err error) {
	req, err := cs.client.newRequest(ctx, http.MethodPost, "/api/v1/configProperty", withBody(config))
	if err != nil {
		return
	}
	_, err = cs.client.doRequest(req, &cp)
	return
}

func (cs ConfigService) UpdateAll(ctx context.Context, configs []ConfigProperty) (cps []ConfigProperty, err error) {
	req, err := cs.client.newRequest(ctx, http.MethodPost, "/api/v1/configProperty/aggregate", withBody(configs))
	if err != nil {
		return
	}
	_, err = cs.client.doRequest(req, &cps)
	return
}

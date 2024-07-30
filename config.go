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
	GroupName     string             `json:"groupName,omitempty"`
	PropertyName  string             `json:"propertyName,omitempty"`
	PropertyValue *string            `json:"propertyValue,omitempty"`
	PropertyType  ConfigPropertyType `json:"propertyType,omitempty"`
	Description   string             `json:"description,omitempty"`
}

type ConfigService struct {
	client *Client
}

func (ps ConfigService) GetAllConfigProperties(ctx context.Context) (cp []ConfigProperty, err error) {
	req, err := ps.client.newRequest(ctx, http.MethodGet, "/api/v1/configProperty")
	if err != nil {
		return
	}

	_, err = ps.client.doRequest(req, &cp)
	if err != nil {
		return
	}

	return
}

type SetConfigPropertyRequest struct {
	GroupName     string `json:"groupName,omitempty"`
	PropertyName  string `json:"propertyName,omitempty"`
	PropertyValue string `json:"propertyValue,omitempty"`
}

func (ps ConfigService) SetConfigProperty(ctx context.Context, setConfigPropertyRequest SetConfigPropertyRequest) (p ConfigProperty, err error) {
	req, err := ps.client.newRequest(ctx, http.MethodPost, "/api/v1/configProperty", withBody(setConfigPropertyRequest))
	if err != nil {
		return
	}

	_, err = ps.client.doRequest(req, &p)
	return
}

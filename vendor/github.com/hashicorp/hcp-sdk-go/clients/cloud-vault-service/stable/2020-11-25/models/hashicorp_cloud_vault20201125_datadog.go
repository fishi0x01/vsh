// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpCloudVault20201125Datadog hashicorp cloud vault 20201125 datadog
//
// swagger:model hashicorp.cloud.vault_20201125.Datadog
type HashicorpCloudVault20201125Datadog struct {

	// api key
	APIKey string `json:"api_key,omitempty"`

	// region
	Region string `json:"region,omitempty"`
}

// Validate validates this hashicorp cloud vault 20201125 datadog
func (m *HashicorpCloudVault20201125Datadog) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hashicorp cloud vault 20201125 datadog based on context it is used
func (m *HashicorpCloudVault20201125Datadog) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpCloudVault20201125Datadog) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpCloudVault20201125Datadog) UnmarshalBinary(b []byte) error {
	var res HashicorpCloudVault20201125Datadog
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

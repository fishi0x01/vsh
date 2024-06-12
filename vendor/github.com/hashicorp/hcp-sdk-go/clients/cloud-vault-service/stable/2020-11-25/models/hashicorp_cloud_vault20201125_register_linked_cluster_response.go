// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpCloudVault20201125RegisterLinkedClusterResponse hashicorp cloud vault 20201125 register linked cluster response
//
// swagger:model hashicorp.cloud.vault_20201125.RegisterLinkedClusterResponse
type HashicorpCloudVault20201125RegisterLinkedClusterResponse struct {

	// client id
	ClientID string `json:"client_id,omitempty"`

	// client secret
	ClientSecret string `json:"client_secret,omitempty"`

	// cluster id
	ClusterID string `json:"cluster_id,omitempty"`

	// resource id
	ResourceID string `json:"resource_id,omitempty"`
}

// Validate validates this hashicorp cloud vault 20201125 register linked cluster response
func (m *HashicorpCloudVault20201125RegisterLinkedClusterResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hashicorp cloud vault 20201125 register linked cluster response based on context it is used
func (m *HashicorpCloudVault20201125RegisterLinkedClusterResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpCloudVault20201125RegisterLinkedClusterResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpCloudVault20201125RegisterLinkedClusterResponse) UnmarshalBinary(b []byte) error {
	var res HashicorpCloudVault20201125RegisterLinkedClusterResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
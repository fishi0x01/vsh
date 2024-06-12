// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpCloudResourcemanagerResourceTestIamPermissionsResponse see ResourceService.TestIamPermissions
//
// swagger:model hashicorp.cloud.resourcemanager.ResourceTestIamPermissionsResponse
type HashicorpCloudResourcemanagerResourceTestIamPermissionsResponse struct {

	// AllowedPermissions are a subset of the request permissions the calling principal has for the resource.
	AllowedPermissions []string `json:"allowed_permissions"`
}

// Validate validates this hashicorp cloud resourcemanager resource test iam permissions response
func (m *HashicorpCloudResourcemanagerResourceTestIamPermissionsResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hashicorp cloud resourcemanager resource test iam permissions response based on context it is used
func (m *HashicorpCloudResourcemanagerResourceTestIamPermissionsResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpCloudResourcemanagerResourceTestIamPermissionsResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpCloudResourcemanagerResourceTestIamPermissionsResponse) UnmarshalBinary(b []byte) error {
	var res HashicorpCloudResourcemanagerResourceTestIamPermissionsResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
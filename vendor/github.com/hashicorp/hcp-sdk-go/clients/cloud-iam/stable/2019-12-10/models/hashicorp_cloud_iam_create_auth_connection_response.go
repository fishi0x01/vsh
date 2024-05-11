// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpCloudIamCreateAuthConnectionResponse CreateAuthConnectionResponse for creating a new auth connection.
//
// swagger:model hashicorp.cloud.iam.CreateAuthConnectionResponse
type HashicorpCloudIamCreateAuthConnectionResponse struct {

	// auth_connection that was created
	AuthConnection *HashicorpCloudIamAuthConnection `json:"auth_connection,omitempty"`
}

// Validate validates this hashicorp cloud iam create auth connection response
func (m *HashicorpCloudIamCreateAuthConnectionResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAuthConnection(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpCloudIamCreateAuthConnectionResponse) validateAuthConnection(formats strfmt.Registry) error {
	if swag.IsZero(m.AuthConnection) { // not required
		return nil
	}

	if m.AuthConnection != nil {
		if err := m.AuthConnection.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("auth_connection")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("auth_connection")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp cloud iam create auth connection response based on the context it is used
func (m *HashicorpCloudIamCreateAuthConnectionResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAuthConnection(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpCloudIamCreateAuthConnectionResponse) contextValidateAuthConnection(ctx context.Context, formats strfmt.Registry) error {

	if m.AuthConnection != nil {

		if swag.IsZero(m.AuthConnection) { // not required
			return nil
		}

		if err := m.AuthConnection.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("auth_connection")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("auth_connection")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpCloudIamCreateAuthConnectionResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpCloudIamCreateAuthConnectionResponse) UnmarshalBinary(b []byte) error {
	var res HashicorpCloudIamCreateAuthConnectionResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

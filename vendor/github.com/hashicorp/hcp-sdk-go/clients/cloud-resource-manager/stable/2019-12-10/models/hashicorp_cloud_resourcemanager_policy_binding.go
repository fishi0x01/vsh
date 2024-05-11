// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpCloudResourcemanagerPolicyBinding PolicyBinding binds a list of members to a role.
//
// swagger:model hashicorp.cloud.resourcemanager.PolicyBinding
type HashicorpCloudResourcemanagerPolicyBinding struct {

	// Members is a set of identities that are associated with the role.
	Members []*HashicorpCloudResourcemanagerPolicyBindingMember `json:"members"`

	// RoleId is the globally unique identifier for the role.
	// e.g. projects/30788e25-a942-470a-9a75-128cc7b97fde/roles/my-custom-project-role
	RoleID string `json:"role_id,omitempty"`
}

// Validate validates this hashicorp cloud resourcemanager policy binding
func (m *HashicorpCloudResourcemanagerPolicyBinding) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMembers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpCloudResourcemanagerPolicyBinding) validateMembers(formats strfmt.Registry) error {
	if swag.IsZero(m.Members) { // not required
		return nil
	}

	for i := 0; i < len(m.Members); i++ {
		if swag.IsZero(m.Members[i]) { // not required
			continue
		}

		if m.Members[i] != nil {
			if err := m.Members[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("members" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("members" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this hashicorp cloud resourcemanager policy binding based on the context it is used
func (m *HashicorpCloudResourcemanagerPolicyBinding) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMembers(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpCloudResourcemanagerPolicyBinding) contextValidateMembers(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Members); i++ {

		if m.Members[i] != nil {

			if swag.IsZero(m.Members[i]) { // not required
				return nil
			}

			if err := m.Members[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("members" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("members" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpCloudResourcemanagerPolicyBinding) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpCloudResourcemanagerPolicyBinding) UnmarshalBinary(b []byte) error {
	var res HashicorpCloudResourcemanagerPolicyBinding
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

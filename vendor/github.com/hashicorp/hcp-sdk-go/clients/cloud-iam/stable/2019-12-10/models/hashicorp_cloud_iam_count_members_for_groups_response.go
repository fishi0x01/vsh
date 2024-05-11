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

// HashicorpCloudIamCountMembersForGroupsResponse hashicorp cloud iam count members for groups response
//
// swagger:model hashicorp.cloud.iam.CountMembersForGroupsResponse
type HashicorpCloudIamCountMembersForGroupsResponse struct {

	// groups_counts is a list of member counts per group.
	GroupsCounts []*HashicorpCloudIamCountMembersForGroupsResponseMembersCount `json:"groups_counts"`
}

// Validate validates this hashicorp cloud iam count members for groups response
func (m *HashicorpCloudIamCountMembersForGroupsResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateGroupsCounts(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpCloudIamCountMembersForGroupsResponse) validateGroupsCounts(formats strfmt.Registry) error {
	if swag.IsZero(m.GroupsCounts) { // not required
		return nil
	}

	for i := 0; i < len(m.GroupsCounts); i++ {
		if swag.IsZero(m.GroupsCounts[i]) { // not required
			continue
		}

		if m.GroupsCounts[i] != nil {
			if err := m.GroupsCounts[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("groups_counts" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("groups_counts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this hashicorp cloud iam count members for groups response based on the context it is used
func (m *HashicorpCloudIamCountMembersForGroupsResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateGroupsCounts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpCloudIamCountMembersForGroupsResponse) contextValidateGroupsCounts(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.GroupsCounts); i++ {

		if m.GroupsCounts[i] != nil {

			if swag.IsZero(m.GroupsCounts[i]) { // not required
				return nil
			}

			if err := m.GroupsCounts[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("groups_counts" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("groups_counts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpCloudIamCountMembersForGroupsResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpCloudIamCountMembersForGroupsResponse) UnmarshalBinary(b []byte) error {
	var res HashicorpCloudIamCountMembersForGroupsResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

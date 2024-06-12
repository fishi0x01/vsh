// Code generated by go-swagger; DO NOT EDIT.

package vault_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/models"
)

// NewLockParams creates a new LockParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewLockParams() *LockParams {
	return &LockParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewLockParamsWithTimeout creates a new LockParams object
// with the ability to set a timeout on a request.
func NewLockParamsWithTimeout(timeout time.Duration) *LockParams {
	return &LockParams{
		timeout: timeout,
	}
}

// NewLockParamsWithContext creates a new LockParams object
// with the ability to set a context for a request.
func NewLockParamsWithContext(ctx context.Context) *LockParams {
	return &LockParams{
		Context: ctx,
	}
}

// NewLockParamsWithHTTPClient creates a new LockParams object
// with the ability to set a custom HTTPClient for a request.
func NewLockParamsWithHTTPClient(client *http.Client) *LockParams {
	return &LockParams{
		HTTPClient: client,
	}
}

/*
LockParams contains all the parameters to send to the API endpoint

	for the lock operation.

	Typically these are written to a http.Request.
*/
type LockParams struct {

	// Body.
	Body *models.HashicorpCloudVault20201125LockRequest

	// ClusterID.
	ClusterID string

	/* LocationOrganizationID.

	   organization_id is the id of the organization.
	*/
	LocationOrganizationID string

	/* LocationProjectID.

	   project_id is the projects id.
	*/
	LocationProjectID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the lock params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *LockParams) WithDefaults() *LockParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the lock params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *LockParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the lock params
func (o *LockParams) WithTimeout(timeout time.Duration) *LockParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the lock params
func (o *LockParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the lock params
func (o *LockParams) WithContext(ctx context.Context) *LockParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the lock params
func (o *LockParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the lock params
func (o *LockParams) WithHTTPClient(client *http.Client) *LockParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the lock params
func (o *LockParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the lock params
func (o *LockParams) WithBody(body *models.HashicorpCloudVault20201125LockRequest) *LockParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the lock params
func (o *LockParams) SetBody(body *models.HashicorpCloudVault20201125LockRequest) {
	o.Body = body
}

// WithClusterID adds the clusterID to the lock params
func (o *LockParams) WithClusterID(clusterID string) *LockParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the lock params
func (o *LockParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithLocationOrganizationID adds the locationOrganizationID to the lock params
func (o *LockParams) WithLocationOrganizationID(locationOrganizationID string) *LockParams {
	o.SetLocationOrganizationID(locationOrganizationID)
	return o
}

// SetLocationOrganizationID adds the locationOrganizationId to the lock params
func (o *LockParams) SetLocationOrganizationID(locationOrganizationID string) {
	o.LocationOrganizationID = locationOrganizationID
}

// WithLocationProjectID adds the locationProjectID to the lock params
func (o *LockParams) WithLocationProjectID(locationProjectID string) *LockParams {
	o.SetLocationProjectID(locationProjectID)
	return o
}

// SetLocationProjectID adds the locationProjectId to the lock params
func (o *LockParams) SetLocationProjectID(locationProjectID string) {
	o.LocationProjectID = locationProjectID
}

// WriteToRequest writes these params to a swagger request
func (o *LockParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param cluster_id
	if err := r.SetPathParam("cluster_id", o.ClusterID); err != nil {
		return err
	}

	// path param location.organization_id
	if err := r.SetPathParam("location.organization_id", o.LocationOrganizationID); err != nil {
		return err
	}

	// path param location.project_id
	if err := r.SetPathParam("location.project_id", o.LocationProjectID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
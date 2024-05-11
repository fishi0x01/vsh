// Code generated by go-swagger; DO NOT EDIT.

package project_service

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
	"github.com/go-openapi/swag"
)

// NewProjectServiceListParams creates a new ProjectServiceListParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewProjectServiceListParams() *ProjectServiceListParams {
	return &ProjectServiceListParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewProjectServiceListParamsWithTimeout creates a new ProjectServiceListParams object
// with the ability to set a timeout on a request.
func NewProjectServiceListParamsWithTimeout(timeout time.Duration) *ProjectServiceListParams {
	return &ProjectServiceListParams{
		timeout: timeout,
	}
}

// NewProjectServiceListParamsWithContext creates a new ProjectServiceListParams object
// with the ability to set a context for a request.
func NewProjectServiceListParamsWithContext(ctx context.Context) *ProjectServiceListParams {
	return &ProjectServiceListParams{
		Context: ctx,
	}
}

// NewProjectServiceListParamsWithHTTPClient creates a new ProjectServiceListParams object
// with the ability to set a custom HTTPClient for a request.
func NewProjectServiceListParamsWithHTTPClient(client *http.Client) *ProjectServiceListParams {
	return &ProjectServiceListParams{
		HTTPClient: client,
	}
}

/*
ProjectServiceListParams contains all the parameters to send to the API endpoint

	for the project service list operation.

	Typically these are written to a http.Request.
*/
type ProjectServiceListParams struct {

	/* PaginationNextPageToken.

	     Specifies a page token to use to retrieve the next page. Set this to the
	`next_page_token` returned by previous list requests to get the next page of
	results. If set, `previous_page_token` must not be set.
	*/
	PaginationNextPageToken *string

	/* PaginationPageSize.

	     The max number of results per page that should be returned. If the number
	of available results is larger than `page_size`, a `next_page_token` is
	returned which can be used to get the next page of results in subsequent
	requests. A value of zero will cause `page_size` to be defaulted.

	     Format: int64
	*/
	PaginationPageSize *int64

	/* PaginationPreviousPageToken.

	     Specifies a page token to use to retrieve the previous page. Set this to
	the `previous_page_token` returned by previous list requests to get the
	previous page of results. If set, `next_page_token` must not be set.
	*/
	PaginationPreviousPageToken *string

	/* ScopeID.

	   id is the id of the object being referenced.
	*/
	ScopeID *string

	/* ScopeType.

	   ResourceType is the type of object being referenced.

	   Default: "UNKNOWN"
	*/
	ScopeType *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the project service list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ProjectServiceListParams) WithDefaults() *ProjectServiceListParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the project service list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ProjectServiceListParams) SetDefaults() {
	var (
		scopeTypeDefault = string("UNKNOWN")
	)

	val := ProjectServiceListParams{
		ScopeType: &scopeTypeDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the project service list params
func (o *ProjectServiceListParams) WithTimeout(timeout time.Duration) *ProjectServiceListParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the project service list params
func (o *ProjectServiceListParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the project service list params
func (o *ProjectServiceListParams) WithContext(ctx context.Context) *ProjectServiceListParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the project service list params
func (o *ProjectServiceListParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the project service list params
func (o *ProjectServiceListParams) WithHTTPClient(client *http.Client) *ProjectServiceListParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the project service list params
func (o *ProjectServiceListParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPaginationNextPageToken adds the paginationNextPageToken to the project service list params
func (o *ProjectServiceListParams) WithPaginationNextPageToken(paginationNextPageToken *string) *ProjectServiceListParams {
	o.SetPaginationNextPageToken(paginationNextPageToken)
	return o
}

// SetPaginationNextPageToken adds the paginationNextPageToken to the project service list params
func (o *ProjectServiceListParams) SetPaginationNextPageToken(paginationNextPageToken *string) {
	o.PaginationNextPageToken = paginationNextPageToken
}

// WithPaginationPageSize adds the paginationPageSize to the project service list params
func (o *ProjectServiceListParams) WithPaginationPageSize(paginationPageSize *int64) *ProjectServiceListParams {
	o.SetPaginationPageSize(paginationPageSize)
	return o
}

// SetPaginationPageSize adds the paginationPageSize to the project service list params
func (o *ProjectServiceListParams) SetPaginationPageSize(paginationPageSize *int64) {
	o.PaginationPageSize = paginationPageSize
}

// WithPaginationPreviousPageToken adds the paginationPreviousPageToken to the project service list params
func (o *ProjectServiceListParams) WithPaginationPreviousPageToken(paginationPreviousPageToken *string) *ProjectServiceListParams {
	o.SetPaginationPreviousPageToken(paginationPreviousPageToken)
	return o
}

// SetPaginationPreviousPageToken adds the paginationPreviousPageToken to the project service list params
func (o *ProjectServiceListParams) SetPaginationPreviousPageToken(paginationPreviousPageToken *string) {
	o.PaginationPreviousPageToken = paginationPreviousPageToken
}

// WithScopeID adds the scopeID to the project service list params
func (o *ProjectServiceListParams) WithScopeID(scopeID *string) *ProjectServiceListParams {
	o.SetScopeID(scopeID)
	return o
}

// SetScopeID adds the scopeId to the project service list params
func (o *ProjectServiceListParams) SetScopeID(scopeID *string) {
	o.ScopeID = scopeID
}

// WithScopeType adds the scopeType to the project service list params
func (o *ProjectServiceListParams) WithScopeType(scopeType *string) *ProjectServiceListParams {
	o.SetScopeType(scopeType)
	return o
}

// SetScopeType adds the scopeType to the project service list params
func (o *ProjectServiceListParams) SetScopeType(scopeType *string) {
	o.ScopeType = scopeType
}

// WriteToRequest writes these params to a swagger request
func (o *ProjectServiceListParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.PaginationNextPageToken != nil {

		// query param pagination.next_page_token
		var qrPaginationNextPageToken string

		if o.PaginationNextPageToken != nil {
			qrPaginationNextPageToken = *o.PaginationNextPageToken
		}
		qPaginationNextPageToken := qrPaginationNextPageToken
		if qPaginationNextPageToken != "" {

			if err := r.SetQueryParam("pagination.next_page_token", qPaginationNextPageToken); err != nil {
				return err
			}
		}
	}

	if o.PaginationPageSize != nil {

		// query param pagination.page_size
		var qrPaginationPageSize int64

		if o.PaginationPageSize != nil {
			qrPaginationPageSize = *o.PaginationPageSize
		}
		qPaginationPageSize := swag.FormatInt64(qrPaginationPageSize)
		if qPaginationPageSize != "" {

			if err := r.SetQueryParam("pagination.page_size", qPaginationPageSize); err != nil {
				return err
			}
		}
	}

	if o.PaginationPreviousPageToken != nil {

		// query param pagination.previous_page_token
		var qrPaginationPreviousPageToken string

		if o.PaginationPreviousPageToken != nil {
			qrPaginationPreviousPageToken = *o.PaginationPreviousPageToken
		}
		qPaginationPreviousPageToken := qrPaginationPreviousPageToken
		if qPaginationPreviousPageToken != "" {

			if err := r.SetQueryParam("pagination.previous_page_token", qPaginationPreviousPageToken); err != nil {
				return err
			}
		}
	}

	if o.ScopeID != nil {

		// query param scope.id
		var qrScopeID string

		if o.ScopeID != nil {
			qrScopeID = *o.ScopeID
		}
		qScopeID := qrScopeID
		if qScopeID != "" {

			if err := r.SetQueryParam("scope.id", qScopeID); err != nil {
				return err
			}
		}
	}

	if o.ScopeType != nil {

		// query param scope.type
		var qrScopeType string

		if o.ScopeType != nil {
			qrScopeType = *o.ScopeType
		}
		qScopeType := qrScopeType
		if qScopeType != "" {

			if err := r.SetQueryParam("scope.type", qScopeType); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

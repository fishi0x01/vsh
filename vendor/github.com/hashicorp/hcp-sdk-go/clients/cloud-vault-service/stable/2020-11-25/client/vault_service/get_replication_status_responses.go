// Code generated by go-swagger; DO NOT EDIT.

package vault_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	cloud "github.com/hashicorp/hcp-sdk-go/clients/cloud-shared/v1/models"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/models"
)

// GetReplicationStatusReader is a Reader for the GetReplicationStatus structure.
type GetReplicationStatusReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetReplicationStatusReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetReplicationStatusOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetReplicationStatusDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetReplicationStatusOK creates a GetReplicationStatusOK with default headers values
func NewGetReplicationStatusOK() *GetReplicationStatusOK {
	return &GetReplicationStatusOK{}
}

/*
GetReplicationStatusOK describes a response with status code 200, with default header values.

A successful response.
*/
type GetReplicationStatusOK struct {
	Payload *models.HashicorpCloudVault20201125GetReplicationStatusResponse
}

// IsSuccess returns true when this get replication status o k response has a 2xx status code
func (o *GetReplicationStatusOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get replication status o k response has a 3xx status code
func (o *GetReplicationStatusOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get replication status o k response has a 4xx status code
func (o *GetReplicationStatusOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get replication status o k response has a 5xx status code
func (o *GetReplicationStatusOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get replication status o k response a status code equal to that given
func (o *GetReplicationStatusOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get replication status o k response
func (o *GetReplicationStatusOK) Code() int {
	return 200
}

func (o *GetReplicationStatusOK) Error() string {
	return fmt.Sprintf("[GET /vault/2020-11-25/organizations/{location.organization_id}/projects/{location.project_id}/clusters/{cluster_id}/replication-status][%d] getReplicationStatusOK  %+v", 200, o.Payload)
}

func (o *GetReplicationStatusOK) String() string {
	return fmt.Sprintf("[GET /vault/2020-11-25/organizations/{location.organization_id}/projects/{location.project_id}/clusters/{cluster_id}/replication-status][%d] getReplicationStatusOK  %+v", 200, o.Payload)
}

func (o *GetReplicationStatusOK) GetPayload() *models.HashicorpCloudVault20201125GetReplicationStatusResponse {
	return o.Payload
}

func (o *GetReplicationStatusOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpCloudVault20201125GetReplicationStatusResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetReplicationStatusDefault creates a GetReplicationStatusDefault with default headers values
func NewGetReplicationStatusDefault(code int) *GetReplicationStatusDefault {
	return &GetReplicationStatusDefault{
		_statusCode: code,
	}
}

/*
GetReplicationStatusDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type GetReplicationStatusDefault struct {
	_statusCode int

	Payload *cloud.GrpcGatewayRuntimeError
}

// IsSuccess returns true when this get replication status default response has a 2xx status code
func (o *GetReplicationStatusDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get replication status default response has a 3xx status code
func (o *GetReplicationStatusDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get replication status default response has a 4xx status code
func (o *GetReplicationStatusDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get replication status default response has a 5xx status code
func (o *GetReplicationStatusDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get replication status default response a status code equal to that given
func (o *GetReplicationStatusDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get replication status default response
func (o *GetReplicationStatusDefault) Code() int {
	return o._statusCode
}

func (o *GetReplicationStatusDefault) Error() string {
	return fmt.Sprintf("[GET /vault/2020-11-25/organizations/{location.organization_id}/projects/{location.project_id}/clusters/{cluster_id}/replication-status][%d] GetReplicationStatus default  %+v", o._statusCode, o.Payload)
}

func (o *GetReplicationStatusDefault) String() string {
	return fmt.Sprintf("[GET /vault/2020-11-25/organizations/{location.organization_id}/projects/{location.project_id}/clusters/{cluster_id}/replication-status][%d] GetReplicationStatus default  %+v", o._statusCode, o.Payload)
}

func (o *GetReplicationStatusDefault) GetPayload() *cloud.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *GetReplicationStatusDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(cloud.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

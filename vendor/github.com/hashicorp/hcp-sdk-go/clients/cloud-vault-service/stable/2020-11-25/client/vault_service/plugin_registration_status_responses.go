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

// PluginRegistrationStatusReader is a Reader for the PluginRegistrationStatus structure.
type PluginRegistrationStatusReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PluginRegistrationStatusReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPluginRegistrationStatusOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewPluginRegistrationStatusDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPluginRegistrationStatusOK creates a PluginRegistrationStatusOK with default headers values
func NewPluginRegistrationStatusOK() *PluginRegistrationStatusOK {
	return &PluginRegistrationStatusOK{}
}

/*
PluginRegistrationStatusOK describes a response with status code 200, with default header values.

A successful response.
*/
type PluginRegistrationStatusOK struct {
	Payload *models.HashicorpCloudVault20201125PluginRegistrationStatusResponse
}

// IsSuccess returns true when this plugin registration status o k response has a 2xx status code
func (o *PluginRegistrationStatusOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this plugin registration status o k response has a 3xx status code
func (o *PluginRegistrationStatusOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this plugin registration status o k response has a 4xx status code
func (o *PluginRegistrationStatusOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this plugin registration status o k response has a 5xx status code
func (o *PluginRegistrationStatusOK) IsServerError() bool {
	return false
}

// IsCode returns true when this plugin registration status o k response a status code equal to that given
func (o *PluginRegistrationStatusOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the plugin registration status o k response
func (o *PluginRegistrationStatusOK) Code() int {
	return 200
}

func (o *PluginRegistrationStatusOK) Error() string {
	return fmt.Sprintf("[GET /vault/2020-11-25/organizations/{location.organization_id}/projects/{location.project_id}/clusters/{cluster_id}/plugin/registration-status][%d] pluginRegistrationStatusOK  %+v", 200, o.Payload)
}

func (o *PluginRegistrationStatusOK) String() string {
	return fmt.Sprintf("[GET /vault/2020-11-25/organizations/{location.organization_id}/projects/{location.project_id}/clusters/{cluster_id}/plugin/registration-status][%d] pluginRegistrationStatusOK  %+v", 200, o.Payload)
}

func (o *PluginRegistrationStatusOK) GetPayload() *models.HashicorpCloudVault20201125PluginRegistrationStatusResponse {
	return o.Payload
}

func (o *PluginRegistrationStatusOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpCloudVault20201125PluginRegistrationStatusResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPluginRegistrationStatusDefault creates a PluginRegistrationStatusDefault with default headers values
func NewPluginRegistrationStatusDefault(code int) *PluginRegistrationStatusDefault {
	return &PluginRegistrationStatusDefault{
		_statusCode: code,
	}
}

/*
PluginRegistrationStatusDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type PluginRegistrationStatusDefault struct {
	_statusCode int

	Payload *cloud.GrpcGatewayRuntimeError
}

// IsSuccess returns true when this plugin registration status default response has a 2xx status code
func (o *PluginRegistrationStatusDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this plugin registration status default response has a 3xx status code
func (o *PluginRegistrationStatusDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this plugin registration status default response has a 4xx status code
func (o *PluginRegistrationStatusDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this plugin registration status default response has a 5xx status code
func (o *PluginRegistrationStatusDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this plugin registration status default response a status code equal to that given
func (o *PluginRegistrationStatusDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the plugin registration status default response
func (o *PluginRegistrationStatusDefault) Code() int {
	return o._statusCode
}

func (o *PluginRegistrationStatusDefault) Error() string {
	return fmt.Sprintf("[GET /vault/2020-11-25/organizations/{location.organization_id}/projects/{location.project_id}/clusters/{cluster_id}/plugin/registration-status][%d] PluginRegistrationStatus default  %+v", o._statusCode, o.Payload)
}

func (o *PluginRegistrationStatusDefault) String() string {
	return fmt.Sprintf("[GET /vault/2020-11-25/organizations/{location.organization_id}/projects/{location.project_id}/clusters/{cluster_id}/plugin/registration-status][%d] PluginRegistrationStatus default  %+v", o._statusCode, o.Payload)
}

func (o *PluginRegistrationStatusDefault) GetPayload() *cloud.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *PluginRegistrationStatusDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(cloud.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

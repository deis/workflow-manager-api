package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/deis/workflow-manager-api/pkg/swagger/models"
)

/*GetComponentByLatestReleaseOK component latest release response

swagger:response getComponentByLatestReleaseOK
*/
type GetComponentByLatestReleaseOK struct {

	// In: body
	Payload *models.ComponentVersion `json:"body,omitempty"`
}

// NewGetComponentByLatestReleaseOK creates GetComponentByLatestReleaseOK with default headers values
func NewGetComponentByLatestReleaseOK() *GetComponentByLatestReleaseOK {
	return &GetComponentByLatestReleaseOK{}
}

// WithPayload adds the payload to the get component by latest release o k response
func (o *GetComponentByLatestReleaseOK) WithPayload(payload *models.ComponentVersion) *GetComponentByLatestReleaseOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get component by latest release o k response
func (o *GetComponentByLatestReleaseOK) SetPayload(payload *models.ComponentVersion) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetComponentByLatestReleaseOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetComponentByLatestReleaseDefault unexpected error

swagger:response getComponentByLatestReleaseDefault
*/
type GetComponentByLatestReleaseDefault struct {
	_statusCode int

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetComponentByLatestReleaseDefault creates GetComponentByLatestReleaseDefault with default headers values
func NewGetComponentByLatestReleaseDefault(code int) *GetComponentByLatestReleaseDefault {
	if code <= 0 {
		code = 500
	}

	return &GetComponentByLatestReleaseDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get component by latest release default response
func (o *GetComponentByLatestReleaseDefault) WithStatusCode(code int) *GetComponentByLatestReleaseDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get component by latest release default response
func (o *GetComponentByLatestReleaseDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get component by latest release default response
func (o *GetComponentByLatestReleaseDefault) WithPayload(payload *models.Error) *GetComponentByLatestReleaseDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get component by latest release default response
func (o *GetComponentByLatestReleaseDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetComponentByLatestReleaseDefault) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/deis/workflow-manager-api/pkg/swagger/models"
)

/*GetClusterByIDOK clusters details response

swagger:response getClusterByIdOK
*/
type GetClusterByIDOK struct {

	// In: body
	Payload *models.Cluster `json:"body,omitempty"`
}

// NewGetClusterByIDOK creates GetClusterByIDOK with default headers values
func NewGetClusterByIDOK() *GetClusterByIDOK {
	return &GetClusterByIDOK{}
}

// WithPayload adds the payload to the get cluster by id o k response
func (o *GetClusterByIDOK) WithPayload(payload *models.Cluster) *GetClusterByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cluster by id o k response
func (o *GetClusterByIDOK) SetPayload(payload *models.Cluster) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetClusterByIDOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetClusterByIDDefault unexpected error

swagger:response getClusterByIdDefault
*/
type GetClusterByIDDefault struct {
	_statusCode int

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetClusterByIDDefault creates GetClusterByIDDefault with default headers values
func NewGetClusterByIDDefault(code int) *GetClusterByIDDefault {
	if code <= 0 {
		code = 500
	}

	return &GetClusterByIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get cluster by id default response
func (o *GetClusterByIDDefault) WithStatusCode(code int) *GetClusterByIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get cluster by id default response
func (o *GetClusterByIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get cluster by id default response
func (o *GetClusterByIDDefault) WithPayload(payload *models.Error) *GetClusterByIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cluster by id default response
func (o *GetClusterByIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetClusterByIDDefault) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
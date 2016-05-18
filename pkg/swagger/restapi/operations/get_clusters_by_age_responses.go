package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/deis/workflow-manager-api/pkg/swagger/models"
)

/*GetClustersByAgeOK clusters details response

swagger:response getClustersByAgeOK
*/
type GetClustersByAgeOK struct {

	// In: body
	Payload GetClustersByAgeOKBodyBody `json:"body,omitempty"`
}

// NewGetClustersByAgeOK creates GetClustersByAgeOK with default headers values
func NewGetClustersByAgeOK() *GetClustersByAgeOK {
	return &GetClustersByAgeOK{}
}

// WithPayload adds the payload to the get clusters by age o k response
func (o *GetClustersByAgeOK) WithPayload(payload GetClustersByAgeOKBodyBody) *GetClustersByAgeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get clusters by age o k response
func (o *GetClustersByAgeOK) SetPayload(payload GetClustersByAgeOKBodyBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetClustersByAgeOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(200)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*GetClustersByAgeDefault unexpected error

swagger:response getClustersByAgeDefault
*/
type GetClustersByAgeDefault struct {
	_statusCode int

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetClustersByAgeDefault creates GetClustersByAgeDefault with default headers values
func NewGetClustersByAgeDefault(code int) *GetClustersByAgeDefault {
	if code <= 0 {
		code = 500
	}

	return &GetClustersByAgeDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get clusters by age default response
func (o *GetClustersByAgeDefault) WithStatusCode(code int) *GetClustersByAgeDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get clusters by age default response
func (o *GetClustersByAgeDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get clusters by age default response
func (o *GetClustersByAgeDefault) WithPayload(payload *models.Error) *GetClustersByAgeDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get clusters by age default response
func (o *GetClustersByAgeDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetClustersByAgeDefault) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
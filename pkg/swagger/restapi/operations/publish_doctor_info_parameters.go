package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	strfmt "github.com/go-swagger/go-swagger/strfmt"

	"github.com/deis/workflow-manager-api/pkg/swagger/models"
)

// NewPublishDoctorInfoParams creates a new PublishDoctorInfoParams object
// with the default values initialized.
func NewPublishDoctorInfoParams() PublishDoctorInfoParams {
	var ()
	return PublishDoctorInfoParams{}
}

// PublishDoctorInfoParams contains all the bound params for the publish doctor info operation
// typically these are obtained from a http.Request
//
// swagger:parameters publishDoctorInfo
type PublishDoctorInfoParams struct {
	/*
	  In: body
	*/
	Body *models.DoctorInfo
	/*A universal Id to represent a sepcific request or report
	  Required: true
	  In: path
	*/
	UUID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *PublishDoctorInfoParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	defer r.Body.Close()
	var body models.DoctorInfo
	if err := route.Consumer.Consume(r.Body, &body); err != nil {
		res = append(res, errors.NewParseError("body", "body", "", err))
	} else {
		if err := body.Validate(route.Formats); err != nil {
			res = append(res, err)
		}

		if len(res) == 0 {
			o.Body = &body
		}
	}

	rUUID, rhkUUID, _ := route.Params.GetOK("uuid")
	if err := o.bindUUID(rUUID, rhkUUID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PublishDoctorInfoParams) bindUUID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.UUID = raw

	return nil
}
package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/go-swagger/go-swagger/errors"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/httpkit/validate"
	"github.com/go-swagger/go-swagger/strfmt"
	"github.com/go-swagger/go-swagger/swag"
)

// GetComponentsByLatestReleaseForV2HandlerFunc turns a function with the right signature into a get components by latest release for v2 handler
type GetComponentsByLatestReleaseForV2HandlerFunc func(GetComponentsByLatestReleaseForV2Params) middleware.Responder

// Handle executing the request and returning a response
func (fn GetComponentsByLatestReleaseForV2HandlerFunc) Handle(params GetComponentsByLatestReleaseForV2Params) middleware.Responder {
	return fn(params)
}

// GetComponentsByLatestReleaseForV2Handler interface for that can handle valid get components by latest release for v2 params
type GetComponentsByLatestReleaseForV2Handler interface {
	Handle(GetComponentsByLatestReleaseForV2Params) middleware.Responder
}

// NewGetComponentsByLatestReleaseForV2 creates a new http.Handler for the get components by latest release for v2 operation
func NewGetComponentsByLatestReleaseForV2(ctx *middleware.Context, handler GetComponentsByLatestReleaseForV2Handler) *GetComponentsByLatestReleaseForV2 {
	return &GetComponentsByLatestReleaseForV2{Context: ctx, Handler: handler}
}

/*GetComponentsByLatestReleaseForV2 swagger:route POST /v2/versions/latest getComponentsByLatestReleaseForV2

list the latest release version of the components.This endpoint is to support old clients

*/
type GetComponentsByLatestReleaseForV2 struct {
	Context *middleware.Context
	Handler GetComponentsByLatestReleaseForV2Handler
}

func (o *GetComponentsByLatestReleaseForV2) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetComponentsByLatestReleaseForV2Params()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

/*GetComponentsByLatestReleaseForV2Body get components by latest release for v2 body

swagger:model GetComponentsByLatestReleaseForV2Body
*/
type GetComponentsByLatestReleaseForV2Body struct {

	/* data
	 */
	Data []*models.ComponentVersion `json:"data,omitempty"`
}

// Validate validates this get components by latest release for v2 body
func (o *GetComponentsByLatestReleaseForV2Body) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetComponentsByLatestReleaseForV2Body) validateData(formats strfmt.Registry) error {

	if swag.IsZero(o.Data) { // not required
		return nil
	}

	for i := 0; i < len(o.Data); i++ {

		if o.Data[i] != nil {

			if err := o.Data[i].Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

/*GetComponentsByLatestReleaseForV2OKBodyBody get components by latest release for v2 o k body body

swagger:model GetComponentsByLatestReleaseForV2OKBodyBody
*/
type GetComponentsByLatestReleaseForV2OKBodyBody struct {

	/* data

	Required: true
	*/
	Data []*models.ComponentVersion `json:"data"`
}

// Validate validates this get components by latest release for v2 o k body body
func (o *GetComponentsByLatestReleaseForV2OKBodyBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetComponentsByLatestReleaseForV2OKBodyBody) validateData(formats strfmt.Registry) error {

	if err := validate.Required("getComponentsByLatestReleaseForV2OK"+"."+"data", "body", o.Data); err != nil {
		return err
	}

	for i := 0; i < len(o.Data); i++ {

		if o.Data[i] != nil {

			if err := o.Data[i].Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
)

// CreateClusterDetailsForV2HandlerFunc turns a function with the right signature into a create cluster details for v2 handler
type CreateClusterDetailsForV2HandlerFunc func(CreateClusterDetailsForV2Params) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateClusterDetailsForV2HandlerFunc) Handle(params CreateClusterDetailsForV2Params) middleware.Responder {
	return fn(params)
}

// CreateClusterDetailsForV2Handler interface for that can handle valid create cluster details for v2 params
type CreateClusterDetailsForV2Handler interface {
	Handle(CreateClusterDetailsForV2Params) middleware.Responder
}

// NewCreateClusterDetailsForV2 creates a new http.Handler for the create cluster details for v2 operation
func NewCreateClusterDetailsForV2(ctx *middleware.Context, handler CreateClusterDetailsForV2Handler) *CreateClusterDetailsForV2 {
	return &CreateClusterDetailsForV2{Context: ctx, Handler: handler}
}

/*CreateClusterDetailsForV2 swagger:route POST /v2/clusters/{id} createClusterDetailsForV2

create a cluster with all components.This endpoint is to support old clients

*/
type CreateClusterDetailsForV2 struct {
	Context *middleware.Context
	Handler CreateClusterDetailsForV2Handler
}

func (o *CreateClusterDetailsForV2) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewCreateClusterDetailsForV2Params()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
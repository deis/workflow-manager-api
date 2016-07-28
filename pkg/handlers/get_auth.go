package handlers

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

const authStaticHTML = "./static_html/auth.html"

type getAuthOK struct {
}

// newGetAuthOK creates GetAuthOK with default headers values
func newGetAuthOK() *getAuthOK {
	return &getAuthOK{}
}

// WriteResponse to the client
func (o *getAuthOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {
	f, err := ioutil.ReadFile(authStaticHTML)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.WriteHeader(200)
	io.WriteString(rw, string(f))
}

// GetAuth is the handler for the auth endpoint
func GetAuth() middleware.Responder {
	return newGetAuthOK()
}

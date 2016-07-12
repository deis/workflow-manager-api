package restapi

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-swagger/go-swagger/swag"
	graceful "github.com/tylerb/graceful"

	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
)

const (
	DefaultRDSType
)

//go:generate swagger generate server -t ../.. -A WorkflowManager -f ./swagger.yml

// NewServer creates a new api workflow manager server but does not configure it
func NewServer(api *operations.WorkflowManagerAPI) *Server {
	s := new(Server)
	s.api = api
	return s
}

// ConfigureAPI configures the API and handlers. Needs to be called before Serve
func (s *Server) ConfigureAPI() {
	if s.api != nil {
		s.handler = configureAPI(s.api)
	}
}

// ConfigureFlags configures the additional flags defined by the handlers. Needs to be called before the parser.Parse
func (s *Server) ConfigureFlags() {
	if s.api != nil {
		configureFlags(s.api)
	}
}

// Server for the workflow manager API
type Server struct {
	Host         string `long:"host" description:"the IP to listen on" default:"localhost" env:"HOST"`
	Port         int    `long:"port" description:"the port to listen on for insecure connections, defaults to a random value" env:"PORT"`
	RDSRegion    string `long:"rds-region" description:"The region of the RDS postgres server"`
	PostgresType string `long:"postgres-type" required description:"The type of the postgres server (either 'incluster' or 'rds')"`
	httpServerL  net.Listener

	api          *operations.WorkflowManagerAPI
	handler      http.Handler
	hasListeners bool
}

// SetAPI configures the server with the specified API. Needs to be called before Serve
func (s *Server) SetAPI(api *operations.WorkflowManagerAPI) {
	if api == nil {
		s.api = nil
		s.handler = nil
		return
	}

	s.api = api
	s.handler = configureAPI(api)
}

// Serve the api
func (s *Server) Serve() (err error) {
	if !s.hasListeners {
		if err := s.Listen(); err != nil {
			return err
		}
	}

	httpServer := &graceful.Server{Server: new(http.Server)}
	httpServer.Handler = s.handler

	fmt.Printf("serving workflow manager at http://%s\n", s.httpServerL.Addr())
	l := s.httpServerL
	if err := httpServer.Serve(tcpKeepAliveListener{l.(*net.TCPListener)}); err != nil {
		return err
	}

	return nil
}

// Listen creates the listeners for the server
func (s *Server) Listen() error {
	if s.hasListeners { // already done this
		return nil
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}

	h, p, err := swag.SplitHostPort(listener.Addr().String())
	if err != nil {
		return err
	}
	s.Host = h
	s.Port = p
	s.httpServerL = listener

	s.hasListeners = true
	return nil
}

// Shutdown server and clean up resources
func (s *Server) Shutdown() error {
	s.api.ServerShutdown()
	return nil
}

// tcpKeepAliveListener is copied from the stdlib net/http package

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

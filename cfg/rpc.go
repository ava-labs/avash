/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package cfg

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

// AvashRPC is the RPCService handler for the avash service
var AvashRPC *RPCService

// RPCService is for maintaining a reference to the root JSON RPC server and the HTTP router
type RPCService struct {
	RPCServer  *rpc.Server
	HTTPRouter *mux.Router
	urlpath    string
	host       string
	port       string
	endpoints  map[string]interface{}
}

type ServiceOption func(*RPCService) error

func WithURLPath(urlpath string) ServiceOption {
	return func(s *RPCService) error {
		s.urlpath = urlpath
		return nil
	}
}

func WithHost(host string) ServiceOption {
	return func(s *RPCService) error {
		if host != "" {
			s.host = host
		}
		return nil
	}
}

func WithPort(port int) ServiceOption {
	return func(s *RPCService) error {
		// empty - return an error ?
		if port > 0 {
			s.port = fmt.Sprintf("%d", port)
		}

		return nil
	}
}

func WithRPCServer(server *rpc.Server) ServiceOption {
	return func(s *RPCService) error {
		if server == nil {
			return fmt.Errorf("server can't be nil")
		}

		s.RPCServer = server
		return nil
	}
}

func WithHTTPRouter(httpRouter *mux.Router) ServiceOption {
	return func(s *RPCService) error {
		if httpRouter == nil {
			return fmt.Errorf("router can't be nil")
		}

		s.HTTPRouter = httpRouter
		return nil
	}
}

// functions helps for unit testing.
var rpcServer = func() *rpc.Server { return rpc.NewServer() }
var httpRouter = func() *mux.Router { return mux.NewRouter() }

// NewRPCService create new RPC serice instance.
func NewRPCService(opts ...ServiceOption) (*RPCService, error) {
	s := &RPCService{
		urlpath:    "/rpc",
		host:       "localhost",
		port:       "9020",
		RPCServer:  rpcServer(),
		HTTPRouter: httpRouter(),
	}

	for _, o := range opts {
		if err := o(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (rs *RPCService) Init() {
	rs.RPCServer.RegisterCodec(json.NewCodec(), "application/json")
	rs.RPCServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	rs.RPCServer.RegisterCodec(json.NewCodec(), "application/json")
	rs.RPCServer.RegisterService(rs, "")

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%s", rs.host, rs.port), rs.HTTPRouter); err != nil {
			// logger.
		}
	}()
}

// RegisterServer registers the adds the rpc and http servers to the plugins service
func (rpcsrv *RPCService) RegisterServer(s *rpc.Server, r *mux.Router) {
	rpcsrv.RPCServer = s
	rpcsrv.HTTPRouter = r
}

// Initialize creates the RPC server at the provided baseurl, hostname, and port
func (rpcsrv *RPCService) Initialize(urlpath string, host string, port string) {
	rpcsrv.urlpath = urlpath
	rpcsrv.host = host
	rpcsrv.port = port
	if rpcsrv.urlpath == "" {
		rpcsrv.urlpath = "/rpc"
	}
	if host == "" {
		rpcsrv.host = "localhost"
	}
	if port == "" {
		rpcsrv.port = "9020"
	}
	s := rpc.NewServer()
	r := mux.NewRouter()
	rpcsrv.RegisterServer(s, r)
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	s.RegisterService(rpcsrv, "")
	r.Handle(rpcsrv.urlpath, s)
	go http.ListenAndServe(rpcsrv.host+":"+rpcsrv.port, r)
}

// AddService registers the appropriate endpoint for the plugin given a symbol
func (rpcsrv *RPCService) AddService(serviceInstance interface{}, endpoint string) {
	rpcsrv.RPCServer.RegisterService(serviceInstance, "")
	rpcsrv.HTTPRouter.Handle(rpcsrv.urlpath+"/"+endpoint, rpcsrv.RPCServer)
}

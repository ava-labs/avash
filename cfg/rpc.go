// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

package cfg

import (
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

package cfg

import (
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
)

func TestNewRPCService(t *testing.T) {
	testServer := rpc.NewServer()
	testRouter := mux.NewRouter()

	customServer := rpc.NewServer()
	customRouter := mux.NewRouter()

	tests := []struct {
		name        string
		input       []ServiceOption
		want        *RPCService
		mockServers bool
		wantErr     bool
	}{{
		name:        "default values case",
		input:       nil,
		mockServers: true,
		want: &RPCService{
			RPCServer:  testServer,
			HTTPRouter: testRouter,
			host:       "localhost",
			port:       "9020",
			urlpath:    "/rpc",
		},
	}, {
		name: "url, host and port case",
		input: []ServiceOption{
			WithHost("my-host"),
			WithPort(9000),
			WithURLPath("/test/path"),
		},
		mockServers: true,
		want: &RPCService{
			RPCServer:  testServer,
			HTTPRouter: testRouter,
			host:       "my-host",
			port:       "9000",
			urlpath:    "/test/path",
		},
	}, {
		name: "custom router and server case",
		input: []ServiceOption{
			WithHTTPRouter(customRouter),
			WithRPCServer(customServer),
		},
		want: &RPCService{
			RPCServer:  customServer,
			HTTPRouter: customRouter,
			host:       "localhost",
			port:       "9020",
			urlpath:    "/rpc",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockServers {

				originalRPCServer := rpcServer
				originalRouter := httpRouter

				rpcServer = func() *rpc.Server { return testServer }
				httpRouter = func() *mux.Router { return testRouter }

				defer func() {
					rpcServer = originalRPCServer
					httpRouter = originalRouter
				}()
			}

			got, err := NewRPCService(tt.input...)
			if err != nil {
				t.Fatalf("NewRPCService(%v) failed: %v", tt.input, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("NewRPCService(%v) failed: got %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/ybbus/jsonrpc"
)

func TestNewKeyCmdRunE(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{{
		name:    "success case",
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newKeyCmdRunE(&cobra.Command{}, []string{})
			if tt.wantErr != (got != nil) {
				t.Fatalf("newKeyCmdRunE() failed: got %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

type fakeRPCClient struct {
	Err         error
	RpcResponse *jsonrpc.RPCResponse
}

func (f fakeRPCClient) Call(method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	return f.RpcResponse, f.Err
}

func (f fakeRPCClient) CallRaw(request *jsonrpc.RPCRequest) (*jsonrpc.RPCResponse, error) {
	return f.RpcResponse, f.Err
}

func (f fakeRPCClient) CallFor(out interface{}, method string, params ...interface{}) error {
	return f.Err
}

func (f fakeRPCClient) CallBatch(requests jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
	return nil, f.Err
}

func (f fakeRPCClient) CallBatchRaw(requests jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
	return nil, f.Err
}

func TestSendCmdRunE(t *testing.T) {
	tests := []struct {
		name      string
		metadata  func(name string) (string, error)
		rpcClient func(host string, port string) jsonrpc.RPCClient
		wantErr   bool
	}{{
		name: "success case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{RpcResponse: &jsonrpc.RPCResponse{}}
		},
		wantErr: false,
	}, {
		name: "metadata error case",
		metadata: func(name string) (string, error) {
			return "", fmt.Errorf("not found")
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{RpcResponse: &jsonrpc.RPCResponse{}}
		},
		wantErr: true,
	}, {
		name: "rpc call error case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{
				Err: fmt.Errorf("fake error"),
			}
		},
		wantErr: true,
	}, {
		name: "issueTx error response case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{
				RpcResponse: &jsonrpc.RPCResponse{
					Error: &jsonrpc.RPCError{
						Code:    1,
						Message: "fake error",
					},
				},
			}
		},
		wantErr: true,
	}, {
		name: "issueTx response unmarshal error case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{
				RpcResponse: &jsonrpc.RPCResponse{
					Result: `{"TxID":what?}`,
				},
			}
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldMetadata := metadata
			oldrpcClient := avmRPCClient

			metadata = tt.metadata
			avmRPCClient = tt.rpcClient
			defer func() {
				metadata = oldMetadata
				avmRPCClient = oldrpcClient
			}()

			got := sendCmdRunE(&cobra.Command{}, []string{"test", "testingid"})
			if tt.wantErr != (got != nil) {
				t.Fatalf("sendCmdRunE() failed: got %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestStatusCmdRunE(t *testing.T) {
	tests := []struct {
		name      string
		metadata  func(name string) (string, error)
		rpcClient func(host string, port string) jsonrpc.RPCClient
		wantErr   bool
	}{{
		name: "success case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{RpcResponse: &jsonrpc.RPCResponse{}}
		},
		wantErr: false,
	}, {
		name: "metadata error case",
		metadata: func(name string) (string, error) {
			return "", fmt.Errorf("not found")
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{RpcResponse: &jsonrpc.RPCResponse{}}
		},
		wantErr: true,
	}, {
		name: "rpc call error case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{
				Err: fmt.Errorf("fake error"),
			}
		},
		wantErr: true,
	}, {
		name: "issueTx error response case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{
				RpcResponse: &jsonrpc.RPCResponse{
					Error: &jsonrpc.RPCError{
						Code:    1,
						Message: "fake error",
					},
				},
			}
		},
		wantErr: true,
	}, {
		name: "status response unmarshal error case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{
				RpcResponse: &jsonrpc.RPCResponse{
					Result: `{"TxID":what?}`,
				},
			}
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldMetadata := metadata
			oldrpcClient := avmRPCClient

			metadata = tt.metadata
			avmRPCClient = tt.rpcClient
			defer func() {
				metadata = oldMetadata
				avmRPCClient = oldrpcClient
			}()

			got := statusCmdRunE(&cobra.Command{}, []string{"test", "testingid"})
			if tt.wantErr != (got != nil) {
				t.Fatalf("sendCmdRunE() failed: got %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestGetBalanceCmdRunE(t *testing.T) {
	tests := []struct {
		name      string
		metadata  func(name string) (string, error)
		rpcClient func(host string, port string) jsonrpc.RPCClient
		wantErr   bool
	}{{
		name: "success case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{RpcResponse: &jsonrpc.RPCResponse{}}
		},
		wantErr: false,
	}, {
		name: "metadata error case",
		metadata: func(name string) (string, error) {
			return "", fmt.Errorf("not found")
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{RpcResponse: &jsonrpc.RPCResponse{}}
		},
		wantErr: true,
	}, {
		name: "rpc call error case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{
				Err: fmt.Errorf("fake error"),
			}
		},
		wantErr: true,
	}, {
		name: "issueTx error response case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{
				RpcResponse: &jsonrpc.RPCResponse{
					Error: &jsonrpc.RPCError{
						Code:    1,
						Message: "fake error",
					},
				},
			}
		},
		wantErr: true,
	}, {
		name: "balance response unmarshal error case",
		metadata: func(name string) (string, error) {
			return `{"public-ip": "testip", "host-port": "8080"}`, nil
		},
		rpcClient: func(host string, port string) jsonrpc.RPCClient {
			return &fakeRPCClient{
				RpcResponse: &jsonrpc.RPCResponse{
					Result: `{"TxID":what?}`,
				},
			}
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldMetadata := metadata
			oldrpcClient := avmRPCClient

			metadata = tt.metadata
			avmRPCClient = tt.rpcClient
			defer func() {
				metadata = oldMetadata
				avmRPCClient = oldrpcClient
			}()

			got := statusCmdRunE(&cobra.Command{}, []string{"test", "testingid"})
			if tt.wantErr != (got != nil) {
				t.Fatalf("sendCmdRunE() failed: got %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

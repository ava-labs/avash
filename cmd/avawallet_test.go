package cmd

import (
	"fmt"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/vms/avm"
	"github.com/ava-labs/avash/node"
	"github.com/spf13/cobra"
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

type fakeAVMClient struct {
	Err             error
	ID              [32]byte
	Status          choices.Status
	GetBalanceReply *avm.GetBalanceReply
}

func (f *fakeAVMClient) IssueTx(txBytes []byte) (ids.ID, error) {
	return f.ID, f.Err
}
func (f *fakeAVMClient) GetTxStatus(txID ids.ID) (choices.Status, error) {
	return f.Status, f.Err
}

func (f *fakeAVMClient) GetBalance(addr string, assetID string, includePartial bool) (*avm.GetBalanceReply, error) {
	if f.Err != nil {
		return nil, f.Err
	}

	return f.GetBalanceReply, nil
}

func TestSendCmdRunE(t *testing.T) {
	tests := []struct {
		name      string
		metadata  func(name string) (*node.Metadata, error)
		avmClient func(host, port string, requestTimeout time.Duration) Client
		wantErr   bool
	}{{
		name: "success case",
		metadata: func(name string) (*node.Metadata, error) {
			return &node.Metadata{Serverhost: "testip", HTTPport: "8080"}, nil
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{ID: ids.Empty}
		},
		wantErr: false,
	}, {
		name: "metadata error case",
		metadata: func(name string) (*node.Metadata, error) {
			return nil, fmt.Errorf("not found")
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{}
		},
		wantErr: true,
	}, {
		name: "avm client error case",
		metadata: func(name string) (*node.Metadata, error) {
			return &node.Metadata{Serverhost: "testip", HTTPport: "8080"}, nil
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{Err: fmt.Errorf("error")}
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldMetadata := metadata
			oldAVMClient := avmClient

			metadata = tt.metadata
			avmClient = tt.avmClient
			defer func() {
				metadata = oldMetadata
				avmClient = oldAVMClient
			}()

			got := sendCmdRunE(&cobra.Command{}, []string{"test", "testingid"})
			if tt.wantErr != (got != nil) {
				t.Fatalf("sendCmdRunE() failed: got %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestStatusCmdRunE(t *testing.T) {
	validIDStr := ids.ID{'a', 'v', 'a', ' ', 'l', 'a', 'b', 's'}.String()
	tests := []struct {
		name      string
		id        string
		metadata  func(name string) (*node.Metadata, error)
		avmClient func(host, port string, requestTimeout time.Duration) Client
		wantErr   bool
	}{{
		name: "success case",
		id:   validIDStr,
		metadata: func(name string) (*node.Metadata, error) {
			return &node.Metadata{Serverhost: "testip", HTTPport: "8080"}, nil
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{ID: [32]byte{0}}
		},
		wantErr: false,
	}, {
		name: "invalid ID case",
		id:   "incorrect",
		metadata: func(name string) (*node.Metadata, error) {
			return &node.Metadata{Serverhost: "testip", HTTPport: "8080"}, nil
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{ID: [32]byte{0}}
		},
		wantErr: true,
	}, {
		name: "metadata error case",
		id:   validIDStr,
		metadata: func(name string) (*node.Metadata, error) {
			return nil, fmt.Errorf("not found")
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{ID: ids.Empty}
		},
		wantErr: true,
	}, {
		name: "avm client error case",
		id:   validIDStr,
		metadata: func(name string) (*node.Metadata, error) {
			return &node.Metadata{Serverhost: "testip", HTTPport: "8080"}, nil
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{Err: fmt.Errorf("error")}
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldMetadata := metadata
			oldAVMClient := avmClient

			metadata = tt.metadata
			avmClient = tt.avmClient
			defer func() {
				metadata = oldMetadata
				avmClient = oldAVMClient
			}()
			got := statusCmdRunE(&cobra.Command{}, []string{"test", tt.id})
			if tt.wantErr != (got != nil) {
				t.Fatalf("sendCmdRunE() failed: got %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestGetBalanceCmdRunE(t *testing.T) {
	tests := []struct {
		name      string
		metadata  func(name string) (*node.Metadata, error)
		avmClient func(host, port string, requestTimeout time.Duration) Client
		wantErr   bool
	}{{
		name: "success case",
		metadata: func(name string) (*node.Metadata, error) {
			return &node.Metadata{Serverhost: "testip", HTTPport: "8080"}, nil
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{GetBalanceReply: &avm.GetBalanceReply{}}
		},
		wantErr: false,
	}, {
		name: "metadata error case",
		metadata: func(name string) (*node.Metadata, error) {
			return nil, fmt.Errorf("not found")
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{GetBalanceReply: &avm.GetBalanceReply{}}
		},
		wantErr: true,
	}, {
		name: "avm client error case",
		metadata: func(name string) (*node.Metadata, error) {
			return &node.Metadata{Serverhost: "testip", HTTPport: "8080"}, nil
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{Err: fmt.Errorf("error")}
		},
		wantErr: true,
	}, {
		name: "nil GetBalanceReply case",
		metadata: func(name string) (*node.Metadata, error) {
			return &node.Metadata{Serverhost: "testip", HTTPport: "8080"}, nil
		},
		avmClient: func(host, port string, requestTimeout time.Duration) Client {
			return &fakeAVMClient{}
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldMetadata := metadata
			oldAVMClient := avmClient

			metadata = tt.metadata
			avmClient = tt.avmClient
			defer func() {
				metadata = oldMetadata
				avmClient = oldAVMClient
			}()

			got := getBalanceCmdRunE(&cobra.Command{}, []string{"test", "testingid"})
			if tt.wantErr != (got != nil) {
				t.Fatalf("getBalanceCmdRunE() failed: got %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

package client

import (
	"github.com/zgfzgf/mid-lotus/api"
	"github.com/zgfzgf/mid-lotus/lib/jsonrpc"
)

// NewRPC creates a new http jsonrpc client.
func NewRPC(addr string) api.API {
	var res api.Struct
	jsonrpc.NewClient(addr, "Filecoin", &res.Internal)
	return &res
}

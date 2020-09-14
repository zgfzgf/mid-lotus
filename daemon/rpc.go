package daemon

import (
	"net/http"

	"github.com/zgfzgf/mid-lotus/api"
	"github.com/zgfzgf/mid-lotus/lib/jsonrpc"
)

func serveRPC(api api.API, addr string) error {
	rpcServer := jsonrpc.NewServer()
	rpcServer.Register("Filecoin", api)
	http.Handle("/rpc/v0", rpcServer)
	return http.ListenAndServe(addr, http.DefaultServeMux)
}

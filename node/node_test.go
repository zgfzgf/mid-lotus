package node

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/zgfzgf/mid-lotus/api"
	"github.com/zgfzgf/mid-lotus/api/client"
	"github.com/zgfzgf/mid-lotus/api/test"
	"github.com/zgfzgf/mid-lotus/lib/jsonrpc"

	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
)

func builder(t *testing.T, n int) []api.API {
	ctx := context.Background()
	mn := mocknet.New(ctx)

	out := make([]api.API, n)

	for i := 0; i < n; i++ {
		var err error
		out[i], err = New(ctx,
			Online(),
			MockHost(mn),
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := mn.LinkAll(); err != nil {
		t.Fatal(err)
	}

	return out
}

func TestAPI(t *testing.T) {
	test.TestApis(t, builder)
}

var nextApi int

func rpcBuilder(t *testing.T, n int) []api.API {
	nodeApis := builder(t, n)
	out := make([]api.API, n)

	for i, a := range nodeApis {
		rpcServer := jsonrpc.NewServer()
		rpcServer.Register("Filecoin", a)
		testServ := httptest.NewServer(rpcServer) //  todo: close
		out[i] = client.NewRPC(testServ.URL)
	}
	return out
}

func TestAPIRPC(t *testing.T) {
	test.TestApis(t, rpcBuilder)
}

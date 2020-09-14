package test

import (
	"context"
	"strings"
	"testing"

	"github.com/zgfzgf/mid-lotus/api"
	"github.com/zgfzgf/mid-lotus/build"
)

// APIBuilder is a function which is invoked in test suite to provide
// test nodes and networks
type APIBuilder func(t *testing.T, n int) []api.API
type testSuite struct {
	makeNodes APIBuilder
}

// TestApis is the entry point to API test suite
func TestApis(t *testing.T, b APIBuilder) {
	ts := testSuite{
		makeNodes: b,
	}

	t.Run("version", ts.testVersion)
	t.Run("id", ts.testID)
	t.Run("testConnectTwo", ts.testConnectTwo)

}

func (ts *testSuite) testVersion(t *testing.T) {
	ctx := context.Background()
	api := ts.makeNodes(t, 1)[0]

	v, err := api.Version(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v.Version != build.Version {
		t.Error("Version didn't work properly")
	}
}

func (ts *testSuite) testID(t *testing.T) {
	ctx := context.Background()
	api := ts.makeNodes(t, 1)[0]

	id, err := api.ID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(id.Pretty(), "Qm") {
		t.Error("expected identity to be Qm..")
	}
}

func (ts *testSuite) testConnectTwo(t *testing.T) {
	ctx := context.Background()
	apis := ts.makeNodes(t, 2)

	p, err := apis[0].NetPeers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 0 {
		t.Error("Node 0 has a peer")
	}

	p, err = apis[1].NetPeers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 0 {
		t.Error("Node 1 has a peer")
	}

	addrs, err := apis[1].NetAddrsListen(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err := apis[0].NetConnect(ctx, addrs); err != nil {
		t.Fatal(err)
	}

	p, err = apis[0].NetPeers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 1 {
		t.Error("Node 0 doesn't have 1 peer")
	}

	p, err = apis[1].NetPeers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 1 {
		t.Error("Node 0 doesn't have 1 peer")
	}
}

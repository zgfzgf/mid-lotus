package api

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Version provides various build-time information
type Version struct {
	Version string
}

// API is a low-level interface to the Filecoin network
type API interface {
	// network

	NetPeers(context.Context) ([]peer.AddrInfo, error) // TODO: check serialization
	NetConnect(context.Context, peer.AddrInfo) error
	NetAddrsListen(context.Context) (peer.AddrInfo, error)

	// ID returns peerID of libp2p node backing this API
	ID(context.Context) (peer.ID, error)

	// Version provides information about API provider
	Version(context.Context) (Version, error)
}

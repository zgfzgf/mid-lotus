package testing

import (
	blockstore "github.com/ipfs/go-ipfs-blockstore"

	"github.com/zgfzgf/mid-lotus/chain"
	"github.com/zgfzgf/mid-lotus/node/modules"
)

func MakeGenesis(bs blockstore.Blockstore, w *chain.Wallet) (modules.Genesis, error) {
	genb, err := chain.MakeGenesisBlock(bs, w)
	if err != nil {
		return nil, err
	}
	return genb.Genesis, nil
}

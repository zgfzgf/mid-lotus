package node

import (
	"errors"

	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/zgfzgf/mid-lotus/node/modules/lp2p"
)

func MockHost(mn mocknet.Mocknet) Option {
	return Options(
		applyIf(func(s *settings) bool { return !s.online },
			Error(errors.New("MockHost must be specified after Online")),
		),

		Override(new(lp2p.RawHost), lp2p.MockHost),
		Override(new(mocknet.Mocknet), mn),
	)
}

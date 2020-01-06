package proxy

import (
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	"github.com/tendermint/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Tendermint interface {
	ValidateBlockMeta(meta *tmtypes.BlockMeta, sh tmtypes.SignedHeader) error
	ValidateBlock(meta *tmtypes.Block, sh tmtypes.SignedHeader) error
}

type TendermintLiteWrapper struct {
}

func NewTendermintLiteWrapper() *TendermintLiteWrapper {
	wrapper := &TendermintLiteWrapper{}
	return wrapper
}

func (cw *TendermintLiteWrapper) ValidateBlockMeta(meta *types.BlockMeta, sh types.SignedHeader) error {
	return tmliteProxy.ValidateBlockMeta(meta, sh)
}

func (cw *TendermintLiteWrapper) ValidateBlock(meta *types.Block, sh types.SignedHeader) error {
	return tmliteProxy.ValidateBlock(meta, sh)
}

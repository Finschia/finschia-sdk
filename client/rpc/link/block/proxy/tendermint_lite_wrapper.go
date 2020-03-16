package proxy

import (
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
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

func (cw *TendermintLiteWrapper) ValidateBlockMeta(meta *tmtypes.BlockMeta, sh tmtypes.SignedHeader) error {
	return tmliteProxy.ValidateBlockMeta(meta, sh)
}

func (cw *TendermintLiteWrapper) ValidateBlock(meta *tmtypes.Block, sh tmtypes.SignedHeader) error {
	return tmliteProxy.ValidateBlock(meta, sh)
}

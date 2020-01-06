package context

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

type CLIContext interface {
	Verify(height int64) (tmtypes.SignedHeader, error)
	TrustNode() bool
	Indent() bool
	Codec() *codec.Codec
	GetNode() (rpcclient.Client, error)
	CosmosCliCtx() context.CLIContext
}

type Wrapper struct {
	cosCliCtx context.CLIContext
}

func NewCLIContextWrapper(cCliCtx context.CLIContext) *Wrapper {
	wrapper := &Wrapper{cosCliCtx: cCliCtx}
	return wrapper
}

func (ctx Wrapper) Verify(height int64) (tmtypes.SignedHeader, error) {
	return ctx.cosCliCtx.Verify(height)
}

func (ctx Wrapper) TrustNode() bool {
	return ctx.cosCliCtx.TrustNode
}

func (ctx Wrapper) Indent() bool {
	return ctx.cosCliCtx.Indent
}

func (ctx Wrapper) Codec() *codec.Codec {
	return ctx.cosCliCtx.Codec
}

func (ctx Wrapper) GetNode() (rpcclient.Client, error) {
	return ctx.cosCliCtx.GetNode()
}

func (ctx Wrapper) CosmosCliCtx() context.CLIContext {
	return ctx.cosCliCtx
}

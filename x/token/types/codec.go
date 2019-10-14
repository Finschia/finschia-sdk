package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgPublishToken{}, "link/MsgPublishToken", nil)
	cdc.RegisterConcrete(MsgMint{}, "link/MsgMint", nil)
	cdc.RegisterConcrete(MsgBurn{}, "link/MsgBurn", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

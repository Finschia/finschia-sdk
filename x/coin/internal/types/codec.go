package types

import (
	"github.com/line/lbm-sdk/codec"
)

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "coin/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "coin/MsgMultiSend", nil)
}

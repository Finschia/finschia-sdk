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

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateAccount{}, "account/MsgCreateAccount", nil)
	cdc.RegisterConcrete(MsgEmpty{}, "account/MsgEmpty", nil)
}

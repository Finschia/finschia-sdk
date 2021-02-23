package types

import "github.com/line/lbm-sdk/codec"

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	ModuleCdc.Seal()
}

package legacytx

import (
	"github.com/line/lfb-sdk/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(StdTx{}, "lfb-sdk/StdTx", nil)
}

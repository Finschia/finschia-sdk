package types

import (
	"github.com/Finschia/finschia-rdk/codec"
	cryptocodec "github.com/Finschia/finschia-rdk/crypto/codec"
	sdk "github.com/Finschia/finschia-rdk/types"
)

var amino = codec.NewLegacyAmino()

func init() {
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}

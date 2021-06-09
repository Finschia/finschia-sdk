package types

import (
	"github.com/line/lfb-sdk/codec"
	cryptocodec "github.com/line/lfb-sdk/crypto/codec"
)

var (
	amino = codec.NewLegacyAmino()
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

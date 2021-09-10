package types

import (
	"github.com/line/lbm-sdk/codec"
	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
)

var (
	amino = codec.NewLegacyAmino()
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

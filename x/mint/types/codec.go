package types

import (
	"github.com/line/lbm-sdk/v2/codec"
	cryptocodec "github.com/line/lbm-sdk/v2/crypto/codec"
)

var (
	amino = codec.NewLegacyAmino()
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

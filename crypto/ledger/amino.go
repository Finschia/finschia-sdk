package ledger

import (
	"github.com/line/lbm-sdk/v2/codec"
	cryptoAmino "github.com/line/lbm-sdk/v2/crypto/codec"
)

var cdc = codec.NewLegacyAmino()

func init() {
	RegisterAmino(cdc)
	cryptoAmino.RegisterCrypto(cdc)
}

// RegisterAmino registers all go-crypto related types in the given (amino) codec.
func RegisterAmino(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(PrivKeyLedgerSecp256k1{},
		"ostracon/PrivKeyLedgerSecp256k1", nil)
}

package multisig

import (
	"github.com/line/ostracon/crypto/sr25519"

	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/v2/crypto/keys/secp256k1"
	cryptotypes "github.com/line/lbm-sdk/v2/crypto/types"
)

// TODO: Figure out API for others to either add their own pubkey types, or
// to make verify / marshal accept a AminoCdc.
const (
	PubKeyAminoRoute = "ostracon/PubKeyMultisigThreshold"
)

var AminoCdc = codec.NewLegacyAmino()

func init() {
	AminoCdc.RegisterInterface((*cryptotypes.PubKey)(nil), nil)
	AminoCdc.RegisterConcrete(ed25519.PubKey{},
		ed25519.PubKeyName, nil)
	AminoCdc.RegisterConcrete(sr25519.PubKey{},
		sr25519.PubKeyName, nil)
	AminoCdc.RegisterConcrete(&secp256k1.PubKey{},
		secp256k1.PubKeyName, nil)
	AminoCdc.RegisterConcrete(&LegacyAminoPubKey{},
		PubKeyAminoRoute, nil)
}

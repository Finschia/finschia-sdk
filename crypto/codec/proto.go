package codec

import (
	codectypes "github.com/line/lbm-sdk/v2/codec/types"
	"github.com/line/lbm-sdk/v2/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/v2/crypto/keys/multisig"
	"github.com/line/lbm-sdk/v2/crypto/keys/secp256k1"
	cryptotypes "github.com/line/lbm-sdk/v2/crypto/types"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface("lbm.crypto.PubKey", (*cryptotypes.PubKey)(nil))
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &ed25519.PubKey{})
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{})
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &multisig.LegacyAminoPubKey{})
}

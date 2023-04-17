package codec

import (
	tmprotocrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"

	occrypto "github.com/Finschia/ostracon/crypto"
	"github.com/Finschia/ostracon/crypto/encoding"

	"github.com/Finschia/finschia-sdk/crypto/keys/ed25519"
	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/Finschia/finschia-sdk/crypto/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// FromOcProtoPublicKey converts a OC's tmprotocrypto.PublicKey into our own PubKey.
func FromOcProtoPublicKey(protoPk tmprotocrypto.PublicKey) (cryptotypes.PubKey, error) {
	switch protoPk := protoPk.Sum.(type) {
	case *tmprotocrypto.PublicKey_Ed25519:
		return &ed25519.PubKey{
			Key: protoPk.Ed25519,
		}, nil
	case *tmprotocrypto.PublicKey_Secp256K1:
		return &secp256k1.PubKey{
			Key: protoPk.Secp256K1,
		}, nil
	default:
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "cannot convert %v from Ostracon public key", protoPk)
	}
}

// ToOcProtoPublicKey converts our own PubKey to OC's tmprotocrypto.PublicKey.
func ToOcProtoPublicKey(pk cryptotypes.PubKey) (tmprotocrypto.PublicKey, error) {
	switch pk := pk.(type) {
	case *ed25519.PubKey:
		return tmprotocrypto.PublicKey{
			Sum: &tmprotocrypto.PublicKey_Ed25519{
				Ed25519: pk.Key,
			},
		}, nil
	case *secp256k1.PubKey:
		return tmprotocrypto.PublicKey{
			Sum: &tmprotocrypto.PublicKey_Secp256K1{
				Secp256K1: pk.Key,
			},
		}, nil
	default:
		return tmprotocrypto.PublicKey{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "cannot convert %v to Ostracon public key", pk)
	}
}

// FromOcPubKeyInterface converts OC's occrypto.PubKey to our own PubKey.
func FromOcPubKeyInterface(tmPk occrypto.PubKey) (cryptotypes.PubKey, error) {
	ocProtoPk, err := encoding.PubKeyToProto(tmPk)
	if err != nil {
		return nil, err
	}

	return FromOcProtoPublicKey(ocProtoPk)
}

// ToOcPubKeyInterface converts our own PubKey to OC's occrypto.PubKey.
func ToOcPubKeyInterface(pk cryptotypes.PubKey) (occrypto.PubKey, error) {
	ocProtoPk, err := ToOcProtoPublicKey(pk)
	if err != nil {
		return nil, err
	}

	return encoding.PubKeyFromProto(&ocProtoPk)
}

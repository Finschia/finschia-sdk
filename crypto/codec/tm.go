package codec

import (
	ostcrypto "github.com/line/ostracon/crypto"
	"github.com/line/ostracon/crypto/encoding"
	ostprotocrypto "github.com/line/ostracon/proto/ostracon/crypto"

	"github.com/line/lfb-sdk/crypto/keys/ed25519"
	"github.com/line/lfb-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/line/lfb-sdk/crypto/types"
	sdkerrors "github.com/line/lfb-sdk/types/errors"
)

// FromTmProtoPublicKey converts a TM's ostprotocrypto.PublicKey into our own PubKey.
func FromTmProtoPublicKey(protoPk ostprotocrypto.PublicKey) (cryptotypes.PubKey, error) {
	switch protoPk := protoPk.Sum.(type) {
	case *ostprotocrypto.PublicKey_Ed25519:
		return &ed25519.PubKey{
			Key: protoPk.Ed25519,
		}, nil
	case *ostprotocrypto.PublicKey_Secp256K1:
		return &secp256k1.PubKey{
			Key: protoPk.Secp256K1,
		}, nil
	default:
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "cannot convert %v from Tendermint public key", protoPk)
	}
}

// ToTmProtoPublicKey converts our own PubKey to TM's ostprotocrypto.PublicKey.
func ToTmProtoPublicKey(pk cryptotypes.PubKey) (ostprotocrypto.PublicKey, error) {
	switch pk := pk.(type) {
	case *ed25519.PubKey:
		return ostprotocrypto.PublicKey{
			Sum: &ostprotocrypto.PublicKey_Ed25519{
				Ed25519: pk.Key,
			},
		}, nil
	case *secp256k1.PubKey:
		return ostprotocrypto.PublicKey{
			Sum: &ostprotocrypto.PublicKey_Secp256K1{
				Secp256K1: pk.Key,
			},
		}, nil
	default:
		return ostprotocrypto.PublicKey{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "cannot convert %v to Tendermint public key", pk)
	}
}

// FromTmPubKeyInterface converts TM's ostcrypto.PubKey to our own PubKey.
func FromTmPubKeyInterface(tmPk ostcrypto.PubKey) (cryptotypes.PubKey, error) {
	tmProtoPk, err := encoding.PubKeyToProto(tmPk)
	if err != nil {
		return nil, err
	}

	return FromTmProtoPublicKey(tmProtoPk)
}

// ToTmPubKeyInterface converts our own PubKey to TM's ostcrypto.PubKey.
func ToTmPubKeyInterface(pk cryptotypes.PubKey) (ostcrypto.PubKey, error) {
	tmProtoPk, err := ToTmProtoPublicKey(pk)
	if err != nil {
		return nil, err
	}

	return encoding.PubKeyFromProto(tmProtoPk)
}

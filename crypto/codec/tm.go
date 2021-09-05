package codec

import (
	occrypto "github.com/line/ostracon/crypto"
	"github.com/line/ostracon/crypto/encoding"
	ocprotocrypto "github.com/line/ostracon/proto/ostracon/crypto"

	"github.com/line/lfb-sdk/crypto/keys/ed25519"
	"github.com/line/lfb-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/line/lfb-sdk/crypto/types"
	sdkerrors "github.com/line/lfb-sdk/types/errors"
)

// FromTmProtoPublicKey converts a TM's ocprotocrypto.PublicKey into our own PubKey.
func FromTmProtoPublicKey(protoPk ocprotocrypto.PublicKey) (cryptotypes.PubKey, error) {
	switch protoPk := protoPk.Sum.(type) {
	case *ocprotocrypto.PublicKey_Ed25519:
		return &ed25519.PubKey{
			Key: protoPk.Ed25519,
		}, nil
	case *ocprotocrypto.PublicKey_Secp256K1:
		return &secp256k1.PubKey{
			Key: protoPk.Secp256K1,
		}, nil
	default:
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "cannot convert %v from Ostracon public key", protoPk)
	}
}

// ToTmProtoPublicKey converts our own PubKey to TM's ocprotocrypto.PublicKey.
func ToTmProtoPublicKey(pk cryptotypes.PubKey) (ocprotocrypto.PublicKey, error) {
	switch pk := pk.(type) {
	case *ed25519.PubKey:
		return ocprotocrypto.PublicKey{
			Sum: &ocprotocrypto.PublicKey_Ed25519{
				Ed25519: pk.Key,
			},
		}, nil
	case *secp256k1.PubKey:
		return ocprotocrypto.PublicKey{
			Sum: &ocprotocrypto.PublicKey_Secp256K1{
				Secp256K1: pk.Key,
			},
		}, nil
	default:
		return ocprotocrypto.PublicKey{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "cannot convert %v to Ostracon public key", pk)
	}
}

// FromTmPubKeyInterface converts TM's occrypto.PubKey to our own PubKey.
func FromTmPubKeyInterface(tmPk occrypto.PubKey) (cryptotypes.PubKey, error) {
	tmProtoPk, err := encoding.PubKeyToProto(tmPk)
	if err != nil {
		return nil, err
	}

	return FromTmProtoPublicKey(tmProtoPk)
}

// ToTmPubKeyInterface converts our own PubKey to TM's occrypto.PubKey.
func ToTmPubKeyInterface(pk cryptotypes.PubKey) (occrypto.PubKey, error) {
	tmProtoPk, err := ToTmProtoPublicKey(pk)
	if err != nil {
		return nil, err
	}

	return encoding.PubKeyFromProto(&tmProtoPk)
}

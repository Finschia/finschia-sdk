package codec

import (
	occrypto "github.com/line/ostracon/crypto"
	"github.com/line/ostracon/crypto/encoding"
	ocprotocrypto "github.com/line/ostracon/proto/ostracon/crypto"

	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/line/lbm-sdk/crypto/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

// FromOcProtoPublicKey converts a OC's ocprotocrypto.PublicKey into our own PubKey.
func FromOcProtoPublicKey(protoPk ocprotocrypto.PublicKey) (cryptotypes.PubKey, error) {
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

// ToOcProtoPublicKey converts our own PubKey to OC's ocprotocrypto.PublicKey.
func ToOcProtoPublicKey(pk cryptotypes.PubKey) (ocprotocrypto.PublicKey, error) {
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
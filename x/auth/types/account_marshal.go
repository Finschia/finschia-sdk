package types

import (
	"bytes"
	"fmt"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/codec"
)

var _, AccountPrefix = amino.NameToDisfix("cosmos-sdk/Account")

func (acc *BaseAccount) MarshalAminoBare(registered bool) ([]byte, error) {
	if acc == nil {
		return nil, nil
	}

	buf := bytes.NewBuffer(nil)

	if registered {
		if _, err := buf.Write(AccountPrefix[:]); err != nil {
			return nil, err
		}
	}

	if err := codec.EncodeFieldByteSlice(buf, 1, acc.Address); err != nil {
		return nil, err
	}

	if err := acc.Coins.MarshalFieldAmino(buf, 2); err != nil {
		return nil, err
	}

	if acc.PubKey != nil {
		if err := codec.EncodeFieldByteSlice(buf, 3, acc.PubKey.Bytes()); err != nil {
			return nil, err
		}
	}

	if err := codec.EncodeFieldUvarint(buf, 4, acc.AccountNumber); err != nil {
		return nil, err
	}

	if err := codec.EncodeFieldUvarint(buf, 5, acc.Sequence); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (acc *BaseAccount) UnmarshalAminoBare(bz []byte) (n int, err error) {
	var _n int
	acc.Address, _n, err = codec.DecodeFieldByteSlice(bz, 1)
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 || err != nil {
		return n, err
	}

	_n, err = acc.Coins.UnmarshalFieldAmino(bz, 2)
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 || err != nil {
		return n, err
	}

	var pubKeyBz []byte
	pubKeyBz, _n, err = codec.DecodeFieldByteSlice(bz, 3)

	if pubKeyBz != nil {
		_, err = acc.unmarshalPubKeyAminoPrefix(pubKeyBz)
		if err != nil {
			return n, err
		}
	}

	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 || err != nil {
		return n, err
	}

	acc.AccountNumber, _n, err = codec.DecodeFieldUvarint(bz, 4)
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 || err != nil {
		return n, err
	}

	acc.Sequence, _n, err = codec.DecodeFieldUvarint(bz, 5)
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 || err != nil {
		return n, err
	}

	return n, err
}

func (acc *BaseAccount) unmarshalPubKeyAminoPrefix(bz []byte) (n int, err error) {
	// disamb, hasDisamb, prefix, hasPrefix, _n, err := amino.DecodeDisambPrefixBytes(bz)
	_, _, prefix, hasPrefix, _n, err := amino.DecodeDisambPrefixBytes(bz)
	if codec.Slide(&bz, &n, _n) && err != nil {
		return n, err
	}
	if !hasPrefix {
		return n, fmt.Errorf("should have prefix")
	}

	var byteslice []byte
	byteslice, _n, err = amino.DecodeByteSlice(bz)
	if codec.Slide(&bz, &n, _n) && err != nil {
		return n, err
	}

	switch prefix {
	case secp256k1.PubKeyPrefix:
		{
			var pubKey secp256k1.PubKeySecp256k1
			copy(pubKey[:], byteslice)
			acc.PubKey = pubKey
		}
	case ed25519.PubKeyPrefix:
		{
			var pubKey ed25519.PubKeyEd25519
			copy(pubKey[:], byteslice)
			acc.PubKey = pubKey
		}
	case multisig.PubKeyPrefix:
		{
			// TODO custom mashaller
			ModuleCdc.MustUnmarshalBinaryBare(byteslice, acc.PubKey)
		}
	}

	return n, err
}

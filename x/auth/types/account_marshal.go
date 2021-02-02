package types

import (
	"bytes"
	"fmt"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (acc *BaseAccount) Marshal() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if _, err := buf.Write(accountPrefix[:]); err != nil {
		return nil, err
	}

	if err := codec.EncodeFieldByteSlice(buf, 1, acc.Address); err != nil {
		return nil, err
	}

	if err := acc.Coins.Marshal(buf, 2); err != nil {
		return nil, err
	}

	if err := codec.EncodeFieldByteSlice(buf, 3, acc.PubKey.Bytes()); err != nil {
		return nil, err
	}

	if err := codec.EncodeFieldUvarint(buf, 4, acc.AccountNumber); err != nil {
		return nil, err
	}

	if err := codec.EncodeFieldUvarint(buf, 5, acc.Sequence); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

//
// func (acc *BaseAccount) encodeCoins(w io.Writer) error {
// 	if len(acc.Coins) == 0 {
// 		return nil
// 	}
//
// 	for _, coin := range acc.Coins {
// 		if err := codec.EncodeFieldNumberAndTyp3(w, 2, amino.Typ3_ByteLength); err != nil {
// 			return err
// 		}
//
// 		bz, err := coin.Marshal()
// 		if err != nil {
// 			return err
// 		}
// 		if err := amino.EncodeByteSlice(w, bz); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (acc *BaseAccount) Unmarshal(bz []byte) (n int, err error) {
	var _n int
	_n, err = acc.decodeAddress(bz)
	if err != nil {
		return
	}

	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 {
		return
	}

	_n, err = acc.decodeCoins(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 {
		return
	}

	_n, err = acc.decodePubKey(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 {
		return
	}

	_n, err = acc.decodeAccountNumber(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 {
		return
	}

	_n, err = acc.decodeSequence(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)

	return
}

func (acc *BaseAccount) decodeAddress(bz []byte) (n int, err error) {
	_n, err := codec.CheckFieldNumberAndTyp3(bz, 1, amino.Typ3_ByteLength)
	if _n == 0 || err != nil {
		return _n, err
	}
	codec.Slide(&bz, &n, _n)

	acc.Address, _n, err = amino.DecodeByteSlice(bz)
	codec.Slide(&bz, &n, _n)

	return
}

func (acc *BaseAccount) decodeCoins(bz []byte) (n int, err error) {
	var _n int

	for {
		if len(bz) == 0 {
			break
		}
		_n, err = codec.CheckFieldNumberAndTyp3(bz, 2, amino.Typ3_ByteLength)
		if _n == 0 || err != nil {
			break
		}
		codec.Slide(&bz, &n, _n)

		// REFACTOR
		if acc.Coins == nil {
			acc.Coins = make(sdk.Coins, 0)
		}

		var u uint64
		u, _n, err = amino.DecodeUvarint(bz)
		codec.Slide(&bz, &n, _n)

		buf := bz[0:u]
		_n, err = acc.decodeCoin(buf)

		codec.Slide(&bz, &n, _n)
	}

	return
}

func (acc *BaseAccount) decodeCoin(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 1, amino.Typ3_ByteLength)
	codec.Slide(&bz, &n, _n)

	var denom, amt string
	denom, _n, err = amino.DecodeString(bz)
	codec.Slide(&bz, &n, _n)

	_n, err = codec.CheckFieldNumberAndTyp3(bz, 2, amino.Typ3_ByteLength)
	codec.Slide(&bz, &n, _n)

	amt, _n, err = amino.DecodeString(bz)
	codec.Slide(&bz, &n, _n)

	amount, _ := sdk.NewIntFromString(amt)

	// acc.Coins.Add(sdk.Coin{Denom: denom, Amount: amount})
	acc.Coins = append(acc.Coins, sdk.Coin{Denom: denom, Amount: amount})
	return n, nil
}

func (acc *BaseAccount) decodePubKey(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 3, amino.Typ3_ByteLength)
	if _n == 0 || err != nil {
		return _n, err
	}
	codec.Slide(&bz, &n, _n)

	var u uint64
	u, _n, err = amino.DecodeUvarint(bz)
	codec.Slide(&bz, &n, _n)

	buf := bz[0:u]
	_n, err = acc.decodePubKeyInterface(buf)
	codec.Slide(&bz, &n, _n)
	return
}

func (acc *BaseAccount) decodePubKeyInterface(bz []byte) (n int, err error) {
	// disamb, hasDisamb, prefix, hasPrefix, _n, err := amino.DecodeDisambPrefixBytes(bz)
	_, _, prefix, hasPrefix, _n, err := amino.DecodeDisambPrefixBytes(bz)
	if codec.Slide(&bz, &n, _n) && err != nil {
		return n, err
	}
	if !hasPrefix {
		return n, fmt.Errorf("should have prefix")
	}

	prefixBytes := prefix.Bytes()

	var byteslice []byte
	byteslice, _n, err = amino.DecodeByteSlice(bz)
	if codec.Slide(&bz, &n, _n) && err != nil {
		return n, err
	}

	switch {
	case secp256k1.PubKeyPrefix.EqualBytes(prefixBytes):
		{
			var pubKey secp256k1.PubKeySecp256k1
			copy(pubKey[:], byteslice)
			acc.PubKey = pubKey
		}
	case ed25519.PubKeyPrefix.EqualBytes(prefixBytes):
		{
			var pubKey ed25519.PubKeyEd25519
			copy(pubKey[:], byteslice)
			acc.PubKey = pubKey
		}
	case multisig.PubKeyPrefix.EqualBytes(prefixBytes):
		{
			// TODO custom mashaller
			ModuleCdc.MustUnmarshalBinaryBare(byteslice, acc.PubKey)
		}
	}

	return n, err
}

func (acc *BaseAccount) decodeAccountNumber(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 4, amino.Typ3_Varint)
	if _n == 0 || err != nil {
		return _n, err
	}
	codec.Slide(&bz, &n, _n)

	u, _n, err := amino.DecodeUvarint(bz)
	acc.AccountNumber = u
	codec.Slide(&bz, &n, _n)
	return
}

func (acc *BaseAccount) decodeSequence(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 5, amino.Typ3_Varint)
	if _n == 0 || err != nil {
		return _n, err
	}
	codec.Slide(&bz, &n, _n)

	u, _n, err := amino.DecodeUvarint(bz)
	acc.Sequence = u
	codec.Slide(&bz, &n, _n)
	return
}

package types

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

func TestBaseAddressPubKey(t *testing.T) {
	_, pub1, addr1 := KeyTestPubAddr()
	_, pub2, addr2 := KeyTestPubAddr()
	acc := NewBaseAccountWithAddress(addr1)

	// check the address (set) and pubkey (not set)
	require.EqualValues(t, addr1, acc.GetAddress())
	require.EqualValues(t, nil, acc.GetPubKey())

	// can't override address
	err := acc.SetAddress(addr2)
	require.NotNil(t, err)
	require.EqualValues(t, addr1, acc.GetAddress())

	// set the pubkey
	err = acc.SetPubKey(pub1)
	require.Nil(t, err)
	require.Equal(t, pub1, acc.GetPubKey())

	// can override pubkey
	err = acc.SetPubKey(pub2)
	require.Nil(t, err)
	require.Equal(t, pub2, acc.GetPubKey())

	//------------------------------------

	// can set address on empty account
	acc2 := BaseAccount{}
	err = acc2.SetAddress(addr2)
	require.Nil(t, err)
	require.EqualValues(t, addr2, acc2.GetAddress())
}

func TestBaseAccountCoins(t *testing.T) {
	_, _, addr := KeyTestPubAddr()
	acc := NewBaseAccountWithAddress(addr)

	someCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 246)}

	err := acc.SetCoins(someCoins)
	require.Nil(t, err)
	require.Equal(t, someCoins, acc.GetCoins())
}

func TestBaseAccountSequence(t *testing.T) {
	_, _, addr := KeyTestPubAddr()
	acc := NewBaseAccountWithAddress(addr)

	seq := uint64(7)

	err := acc.SetSequence(seq)
	require.Nil(t, err)
	require.Equal(t, seq, acc.GetSequence())
}

func TestBaseAccountMarshal(t *testing.T) {
	_, pub, addr := KeyTestPubAddr()
	acc := NewBaseAccountWithAddress(addr)

	someCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 246)}
	seq := uint64(7)

	// set everything on the account
	err := acc.SetPubKey(pub)
	require.Nil(t, err)
	err = acc.SetSequence(seq)
	require.Nil(t, err)
	err = acc.SetCoins(someCoins)
	require.Nil(t, err)

	// need a codec for marshaling
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)

	exp, err := cdc.MarshalBinaryBare(acc)
	require.Nil(t, err)

	b, err := acc.Marshal()
	require.Nil(t, err)
	require.Equal(t, exp, b)

	acc2 := BaseAccount{}
	err = cdc.UnmarshalBinaryBare(b, &acc2)
	require.Nil(t, err)
	require.Equal(t, acc, acc2)

	// error on bad bytes
	acc2 = BaseAccount{}
	err = cdc.UnmarshalBinaryBare(b[:len(b)/2], &acc2)
	require.NotNil(t, err)
}

func TestBaseAccountUnmarshal(t *testing.T) {
	_, pub, addr := KeyTestPubAddr()
	acc := NewBaseAccountWithAddress(addr)

	someCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 246), sdk.NewInt64Coin("btc", 1000)}
	seq := uint64(7)

	// set everything on the account
	err := acc.SetPubKey(pub)
	require.Nil(t, err)
	err = acc.SetSequence(seq)
	require.Nil(t, err)
	err = acc.SetCoins(someCoins)
	require.Nil(t, err)

	// need a codec for marshaling
	cdc := codec.New()
	codec.RegisterCrypto(cdc)

	b, err := cdc.MarshalBinaryBare(acc)
	require.Nil(t, err)

	acc2 := BaseAccount{}
	err = cdc.UnmarshalBinaryBare(b, &acc2)
	require.Nil(t, err)
	require.Equal(t, acc, acc2)

	acc3 := BaseAccount{}
	_, err = acc3.Unmarshal(b)
	require.Nil(t, err)
	require.Equal(t, acc, acc3)
}

func TestGenesisAccountValidate(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	baseAcc := NewBaseAccount(addr, nil, pubkey, 0, 0)
	tests := []struct {
		name   string
		acc    exported.GenesisAccount
		expErr error
	}{
		{
			"valid base account",
			baseAcc,
			nil,
		},
		{
			"invalid base valid account",
			NewBaseAccount(addr, sdk.NewCoins(), secp256k1.GenPrivKey().PubKey(), 0, 0),
			errors.New("pubkey and address pair is invalid"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.acc.Validate()
			require.Equal(t, tt.expErr, err)
		})
	}
}

func TestBaseAccountJSON(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	coins := sdk.NewCoins(sdk.NewInt64Coin("test", 5))
	baseAcc := NewBaseAccount(addr, coins, pubkey, 10, 50)

	bz, err := json.Marshal(baseAcc)
	require.NoError(t, err)

	bz1, err := baseAcc.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz1), string(bz))

	var a BaseAccount
	require.NoError(t, json.Unmarshal(bz, &a))
	require.Equal(t, baseAcc.String(), a.String())

	bz, err = ModuleCdc.MarshalJSON(baseAcc)
	require.NoError(t, err)

	var b BaseAccount
	require.NoError(t, ModuleCdc.UnmarshalJSON(bz, &b))
	require.Equal(t, baseAcc.String(), b.String())
}

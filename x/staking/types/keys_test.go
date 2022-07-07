package types_test

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/staking/types"
)

var (
	keysPK1   = ed25519.GenPrivKeyFromSecret([]byte{1}).PubKey()
	keysPK2   = ed25519.GenPrivKeyFromSecret([]byte{2}).PubKey()
	keysPK3   = ed25519.GenPrivKeyFromSecret([]byte{3}).PubKey()
	keysAddr1 = keysPK1.Address()
	keysAddr2 = keysPK2.Address()
	keysAddr3 = keysPK3.Address()
)

func TestGetValidatorPowerRank(t *testing.T) {
	valAddr1 := sdk.ValAddress(keysAddr1)
	val1 := newValidator(t, valAddr1, keysPK1)
	val1.Tokens = sdk.ZeroInt()
	val2, val3, val4 := val1, val1, val1
	val2.Tokens = sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)
	val3.Tokens = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	x := new(big.Int).Exp(big.NewInt(2), big.NewInt(40), big.NewInt(0))
	val4.Tokens = sdk.TokensFromConsensusPower(x.Int64(), sdk.DefaultPowerReduction)

	tests := []struct {
		validator types.Validator
		wantHex   string
	}{
		{val1, "2300000000000000003293969194899e93908f9a8dce89cf8b97859889858fc7898bc98ec8868c8b929992c89ec6888998cf8f8f8c9992c6cd9998c6"},
		{val2, "2300000000000000013293969194899e93908f9a8dce89cf8b97859889858fc7898bc98ec8868c8b929992c89ec6888998cf8f8f8c9992c6cd9998c6"},
		{val3, "23000000000000000a3293969194899e93908f9a8dce89cf8b97859889858fc7898bc98ec8868c8b929992c89ec6888998cf8f8f8c9992c6cd9998c6"},
		{val4, "2300000100000000003293969194899e93908f9a8dce89cf8b97859889858fc7898bc98ec8868c8b929992c89ec6888998cf8f8f8c9992c6cd9998c6"},
	}
	for i, tt := range tests {
		got := hex.EncodeToString(types.GetValidatorsByPowerIndexKey(tt.validator, sdk.DefaultPowerReduction))

		require.Equal(t, tt.wantHex, got, "Keys did not match on test case %d", i)
	}
}

func TestGetREDByValDstIndexKey(t *testing.T) {
	tests := []struct {
		delAddr    sdk.AccAddress
		valSrcAddr sdk.ValAddress
		valDstAddr sdk.ValAddress
		wantHex    string
	}{
		{sdk.AccAddress(keysAddr1), sdk.ValAddress(keysAddr1), sdk.ValAddress(keysAddr1),
			"36326c696e6b76616c6f70657231763074687a67767a703876743671377973746d666d37613977766730707073666d39326667392b6c696e6b31763074687a67767a703876743671377973746d666d376139777667307070736666336735786b326c696e6b76616c6f70657231763074687a67767a703876743671377973746d666d37613977766730707073666d3932666739"},
		{sdk.AccAddress(keysAddr1), sdk.ValAddress(keysAddr2), sdk.ValAddress(keysAddr3),
			"36326c696e6b76616c6f7065723138326d7a3772766e736a64376639307a72636c66717961397a757063373364613074633273342b6c696e6b31763074687a67767a703876743671377973746d666d376139777667307070736666336735786b326c696e6b76616c6f70657231746d656d74756a75326a3278366a35666c7378736e356833796573353273386a7a37736e7777"},
		{sdk.AccAddress(keysAddr2), sdk.ValAddress(keysAddr1), sdk.ValAddress(keysAddr3),
			"36326c696e6b76616c6f7065723138326d7a3772766e736a64376639307a72636c66717961397a757063373364613074633273342b6c696e6b31746d656d74756a75326a3278366a35666c7378736e356833796573353273386a73326a777161326c696e6b76616c6f70657231763074687a67767a703876743671377973746d666d37613977766730707073666d3932666739"},
	}
	for i, tt := range tests {
		got := hex.EncodeToString(types.GetREDByValDstIndexKey(tt.delAddr, tt.valSrcAddr, tt.valDstAddr))

		require.Equal(t, tt.wantHex, got, "Keys did not match on test case %d", i)
	}
}

func TestGetREDByValSrcIndexKey(t *testing.T) {
	tests := []struct {
		delAddr    sdk.AccAddress
		valSrcAddr sdk.ValAddress
		valDstAddr sdk.ValAddress
		wantHex    string
	}{
		{sdk.AccAddress(keysAddr1), sdk.ValAddress(keysAddr1), sdk.ValAddress(keysAddr1),
			"35326c696e6b76616c6f70657231763074687a67767a703876743671377973746d666d37613977766730707073666d39326667392b6c696e6b31763074687a67767a703876743671377973746d666d376139777667307070736666336735786b326c696e6b76616c6f70657231763074687a67767a703876743671377973746d666d37613977766730707073666d3932666739"},
		{sdk.AccAddress(keysAddr1), sdk.ValAddress(keysAddr2), sdk.ValAddress(keysAddr3),
			"35326c696e6b76616c6f70657231746d656d74756a75326a3278366a35666c7378736e356833796573353273386a7a37736e77772b6c696e6b31763074687a67767a703876743671377973746d666d376139777667307070736666336735786b326c696e6b76616c6f7065723138326d7a3772766e736a64376639307a72636c66717961397a75706337336461307463327334"},
		{sdk.AccAddress(keysAddr2), sdk.ValAddress(keysAddr1), sdk.ValAddress(keysAddr3),
			"35326c696e6b76616c6f70657231763074687a67767a703876743671377973746d666d37613977766730707073666d39326667392b6c696e6b31746d656d74756a75326a3278366a35666c7378736e356833796573353273386a73326a777161326c696e6b76616c6f7065723138326d7a3772766e736a64376639307a72636c66717961397a75706337336461307463327334"},
	}
	for i, tt := range tests {
		got := hex.EncodeToString(types.GetREDByValSrcIndexKey(tt.delAddr, tt.valSrcAddr, tt.valDstAddr))

		require.Equal(t, tt.wantHex, got, "Keys did not match on test case %d", i)
	}
}

func TestGetValidatorQueueKey(t *testing.T) {
	ts := time.Now()
	height := int64(1024)

	bz := types.GetValidatorQueueKey(ts, height)
	rTs, rHeight, err := types.ParseValidatorQueueKey(bz)
	require.NoError(t, err)
	require.Equal(t, ts.UTC(), rTs.UTC())
	require.Equal(t, rHeight, height)
}

func TestTestGetValidatorQueueKeyOrder(t *testing.T) {
	ts := time.Now().UTC()
	height := int64(1000)

	endKey := types.GetValidatorQueueKey(ts, height)

	keyA := types.GetValidatorQueueKey(ts.Add(-10*time.Minute), height-10)
	keyB := types.GetValidatorQueueKey(ts.Add(-5*time.Minute), height+50)
	keyC := types.GetValidatorQueueKey(ts.Add(10*time.Minute), height+100)

	require.Equal(t, -1, bytes.Compare(keyA, endKey)) // keyA <= endKey
	require.Equal(t, -1, bytes.Compare(keyB, endKey)) // keyB <= endKey
	require.Equal(t, 1, bytes.Compare(keyC, endKey))  // keyB >= endKey
}

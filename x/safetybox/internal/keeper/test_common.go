package keeper

import (
	"testing"

	"github.com/line/link/x/safetybox/internal/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"

	cbank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/line/link/x/contract"
	"github.com/line/link/x/iam"
	"github.com/line/link/x/token"
)

type TestInput struct {
	Cdc    *codec.Codec
	Ctx    sdk.Context
	Keeper Keeper
	Tk     token.Keeper
	Iam    iam.Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	cbank.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	iam.RegisterCodec(cdc)
	token.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func SetupTestInput(t *testing.T) TestInput {
	keyAuth := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyIam := sdk.NewKVStoreKey(iam.StoreKey)
	keyContract := sdk.NewKVStoreKey(contract.StoreKey)
	keyToken := sdk.NewKVStoreKey(token.StoreKey)
	keySafetyBox := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAuth, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIam, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyContract, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyToken, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySafetyBox, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.NoError(t, err)

	cdc := newTestCodec()

	// init params keeper and subspaces
	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams)
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	cbankSubspace := paramsKeeper.Subspace(cbank.DefaultParamspace)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[sdk.AccAddress([]byte("moduleAcc")).String()] = true

	// module account permissions
	maccPerms := map[string][]string{
		types.ModuleName: {supply.Burner, supply.Minter, supply.Staking},
	}

	// add keepers
	accountKeeper := auth.NewAccountKeeper(cdc, keyAuth, authSubspace, auth.ProtoBaseAccount)
	cbankKeeper := cbank.NewBaseKeeper(accountKeeper, cbankSubspace, blacklistedAddrs)
	iamKeeper := iam.NewKeeper(cdc, keyIam)
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, cbankKeeper, maccPerms)
	contractKeeper := contract.NewContractKeeper(cdc, keyContract)
	tokenKeeper := token.NewKeeper(cdc, accountKeeper, iamKeeper.WithPrefix(token.ModuleName), contractKeeper, keyToken)

	keeper := NewKeeper(cdc, iamKeeper.WithPrefix(types.ModuleName), tokenKeeper, keySafetyBox)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.NewCoins()))

	return TestInput{Cdc: cdc, Ctx: ctx, Keeper: keeper, Tk: tokenKeeper}
}

//nolint:unparam
func checkSafetyBoxBalance(t *testing.T, k Keeper, ctx sdk.Context, safetyBoxID string, ta, ca, ti int64) {
	sb, err := k.GetSafetyBox(ctx, safetyBoxID)
	require.NoError(t, err)

	var taExpected, caExpected, tiExpected sdk.Int
	if ta == 0 {
		taExpected = sdk.ZeroInt()
	} else {
		taExpected = sdk.NewInt(ta)
	}
	if ca == 0 {
		caExpected = sdk.ZeroInt()
	} else {
		caExpected = sdk.NewInt(ca)
	}
	if ti == 0 {
		tiExpected = sdk.ZeroInt()
	} else {
		tiExpected = sdk.NewInt(ti)
	}

	require.Equal(t, taExpected, sb.TotalAllocation)
	require.Equal(t, caExpected, sb.CumulativeAllocation)
	require.Equal(t, tiExpected, sb.TotalIssuance)
}

// testing events - the order of events only matter in the test
func VerifyEventFunc(t *testing.T, expected sdk.Events, actual sdk.Events) {
	require.Equal(
		t,
		sdk.StringifyEvents(expected.ToABCIEvents()).String(),
		sdk.StringifyEvents(actual.ToABCIEvents()).String(),
	)
}

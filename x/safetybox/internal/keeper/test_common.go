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

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/line/link/x/iam"
)

type TestInput struct {
	Cdc    *codec.Codec
	Ctx    sdk.Context
	Keeper Keeper
	Ak     auth.AccountKeeper
	Bk     bank.BaseKeeper
	Iam    iam.Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	iam.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func SetupTestInput(t *testing.T) TestInput {

	keyAuth := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyIam := sdk.NewKVStoreKey(iam.StoreKey)
	keySafetyBox := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAuth, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIam, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySafetyBox, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.NoError(t, err)

	cdc := newTestCodec()

	// init params keeper and subspaces
	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := paramsKeeper.Subspace(bank.DefaultParamspace)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[sdk.AccAddress([]byte("moduleAcc")).String()] = true

	// add keepers
	accountKeeper := auth.NewAccountKeeper(cdc, keyAuth, authSubspace, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, bankSubspace, bank.DefaultCodespace, blacklistedAddrs)
	iamKeeper := iam.NewKeeper(cdc, keyIam)

	// module account permissions
	maccPerms := map[string][]string{
		types.ModuleName: {supply.Burner, supply.Minter, supply.Staking},
	}

	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	keeper := NewKeeper(cdc, iamKeeper.WithPrefix(types.ModuleName), bankKeeper, accountKeeper, keySafetyBox)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.NewCoins()))

	return TestInput{Cdc: cdc, Ctx: ctx, Keeper: keeper, Ak: accountKeeper, Bk: bankKeeper}
}

func checkSafetyBoxBalance(t *testing.T, k Keeper, ctx sdk.Context, safetyBoxId string, ta, ca, ti int64) {
	sb, err := k.GetSafetyBox(ctx, safetyBoxId)
	require.NoError(t, err)

	// only accepts one type of coins
	require.LessOrEqual(t, len(sb.TotalAllocation), 1)
	require.LessOrEqual(t, len(sb.CumulativeAllocation), 1)
	require.LessOrEqual(t, len(sb.TotalIssuance), 1)

	var taExpected, caExpected, tiExpected sdk.Coins
	if ta == 0 {
		taExpected = sdk.Coins(nil)
	} else {
		taExpected = sdk.Coins{sdk.Coin{Denom: sb.Denoms[0], Amount: sdk.NewInt(ta)}}
	}
	if ca == 0 {
		caExpected = sdk.Coins(nil)
	} else {
		caExpected = sdk.Coins{sdk.Coin{Denom: sb.Denoms[0], Amount: sdk.NewInt(ca)}}
	}
	if ti == 0 {
		tiExpected = sdk.Coins(nil)
	} else {
		tiExpected = sdk.Coins{sdk.Coin{Denom: sb.Denoms[0], Amount: sdk.NewInt(ti)}}
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

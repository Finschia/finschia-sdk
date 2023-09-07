package testutil

import (
	"testing"

	"github.com/Finschia/ostracon/libs/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/Finschia/finschia-sdk/codec"
	codectypes "github.com/Finschia/finschia-sdk/codec/types"
	simappparams "github.com/Finschia/finschia-sdk/simapp/params"
	"github.com/Finschia/finschia-sdk/store"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
	"github.com/Finschia/finschia-sdk/x/or/da/keeper"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
	rutypes "github.com/Finschia/finschia-sdk/x/or/rollup/types"
)

func DaKeeper(t testing.TB, encCfg simappparams.EncodingConfig) (keeper.Keeper, sdk.Context, sdk.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	ctrl := gomock.NewController(t)
	accountKeeprMock := NewMockAccountKeeper(ctrl)
	accountKeeprMock.EXPECT().GetParams(gomock.Any()).Return(authtypes.DefaultParams()).AnyTimes()
	rollupKeeperMock := NewMockRollupKeeper(ctrl)
	rollupKeeperMock.EXPECT().GetAllRollup(gomock.Any()).Return([]rutypes.Rollup{
		{
			RollupName:     "rollup1",
			L1ToL2GasRatio: 100,
			PermissionedAddresses: rutypes.Sequencers{
				Addresses: []string{"addr1", "addr2"},
			},
			Creator: "creator1",
		},
	}).AnyTimes()
	k := keeper.NewKeeper(
		encCfg.TxConfig,
		cdc,
		storeKey,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		accountKeeprMock,
		rollupKeeperMock,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, ctx, storeKey
}

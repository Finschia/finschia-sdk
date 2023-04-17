package keeper

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/line/ostracon/libs/log"

	"github.com/Finschia/finschia-sdk/codec"
	codectypes "github.com/Finschia/finschia-sdk/codec/types"
	"github.com/Finschia/finschia-sdk/store"
	sdk "github.com/Finschia/finschia-sdk/types"
	accountkeeper "github.com/Finschia/finschia-sdk/x/auth/keeper"
	accounttypes "github.com/Finschia/finschia-sdk/x/auth/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	paramtypes "github.com/Finschia/finschia-sdk/x/params/types"
)

func genAddress() sdk.AccAddress {
	b := make([]byte, 20)
	rand.Read(b)
	return b
}

func setupKeeper(storeKey *sdk.KVStoreKey) BaseKeeper {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	amino := codec.NewLegacyAmino()
	testTransientStoreKey := sdk.NewTransientStoreKey("test")

	accountStoreKey := sdk.NewKVStoreKey(accounttypes.StoreKey)
	accountSubspace := paramtypes.NewSubspace(cdc, amino, accountStoreKey, testTransientStoreKey, accounttypes.ModuleName)
	accountKeeper := accountkeeper.NewAccountKeeper(cdc, accountStoreKey, accountSubspace, accounttypes.ProtoBaseAccount, nil)

	bankSubspace := paramtypes.NewSubspace(cdc, amino, storeKey, testTransientStoreKey, banktypes.StoreKey)
	return NewBaseKeeper(cdc, storeKey, accountKeeper, bankSubspace, nil, false)
}

func setupContext(t *testing.T, storeKey *sdk.KVStoreKey) sdk.Context {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	return sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
}

func TestInactiveAddr(t *testing.T) {
	storeKey := sdk.NewKVStoreKey(banktypes.StoreKey)
	bankKeeper := setupKeeper(storeKey)
	ctx := setupContext(t, storeKey)

	addr := genAddress()

	require.Equal(t, 0, len(bankKeeper.inactiveAddrs))

	bankKeeper.addToInactiveAddr(ctx, addr)
	require.True(t, bankKeeper.isStoredInactiveAddr(ctx, addr))

	// duplicate addition, no error expects.
	bankKeeper.addToInactiveAddr(ctx, addr)
	require.True(t, bankKeeper.isStoredInactiveAddr(ctx, addr))

	bankKeeper.deleteFromInactiveAddr(ctx, addr)
	require.False(t, bankKeeper.isStoredInactiveAddr(ctx, addr))

	addr2 := genAddress()
	require.False(t, bankKeeper.isStoredInactiveAddr(ctx, addr2))

	// expect no error
	bankKeeper.deleteFromInactiveAddr(ctx, addr2)

	// test loadAllInactiveAddrs
	bankKeeper.addToInactiveAddr(ctx, addr)
	bankKeeper.addToInactiveAddr(ctx, addr2)
	require.Equal(t, 0, len(bankKeeper.inactiveAddrs))
	bankKeeper.loadAllInactiveAddrs(ctx)
	require.Equal(t, 2, len(bankKeeper.inactiveAddrs))
}

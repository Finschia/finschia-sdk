package keeper

//import (
//	"math/rand"
//	"testing"
//
//	storetypes "cosmossdk.io/store/types"
//	"github.com/cosmos/cosmos-sdk/testutil"
//	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
//	"github.com/stretchr/testify/require"
//
//	"github.com/cosmos/cosmos-sdk/codec"
//	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
//	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
//)
//
//func genAddress() sdk.AccAddress {
//	b := make([]byte, 20)
//	rand.Read(b)
//	return b
//}
//
//func setupKeeper(storeKey *storetypes.KVStoreKey) BaseKeeper {
//	registry := codectypes.NewInterfaceRegistry()
//	cdc := codec.NewProtoCodec(registry)
//	amino := codec.NewLegacyAmino()
//	testTransientStoreKey := storetypes.NewTransientStoreKey("test")
//
//	//accountStoreKey := storetypes.NewKVStoreKey(accounttypes.StoreKey) // sdk.NewKVStoreKey(accounttypes.StoreKey)
//	//accountSubspace := paramtypes.NewSubspace(cdc, amino, accountStoreKey, testTransientStoreKey, accounttypes.ModuleName)
//	stubAccKeeper := accountkeeper.AccountKeeper{}
//	//accountKeeper := accountkeeper.NewAccountKeeper(cdc, runtime.NewKVStoreService(accountStoreKey), accounttypes.ProtoBaseAccount, accountSubspace, accounttypes.ProtoBaseAccount, "link", nil)
//
//	bankSubspace := paramtypes.NewSubspace(cdc, amino, storeKey, testTransientStoreKey, banktypes.StoreKey)
//
//	//cdc codec.Codec, storeService store.KVStoreService, ak types.AccountKeeper, paramSpace paramtypes.Subspace,
//	//	blockedAddr map[string]bool, deactMultiSend bool, authority string, logger log.Logger,
//	return NewBaseKeeper(cdc, storeKey, stubAccKeeper, bankSubspace, nil, false)
//}
//
//func setupContext(t *testing.T, storeKey *storetypes.KVStoreKey) sdk.Context {
//	tkey := storetypes.NewTransientStoreKey("transient_test")
//	testCtx := testutil.DefaultContextWithDB(t, storeKey, tkey)
//
//	return testCtx.Ctx
//	//db := dbm.NewMemDB()
//	//stateStore := store.NewCommitMultiStore(db)
//	//stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
//	//require.NoError(t, stateStore.LoadLatestVersion())
//	//
//	//return sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
//}
//
//func TestInactiveAddr(t *testing.T) {
//	storeKey := storetypes.NewKVStoreKey(banktypes.StoreKey)
//	bankKeeper := setupKeeper(storeKey)
//	ctx := setupContext(t, storeKey)
//
//	addr := genAddress()
//
//	require.Equal(t, 0, len(bankKeeper.inactiveAddrs))
//
//	bankKeeper.addToInactiveAddr(ctx, addr)
//	require.True(t, bankKeeper.isStoredInactiveAddr(ctx, addr))
//
//	// duplicate addition, no error expects.
//	bankKeeper.addToInactiveAddr(ctx, addr)
//	require.True(t, bankKeeper.isStoredInactiveAddr(ctx, addr))
//
//	bankKeeper.deleteFromInactiveAddr(ctx, addr)
//	require.False(t, bankKeeper.isStoredInactiveAddr(ctx, addr))
//
//	addr2 := genAddress()
//	require.False(t, bankKeeper.isStoredInactiveAddr(ctx, addr2))
//
//	// expect no error
//	bankKeeper.deleteFromInactiveAddr(ctx, addr2)
//
//	// test loadAllInactiveAddrs
//	bankKeeper.addToInactiveAddr(ctx, addr)
//	bankKeeper.addToInactiveAddr(ctx, addr2)
//	require.Equal(t, 0, len(bankKeeper.inactiveAddrs))
//	bankKeeper.loadAllInactiveAddrs(ctx)
//	require.Equal(t, 2, len(bankKeeper.inactiveAddrs))
//}

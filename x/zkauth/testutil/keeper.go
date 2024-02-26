package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Finschia/finschia-sdk/x/zkauth/keeper"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"

	"github.com/Finschia/ostracon/libs/log"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/Finschia/finschia-sdk/codec"
	codectypes "github.com/Finschia/finschia-sdk/codec/types"
	"github.com/Finschia/finschia-sdk/store"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
)

func ZkAuthKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	verificationKey, err := os.ReadFile("../testutil/testdata/verification_key.json")
	require.NoError(t, err)

	zwksMap := types.NewJWKs()

	homeDir := filepath.Join(t.TempDir(), "x_zkauth_keeper_test")
	t.Log("home dir: ", homeDir)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		zwksMap,
		types.NewZKAuthVerifier(verificationKey),
		homeDir,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return k, ctx
}

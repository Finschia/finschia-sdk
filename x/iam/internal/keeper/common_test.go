package keeper

import (
	"github.com/line/link/x/iam/internal/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type testInput struct {
	cdc    *codec.Codec
	ctx    sdk.Context
	keeper Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()
	types.RegisterCodec(cdc)
	cdc.RegisterConcrete(&permission{}, "link/Permission", nil)
	codec.RegisterCrypto(cdc)
	return cdc
}

func setupTestInput(t *testing.T) testInput {

	keyIam := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyIam, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.NoError(t, err)

	cdc := newTestCodec()

	keeper := NewKeeper(cdc, keyIam)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	return testInput{cdc: cdc, ctx: ctx, keeper: keeper}
}

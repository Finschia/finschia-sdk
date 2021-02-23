package keeper

import (
	"testing"

	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/store"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/contract/internal/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type testInput struct {
	cdc    *codec.Codec
	ctx    sdk.Context
	keeper ContractKeeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func setupTestInput(t *testing.T) testInput {
	keyContract := sdk.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyContract, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.NoError(t, err)

	cdc := newTestCodec()

	keeper := NewContractKeeper(cdc, keyContract)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	return testInput{cdc: cdc, ctx: ctx, keeper: keeper}
}

func TestKeeper(t *testing.T) {
	testInput := setupTestInput(t)
	_, ctx, keeper := testInput.cdc, testInput.ctx, testInput.keeper
	for i := 0; i < 10000; i++ {
		contractID := keeper.NewContractID(ctx)
		require.True(t, VerifyContractID(contractID))
	}
}

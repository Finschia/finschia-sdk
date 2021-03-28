package account

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/line/lbm-sdk/v2/x/account/internal/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
)

type TestInput struct {
	Cdc *codec.Codec
	Ctx sdk.Context
	Ak  auth.AccountKeeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func setupTestInput(t *testing.T) TestInput {
	keyAuth := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAuth, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.NoError(t, err)

	cdc := newTestCodec()

	// init params keeper and subspaces
	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams)
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	// add keepers
	accountKeeper := auth.NewAccountKeeper(cdc, keyAuth, authSubspace, auth.ProtoBaseAccount)
	accountKeeper.NewAccountWithAddress(ctx, addr1)

	return TestInput{Cdc: cdc, Ctx: ctx, Ak: accountKeeper}
}

func TestHandlerCreateAccount(t *testing.T) {
	input := setupTestInput(t)
	ctx, keeper := input.Ctx, input.Ak

	h := NewHandler(keeper)

	// creating the account addr2 succeeds at first
	{
		msg := types.NewMsgCreateAccount(addr1, addr2)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}

	// creating the account addr2 twice fails
	{
		msg := types.NewMsgCreateAccount(addr1, addr2)
		_, err := h(ctx, msg)
		require.Error(t, err)
	}
}

func TestHandlerEmpty(t *testing.T) {
	input := setupTestInput(t)
	ctx, keeper := input.Ctx, input.Ak

	h := NewHandler(keeper)

	// message test
	{
		msg := types.MsgEmpty{From: addr1}
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}

	// invalid message
	{
		msg := types.MsgEmpty{From: nil}
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing signer address").Error())
	}
}

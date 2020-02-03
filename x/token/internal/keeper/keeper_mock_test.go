package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/line/link/x/token/internal/types"
	mock_types "github.com/line/link/x/token/mock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
)

func TestKeeper_ModifyTokenURI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Log("modify succeed")
	{
		addr1, cd, mockIamKeeper, keeper, mockMultiStore, cliCtx, mockKvStore, persistedToken, persistedTokenSymbolKey,
			persistedTokenJSON, tokenURIModifiedToken := prepare(ctrl)

		mockMultiStore.EXPECT().GetKVStore(gomock.Any()).Return(mockKvStore).Times(2)
		mockKvStore.EXPECT().Get(persistedTokenSymbolKey).Return(persistedTokenJSON).Times(1)
		mockKvStore.EXPECT().Has(persistedTokenSymbolKey).Return(true).Times(1)
		mockIamKeeper.EXPECT().HasPermission(cliCtx, addr1, types.NewModifyTokenURIPermission(persistedToken.GetDenom())).
			Return(true).Times(1)
		mockKvStore.EXPECT().Set(persistedTokenSymbolKey, cd.MustMarshalBinaryBare(tokenURIModifiedToken)).Times(1)

		err := keeper.ModifyTokenURI(cliCtx, addr1, persistedToken.GetSymbol(), persistedToken.GetTokenID(), tokenURIModifiedToken.GetTokenURI())
		require.NoError(t, err)
	}
	t.Log("could not found persisted token")
	{
		addr1, _, mockIamKeeper, keeper, mockMultiStore, cliCtx, mockKvStore, persistedToken, persistedTokenSymbolKey, _, _ := prepare(ctrl)

		mockMultiStore.EXPECT().GetKVStore(gomock.Any()).Return(mockKvStore).Times(1)
		mockKvStore.EXPECT().Get(persistedTokenSymbolKey).Return(nil).Times(1)
		mockIamKeeper.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(0)
		mockKvStore.EXPECT().Has(gomock.Any()).Times(0)
		mockKvStore.EXPECT().Set(gomock.Any(), gomock.Any()).Times(0)
		err := keeper.ModifyTokenURI(cliCtx, addr1, persistedToken.GetSymbol(), persistedToken.GetTokenID(), persistedToken.GetTokenURI())
		require.Error(t, err)
	}
	t.Log("missing tokenSymbolKey in modifying step")
	{
		addr1, _, mockIamKeeper, keeper, mockMultiStore, cliCtx, mockKvStore, tt, tsk, ttJSON, _ := prepare(ctrl)

		mockMultiStore.EXPECT().GetKVStore(gomock.Any()).Return(mockKvStore).Times(2)
		mockKvStore.EXPECT().Get(tsk).Return(ttJSON).Times(1)
		mockIamKeeper.EXPECT().HasPermission(cliCtx, addr1, types.NewModifyTokenURIPermission(tt.GetDenom())).Return(true).Times(1)
		mockKvStore.EXPECT().Has(tsk).Return(false).Times(1)
		mockKvStore.EXPECT().Set(gomock.Any(), gomock.Any()).Times(0)

		err := keeper.ModifyTokenURI(cliCtx, addr1, tt.GetSymbol(), tt.GetTokenID(), tt.GetTokenURI())
		require.Error(t, err)
	}
	t.Log("failed - no permission")
	{
		addr1, _, mockIamKeeper, keeper, mockMultiStore, cliCtx, mockKvStore, tt, tsk, ttJSON, _ := prepare(ctrl)

		mockMultiStore.EXPECT().GetKVStore(gomock.Any()).Return(mockKvStore).Times(1)
		mockKvStore.EXPECT().Get(tsk).Return(ttJSON).Times(1)
		mockIamKeeper.EXPECT().HasPermission(cliCtx, addr1, types.NewModifyTokenURIPermission(tt.GetDenom())).Return(false).Times(1)
		mockKvStore.EXPECT().Set(gomock.Any(), gomock.Any()).Times(0)

		err := keeper.ModifyTokenURI(cliCtx, addr1, tt.GetSymbol(), tt.GetTokenID(), tt.GetTokenURI())
		require.Error(t, err)
	}
}

func prepare(ctrl *gomock.Controller) (sdk.AccAddress, *codec.Codec, *mock_types.MockIamKeeper, Keeper,
	*mock_types.MockMultiStore, sdk.Context, *mock_types.MockKVStore, types.NFT, []byte, []byte, types.NFT) {
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	cd := types.ModuleCdc
	mockIamKeeper := mock_types.NewMockIamKeeper(ctrl)
	keeper := Keeper{
		supplyKeeper:  mock_types.NewMockSupplyKeeper(ctrl),
		iamKeeper:     mockIamKeeper,
		accountKeeper: mock_types.NewMockAccountKeeper(ctrl),
		storeKey:      mock_types.NewMockStoreKey(ctrl),
		cdc:           cd,
	}
	mockMultiStore := mock_types.NewMockMultiStore(ctrl)

	cliCtx := sdk.NewContext(mockMultiStore, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	mockKvStore := mock_types.NewMockKVStore(ctrl)
	tokenSymbol := "testSymbol"
	tt := types.NewNFT(
		"testToken",
		tokenSymbol,
		"",
		addr1,
	)
	tokenURIModifiedToken := types.NewNFT(
		"testToken",
		tokenSymbol,
		"modified",
		addr1,
	)

	tsk := types.TokenDenomKey(tt.GetSymbol())
	ttJSON := cd.MustMarshalBinaryBare(tt)
	return addr1, cd, mockIamKeeper, keeper, mockMultiStore, cliCtx, mockKvStore, tt, tsk, ttJSON, tokenURIModifiedToken
}

package scenario

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/types"
	linktypes "github.com/line/link/types"
	"github.com/line/link/x/collection"
	"github.com/stretchr/testify/require"
)

func TestNewMsgBuilder(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet),
		linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given info
	testAddress, err := sdk.AccAddressFromBech32("link1jrhv7xwgt50hn4nyf5m2clnyq0n4d4akam56rf")
	require.NoError(t, err)
	numNFTPerUser := 10
	info := Info{
		config: types.Config{CoinName: tests.TestCoinName},
		stateParams: map[string]string{"token_contract_id": tests.TokenContractID,
			"collection_contract_id": tests.CollectionContractID, "ft_token_id": tests.FTTokenID,
			"nft_token_type": tests.NFTTokenType, "tx_hash": tests.TxHash,
			"num_nft_per_user": fmt.Sprintf("%d", numNFTPerUser)},
	}

	t.Log("Happy path")
	{
		// When
		msgBuilder, err := NewMsgBuilder(testAddress, info, 0, nil)
		require.NoError(t, err)

		// Then
		require.Equal(t, numNFTPerUser, msgBuilder.numNFTPerUser)
		require.Len(t, msgBuilder.nftTokenIDs, numNFTPerUser)
		require.Len(t, msgBuilder.handlers, 26)
	}
	t.Log("No num_nft_per_user info")
	{
		// Given no num_nft_per_user info
		delete(info.stateParams, "num_nft_per_user")

		// When
		msgBuilder, err := NewMsgBuilder(testAddress, info, 0, nil)
		require.NoError(t, err)

		// Then
		require.Equal(t, 0, msgBuilder.numNFTPerUser)
		require.Len(t, msgBuilder.nftTokenIDs, 0)
		require.Len(t, msgBuilder.handlers, 26)
	}
	t.Log("Invalid num_nft_per_user info")
	{
		// Given no num_nft_per_user info
		info.stateParams["num_nft_per_user"] = "invalid"

		// When
		_, err := NewMsgBuilder(testAddress, info, 0, nil)

		// Then
		require.EqualError(t, err, fmt.Sprintf("strconv.Atoi: parsing \"%s\": invalid syntax", info.stateParams["num_nft_per_user"]))
	}
	t.Log("num_nft_per_user overflowed")
	{
		// Given no num_nft_per_user info
		info.stateParams["num_nft_per_user"] = strconv.Itoa(math.MaxInt32)

		// When
		_, err := NewMsgBuilder(testAddress, info, 1, nil)

		// Then
		require.EqualError(t, err, types.NFTTokenIDOverFlowError{}.Error())
	}
}

func TestMsgBuilder_GetHandler(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet),
		linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given MsgBuilder
	testAddress, err := sdk.AccAddressFromBech32("link1jrhv7xwgt50hn4nyf5m2clnyq0n4d4akam56rf")
	require.NoError(t, err)
	numNFTPerUser := 10
	info := Info{
		config: types.Config{CoinName: tests.TestCoinName},
		stateParams: map[string]string{"token_contract_id": tests.TokenContractID,
			"collection_contract_id": tests.CollectionContractID, "ft_token_id": tests.FTTokenID,
			"nft_token_type": tests.NFTTokenType, "tx_hash": tests.TxHash,
			"num_nft_per_user": fmt.Sprintf("%d", numNFTPerUser)},
	}
	m, err := NewMsgBuilder(testAddress, info, 0, nil)
	require.NoError(t, err)

	var cases = []struct {
		handlerType string
		handler     func() (sdk.Msg, error)
	}{
		{"MsgEmpty", m.newMsgEmpty},
		{"MsgSend", m.newMsgSend},
		{"MsgIssue", m.newMsgIssue},
		{"MsgMint", m.newMsgMint},
		{"MsgTransfer", m.newMsgTransfer},
		{"MsgModifyToken", m.newMsgModifyToken},
		{"MsgModifyTokenName", m.newMsgModifyTokenName},
		{"MsgModifyTokenURI", m.newMsgModifyTokenURI},
		{"MsgBurn", m.newMsgBurn},
		{"MsgCreateCollection", m.newMsgCreateCollection},
		{"MsgApprove", m.newMsgApprove},
		{"MsgIssueFT", m.newMsgIssueFT},
		{"MsgMintFT", m.newMsgMintFT},
		{"MsgTransferFT", m.newMsgTransferFT},
		{"MsgBurnFT", m.newMsgBurnFT},
		{"MsgModifyCollection", m.newMsgModifyToken},
		{"MsgIssueNFT", m.newMsgIssueNFT},
		{"MsgMintNFT", m.newMsgMintNFT},
		{"MsgMintOneNFT", m.newMsgMintOneNFT},
		{"MsgMintFiveNFT", m.newMsgMintFiveNFTs},
		{"MsgAttach", m.newMsgAttach},
		{"MsgDetach", m.newMsgDetach},
		{"MsgTransferNFT", m.newMsgTransferFT},
		{"MsgMultiTransferNFT", m.newMsgTransferNFT},
		{"MsgBurnNFT", m.newMsgBurnNFT},
		{"MsgGrantPermission", m.newMsgGrantPermission},
	}

	for _, tt := range cases {
		// When
		handler, err := m.GetHandler(tt.handlerType)
		require.NoError(t, err)

		// Then
		require.IsType(t, tt.handler, handler)
	}
}

func TestMsgBuilder(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet),
		linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given MsgBuilder
	testAddress, err := sdk.AccAddressFromBech32("link1jrhv7xwgt50hn4nyf5m2clnyq0n4d4akam56rf")
	require.NoError(t, err)
	numNFTPerUser := 10
	info := Info{
		config: types.Config{CoinName: tests.TestCoinName},
		stateParams: map[string]string{"token_contract_id": tests.TokenContractID,
			"collection_contract_id": tests.CollectionContractID, "ft_token_id": tests.FTTokenID,
			"nft_token_type": tests.NFTTokenType, "tx_hash": tests.TxHash,
			"num_nft_per_user": fmt.Sprintf("%d", numNFTPerUser)},
	}
	m, err := NewMsgBuilder(testAddress, info, 0, nil)
	require.NoError(t, err)

	var cases = []struct {
		handlerType string
		msgType     string
	}{
		{"MsgEmpty", "empty"},
		{"MsgSend", "send"},
		{"MsgIssue", "issue_token"},
		{"MsgMint", "mint"},
		{"MsgTransfer", "transfer_ft"},
		{"MsgModifyToken", "modify_token"},
		{"MsgModifyTokenName", "modify_token"},
		{"MsgModifyTokenURI", "modify_token"},
		{"MsgBurn", "burn"},
		{"MsgCreateCollection", "create_collection"},
		{"MsgApprove", "approve_collection"},
		{"MsgIssueFT", "issue_ft"},
		{"MsgMintFT", "mint_ft"},
		{"MsgTransferFT", "transfer_ft"},
		{"MsgBurnFT", "burn_ft"},
		{"MsgModifyCollection", "modify_token"},
		{"MsgIssueNFT", "issue_nft"},
		{"MsgMintOneNFT", "mint_nft"},
		{"MsgMintFiveNFT", "mint_nft"},
		{"MsgAttach", "attach"},
		{"MsgDetach", "detach"},
		{"MsgTransferNFT", "transfer_nft"},
		{"MsgMultiTransferNFT", "transfer_nft"},
		{"MsgBurnNFT", "burn_nft"},
		{"MsgGrantPermission", "grant_perm"},
	}

	for _, tt := range cases {
		// When
		handler, err := m.GetHandler(tt.handlerType)
		require.NoError(t, err)
		msg, err := handler()
		require.NoError(t, err)

		// Then
		require.Equal(t, tt.msgType, msg.Type())
	}
}

func TestMsgBuilder_NewMsgMintNFT(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet),
		linktypes.Bech32PrefixAccPub(tests.TestNet))
	testAddress, err := sdk.AccAddressFromBech32("link1jrhv7xwgt50hn4nyf5m2clnyq0n4d4akam56rf")
	require.NoError(t, err)

	t.Log("Happy Path")
	{
		var cases = []struct {
			numNFT            []string
			expectedLenParams int
		}{
			{[]string{"1"}, 1},
			{[]string{"5"}, 5},
			{[]string{"10"}, 10},
		}

		for _, tt := range cases {
			// Given MsgBuilder
			info := Info{
				config: types.Config{CoinName: tests.TestCoinName},
				stateParams: map[string]string{"collection_contract_id": tests.CollectionContractID,
					"nft_token_type": tests.NFTTokenType},
			}
			m, err := NewMsgBuilder(testAddress, info, 0, tt.numNFT)
			require.NoError(t, err)

			// When
			msg, err := m.newMsgMintNFT()
			require.NoError(t, err)

			// Then
			msgMintNFT, ok := msg.(collection.MsgMintNFT)
			require.True(t, ok)
			require.Len(t, msgMintNFT.MintNFTParams, tt.expectedLenParams)
			require.NoError(t, msgMintNFT.ValidateBasic())
		}
	}
	t.Log("Invalid scenario params")
	{
		// Given MsgBuilder
		invalidParam := "invalid"
		info := Info{
			config: types.Config{CoinName: tests.TestCoinName},
			stateParams: map[string]string{"collection_contract_id": tests.CollectionContractID,
				"nft_token_type": tests.NFTTokenType},
		}
		m, err := NewMsgBuilder(testAddress, info, 0, []string{invalidParam})
		require.NoError(t, err)

		// When
		_, err = m.newMsgMintNFT()

		// Then
		require.EqualError(t, err, fmt.Sprintf("strconv.Atoi: parsing \"%s\": invalid syntax", invalidParam))
	}
	t.Log("negative num NFT")
	{
		// Given MsgBuilder
		info := Info{
			config: types.Config{CoinName: tests.TestCoinName},
			stateParams: map[string]string{"collection_contract_id": tests.CollectionContractID,
				"nft_token_type": tests.NFTTokenType},
		}
		m, err := NewMsgBuilder(testAddress, info, 0, []string{"-1"})
		require.NoError(t, err)

		// When
		msg, err := m.newMsgMintNFT()
		require.NoError(t, err)

		// Then
		msgMintNFT, ok := msg.(collection.MsgMintNFT)
		require.True(t, ok)
		require.EqualError(t, msgMintNFT.ValidateBasic(), "required field cannot be empty: params cannot be empty")
	}
}

package scenario

import (
	"fmt"
	"net/http"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	linktypes "github.com/line/link/types"
	"github.com/stretchr/testify/require"
)

func TestGenerateRegisterAccountMsgs(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Master Key Wallet
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	require.NoError(t, err)
	masterWallet, err := hdWallet.GetKeyWallet(1, 0)
	require.NoError(t, err)
	// Given Config
	config := types.Config{
		MsgsPerTxPrepare:  tests.TestMsgsPerTxPrepare,
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		TargetURL:         tests.TestTargetURL,
		Duration:          tests.TestDuration,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
	}

	msgs, err := GenerateRegisterAccountMsgs(masterWallet.Address(), hdWallet, config)
	require.NoError(t, err)

	require.Len(t, msgs, tests.TestTPS*tests.TestDuration)
	for _, msg := range msgs {
		require.Equal(t, "send", msg.Type())
		require.Equal(t, masterWallet.Address(), msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())
	}
}

func TestGenerateGrantPermissionMsgs(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Master Key Wallet
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	require.NoError(t, err)
	masterWallet, err := hdWallet.GetKeyWallet(1, 0)
	require.NoError(t, err)
	// Given Config
	config := types.Config{
		MsgsPerTxPrepare: tests.TestMsgsPerTxPrepare,
		TPS:              tests.TestTPS,
		TargetURL:        tests.TestTargetURL,
		Duration:         tests.TestDuration,
		ChainID:          tests.TestChainID,
		CoinName:         tests.TestCoinName,
	}

	testCases := []struct {
		contractID string
		module     string
	}{
		{"9be17165", "token"},
		{"678c146a", "collection"},
	}
	for _, tt := range testCases {
		msgs, err := GenerateGrantPermissionMsgs(masterWallet.Address(), hdWallet, config, tt.contractID, tt.module,
			[]string{"mint", "burn", "modify"})
		require.NoError(t, err)

		require.Len(t, msgs, tests.TestTPS*tests.TestDuration*3)
		for _, msg := range msgs {
			require.Equal(t, "grant_perm", msg.Type())
			require.Equal(t, masterWallet.Address(), msg.GetSigners()[0])
			require.NoError(t, msg.ValidateBasic())
		}
	}
}

func TestGenerateMintFTMsgs(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	TestContractID := "678c146a"
	// Given Master Key Wallet
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	require.NoError(t, err)
	masterWallet, err := hdWallet.GetKeyWallet(1, 0)
	require.NoError(t, err)
	// Given Config
	config := types.Config{
		MsgsPerTxPrepare: tests.TestMsgsPerTxPrepare,
		TPS:              tests.TestTPS,
		TargetURL:        tests.TestTargetURL,
		Duration:         tests.TestDuration,
		ChainID:          tests.TestChainID,
		CoinName:         tests.TestCoinName,
	}

	msgs, err := GenerateMintFTMsgs(masterWallet.Address(), hdWallet, config, TestContractID, "0000000100000000")
	require.NoError(t, err)

	require.Len(t, msgs, tests.TestTPS*tests.TestDuration)
	for _, msg := range msgs {
		require.Equal(t, "mint_ft", msg.Type())
		require.Equal(t, masterWallet.Address(), msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())
	}
}

func TestGenerateMintNFTMsgs(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	TestContractID := "678c146a"
	TestTokenType := "10000001"
	// Given Master Key Wallet
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	require.NoError(t, err)
	masterWallet, err := hdWallet.GetKeyWallet(1, 0)
	require.NoError(t, err)
	// Given Config
	config := types.Config{
		MsgsPerTxPrepare: tests.TestMsgsPerTxPrepare,
		TPS:              tests.TestTPS,
		TargetURL:        tests.TestTargetURL,
		Duration:         tests.TestDuration,
		ChainID:          tests.TestChainID,
		CoinName:         tests.TestCoinName,
	}

	msgs, err := GenerateMintNFTMsgs(masterWallet.Address(), hdWallet, config, TestContractID, TestTokenType, 2)
	require.NoError(t, err)

	require.Len(t, msgs, tests.TestTPS*tests.TestDuration*2)
	for _, msg := range msgs {
		require.Equal(t, "mint_nft", msg.Type())
		require.Equal(t, masterWallet.Address(), msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())
	}
}

func TestIssueToken(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given LinkService
	linkService := service.NewLinkService(&http.Client{}, app.MakeCodec(), server.URL)
	txBuilder := transaction.NewTxBuilder(tests.TestMaxGasPrepare).WithChainID(tests.TestChainID)
	// Given Master Key Wallet
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	require.NoError(t, err)
	masterWallet, err := hdWallet.GetKeyWallet(0, 0)
	require.NoError(t, err)

	contractID, _, err := IssueToken(linkService, txBuilder, masterWallet)
	require.NoError(t, err)

	require.Equal(t, "9be17165", contractID)
}

func TestCreateCollection(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given LinkService
	linkService := service.NewLinkService(&http.Client{}, app.MakeCodec(), server.URL)
	txBuilder := transaction.NewTxBuilder(tests.TestMaxGasPrepare).WithChainID(tests.TestChainID)
	// Given Master Key Wallet
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	require.NoError(t, err)
	masterWallet, err := hdWallet.GetKeyWallet(0, 0)
	require.NoError(t, err)

	contractID, err := CreateCollection(linkService, txBuilder, masterWallet)
	require.NoError(t, err)

	require.Equal(t, "678c146a", contractID)
}

func TestIssueFT(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given LinkService
	linkService := service.NewLinkService(&http.Client{}, app.MakeCodec(), server.URL)
	txBuilder := transaction.NewTxBuilder(tests.TestMaxGasPrepare).WithChainID(tests.TestChainID)
	// Given Master Key Wallet
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	require.NoError(t, err)
	masterWallet, err := hdWallet.GetKeyWallet(0, 0)
	require.NoError(t, err)

	ftTokenID, err := IssueFT(linkService, txBuilder, masterWallet, "678c146a")
	require.NoError(t, err)

	require.Equal(t, "0000000100000000", ftTokenID)
}

func TestIssueNFT(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given LinkService
	linkService := service.NewLinkService(&http.Client{}, app.MakeCodec(), server.URL)
	txBuilder := transaction.NewTxBuilder(tests.TestMaxGasPrepare).WithChainID(tests.TestChainID)
	// Given Master Key Wallet
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	require.NoError(t, err)
	masterWallet, err := hdWallet.GetKeyWallet(0, 0)
	require.NoError(t, err)

	nftTokenType, err := IssueNFT(linkService, txBuilder, masterWallet, "678c146a")
	require.NoError(t, err)

	require.Equal(t, "10000001", nftTokenType)
}

func TestBuildTxTarget(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet),
		linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Master Key Wallet
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	require.NoError(t, err)
	keyWallet, err := hdWallet.GetKeyWallet(1, 0)
	require.NoError(t, err)
	// Given Config
	config := types.Config{
		MsgsPerTxPrepare:  tests.TestMsgsPerTxPrepare,
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		TargetURL:         server.URL,
		Duration:          tests.TestDuration,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
	}
	// Given info
	numNFTPerUser := 30
	info := Info{
		linkService:   service.NewLinkService(&http.Client{}, app.MakeCodec(), config.TargetURL),
		targetBuilder: NewTargetBuilder(app.MakeCodec(), config.TargetURL),
		txBuilder:     transaction.NewTxBuilder(uint64(config.MaxGasLoadTest)).WithChainID(config.ChainID),
		config:        config,
		stateParams: map[string]string{"token_contract_id": tests.TokenContractID,
			"collection_contract_id": tests.CollectionContractID, "ft_token_id": tests.FTTokenID,
			"nft_token_type": tests.NFTTokenType, "tx_hash": tests.TxHash,
			"num_nft_per_user": fmt.Sprintf("%d", numNFTPerUser)},
	}
	msgList := []string{"MsgEmpty", "MsgSend", "MsgIssue", "MsgMint", "MsgTransfer", "MsgModifyToken", "MsgModifyTokenName",
		"MsgModifyTokenURI", "MsgBurn", "MsgCreateCollection", "MsgApprove", "MsgIssueFT", "MsgMintFT", "MsgTransferFT",
		"MsgBurnFT", "MsgModifyCollection", "MsgIssueNFT", "MsgMintOneNFT", "MsgAttach", "MsgDetach",
		"MsgTransferNFT", "MsgMultiTransferNFT", "MsgBurnNFT", "MsgGrantPermission"}

	// When
	target, err := BuildTxTarget(info, keyWallet, 0, msgList, nil)
	require.NoError(t, err)

	// Then
	require.Equal(t, "POST", target.Method)
	require.Equal(t, fmt.Sprintf("%s%s", server.URL, TxURL), target.URL)
	// Then tx target is valid
	var req rest.BroadcastReq
	require.NoError(t, app.MakeCodec().UnmarshalJSON(target.Body, &req))
	require.Equal(t, service.BroadcastSync, req.Mode)

	require.Len(t, req.Tx.Msgs, len(msgList)*tests.TestMsgsPerTxLoadTest)
	for _, msg := range req.Tx.Msgs {
		require.NoError(t, msg.ValidateBasic())
	}
}

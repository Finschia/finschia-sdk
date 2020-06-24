package scenario

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	linktypes "github.com/line/link/types"
	"github.com/line/link/x/coin"
	"github.com/line/link/x/collection"
	"github.com/line/link/x/token"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxAndQueryAllScenario struct {
	linkService   *service.LinkService
	targetBuilder *TargetBuilder
	txBuilder     transaction.TxBuilderWithoutKeybase
	config        types.Config
	params        map[string]string
}

func (s *TxAndQueryAllScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet) ([]sdk.Msg, map[string]string, error) {
	masterAddress := masterKeyWallet.Address()
	tokenContractID, txHash, err := IssueToken(s.linkService, s.txBuilder, masterKeyWallet)
	if err != nil {
		return nil, nil, err
	}
	collectionContractID, err := CreateCollection(s.linkService, s.txBuilder, masterKeyWallet)
	if err != nil {
		return nil, nil, err
	}
	ftTokenID, err := IssueFT(s.linkService, s.txBuilder, masterKeyWallet, collectionContractID)
	if err != nil {
		return nil, nil, err
	}
	nftTokenType, err := IssueNFT(s.linkService, s.txBuilder, masterKeyWallet, collectionContractID)
	if err != nil {
		return nil, nil, err
	}

	msgs, err := GenerateRegisterAccountMsgs(masterAddress, hdWallet, s.config)
	if err != nil {
		return nil, nil, err
	}
	grantTokenPermMsgs, err := GenerateGrantPermissionMsgs(masterAddress, hdWallet, s.config, tokenContractID,
		"token", []string{"mint", "burn", "modify"})
	if err != nil {
		return nil, nil, err
	}
	grantCollectionPermMsgs, err := GenerateGrantPermissionMsgs(masterAddress, hdWallet, s.config,
		collectionContractID, "collection", []string{"issue", "mint", "burn", "modify"})
	if err != nil {
		return nil, nil, err
	}
	mintNFTMsgs, err := GenerateMintNFTMsgs(masterAddress, hdWallet, s.config, collectionContractID, nftTokenType, 13)
	if err != nil {
		return nil, nil, err
	}

	msgs = append(msgs, grantTokenPermMsgs...)
	msgs = append(msgs, grantCollectionPermMsgs...)
	msgs = append(msgs, mintNFTMsgs...)

	return msgs, map[string]string{"token_contract_id": tokenContractID, "collection_contract_id": collectionContractID,
		"ft_token_id": ftTokenID, "nft_token_type": nftTokenType, "tx_hash": txHash}, nil
}

func (s *TxAndQueryAllScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	account, err := s.linkService.GetAccount(keyWallet.Address().String())
	if err != nil {
		return nil, 0, err
	}

	numMsgIncrement := 29
	repeatCount := (s.config.MsgsPerTxLoadTest + numMsgIncrement - 1) / numMsgIncrement // round up

	msgs := make([]sdk.Msg, numMsgIncrement*repeatCount)
	tokenContractID := s.params["token_contract_id"]
	collectionContractID := s.params["collection_contract_id"]
	ftTokenID := s.params["ft_token_id"]
	nftTokenType := s.params["nft_token_type"]
	nftTokenIDs := make([]string, 13)
	for i := 0; i < 13; i++ {
		nftTokenIDs[i] = fmt.Sprintf("%s%08x", nftTokenType, 13*walletIndex+i+1)
	}
	from := account.Address
	coins := sdk.NewCoins(sdk.NewCoin(s.config.CoinName, sdk.NewInt(1)))
	for i := 0; i < repeatCount; i++ {
		to := secp256k1.GenPrivKey().PubKey().Address().Bytes()
		base := numMsgIncrement * i

		msgs[base] = coin.NewMsgSend(from, to, coins)
		msgs[base+1] = coin.NewMsgSend(from, to, coins)

		msgs[base+2] = token.NewMsgIssue(from, to, "token", "TOK", "{}", "uri", sdk.NewInt(1),
			sdk.NewInt(8), true)
		msgs[base+3] = token.NewMsgMint(from, tokenContractID, from, sdk.NewInt(2))
		msgs[base+4] = token.NewMsgMint(from, tokenContractID, from, sdk.NewInt(2))
		msgs[base+5] = token.NewMsgTransfer(from, to, tokenContractID, sdk.NewInt(1))
		msgs[base+6] = token.NewMsgTransfer(from, to, tokenContractID, sdk.NewInt(1))
		msgs[base+7] = token.NewMsgModify(from, tokenContractID,
			linktypes.NewChangesWithMap(map[string]string{"name": fmt.Sprintf("token%d-%d", walletIndex, i)}))
		msgs[base+8] = token.NewMsgModify(from, tokenContractID,
			linktypes.NewChangesWithMap(map[string]string{"img_uri": fmt.Sprintf("uri%d-%d", walletIndex, i)}))
		msgs[base+9] = token.NewMsgBurn(from, tokenContractID, sdk.NewInt(1))
		msgs[base+10] = token.NewMsgBurn(from, tokenContractID, sdk.NewInt(1))

		msgs[base+11] = collection.NewMsgCreateCollection(from, "name", "{}", "uri")
		msgs[base+12] = collection.NewMsgApprove(from, collectionContractID, to)
		msgs[base+13] = collection.NewMsgIssueFT(from, to, collectionContractID, "collection", "{}", sdk.NewInt(1),
			sdk.NewInt(8), true)
		msgs[base+14] = collection.NewMsgMintFT(from, collectionContractID, from, collection.NewCoin(ftTokenID,
			sdk.NewInt(4)))
		msgs[base+15] = collection.NewMsgTransferFT(from, collectionContractID, to, collection.NewCoin(ftTokenID,
			sdk.NewInt(1)))
		msgs[base+16] = collection.NewMsgTransferFT(from, collectionContractID, to, collection.NewCoin(ftTokenID,
			sdk.NewInt(1)))
		msgs[base+17] = collection.NewMsgBurnFT(from, collectionContractID, collection.NewCoin(ftTokenID,
			sdk.NewInt(1)))

		msgs[base+18] = collection.NewMsgModify(from, collectionContractID, ftTokenID[:8], ftTokenID[8:],
			linktypes.NewChangesWithMap(map[string]string{"name": fmt.Sprintf("name%d-%d", walletIndex, i)}))

		msgs[base+19] = collection.NewMsgIssueNFT(from, collectionContractID, "ft", "{}")
		msgs[base+20] = collection.NewMsgMintNFT(from, collectionContractID, to, collection.NewMintNFTParam("name",
			"{}", nftTokenType))
		var params []collection.MintNFTParam
		for j := 0; j < 5; j++ {
			params = append(params, collection.NewMintNFTParam("name", "{}", nftTokenType))
		}
		msgs[base+21] = collection.NewMsgMintNFT(from, collectionContractID, to, params...)
		msgs[base+22] = collection.NewMsgAttach(from, collectionContractID, nftTokenIDs[0], nftTokenIDs[1])
		msgs[base+23] = collection.NewMsgDetach(from, collectionContractID, nftTokenIDs[1])
		msgs[base+24] = collection.NewMsgTransferNFT(from, collectionContractID, to, nftTokenIDs[1])
		msgs[base+25] = collection.NewMsgTransferNFT(from, collectionContractID, to, nftTokenIDs[2])
		msgs[base+26] = collection.NewMsgTransferNFT(from, collectionContractID, to, nftTokenIDs[3:8]...) // array
		msgs[base+27] = collection.NewMsgTransferNFT(from, collectionContractID, to, nftTokenIDs[8:]...)  // array
		msgs[base+28] = collection.NewMsgBurnNFT(from, collectionContractID, nftTokenIDs[0])
	}

	stdTx, err := s.txBuilder.WithAccountNumber(account.AccountNumber).WithSequence(account.Sequence).
		BuildAndSign(keyWallet.PrivateKey(), msgs)
	if err != nil {
		return nil, 0, err
	}

	target, err := s.targetBuilder.MakeTxTarget(stdTx, service.BroadcastSync)
	targets := []*vegeta.Target{target}
	targets = s.addQueryTargets(targets, from, nftTokenIDs[1])

	return &targets, len(targets), err
}

func (s *TxAndQueryAllScenario) addQueryTargets(targets []*vegeta.Target, from sdk.AccAddress, nftTokenID string) []*vegeta.Target {
	tokenContractID := s.params["token_contract_id"]
	collectionContractID := s.params["collection_contract_id"]
	ftTokenID := s.params["ft_token_id"]
	nftTokenType := s.params["nft_token_type"]
	txHash := s.params["tx_hash"]

	targets = s.addQueriesWithWeights(
		targets,
		[]string{
			"/supply/total",
			fmt.Sprintf("/token/%s/supply", tokenContractID),
			fmt.Sprintf("/token/%s/token", tokenContractID),
			fmt.Sprintf("/token/%s/accounts/%s/balance", tokenContractID, from),
		},
		1,
	)
	targets = s.addQueriesWithWeights(
		targets,
		[]string{
			fmt.Sprintf("/coin/balances/%s", from),
			"/genesis/app_state/accounts",
			fmt.Sprintf("/token/%s/supply", tokenContractID),
			fmt.Sprintf("/collection/%s/fts/%s/supply", collectionContractID, ftTokenID),
			fmt.Sprintf("/collection/%s/tokens", collectionContractID),
			"/coin/tcony",
			fmt.Sprintf("/token/%s/supply", tokenContractID),
			fmt.Sprintf("/token/%s/token", tokenContractID),
			fmt.Sprintf("/token/%s/accounts/%s/balance", tokenContractID, from),
			fmt.Sprintf("/collection/%s/tokentypes", collectionContractID),
			fmt.Sprintf("/collection/%s/collection", collectionContractID),
			fmt.Sprintf("/collection/%s/tokentypes/%s/count", collectionContractID, nftTokenType),
			fmt.Sprintf("/collection/%s/fts/%s/mint", collectionContractID, ftTokenID),
			fmt.Sprintf("/collection/%s/fts/%s/burn", collectionContractID, ftTokenID),
			fmt.Sprintf("/collection/%s/fts/%s/supply", collectionContractID, ftTokenID),
		},
		3,
	)
	targets = s.addQueriesWithWeights(
		targets,
		[]string{
			fmt.Sprintf("/txs/%s", txHash),
			"/unconfirmed_txs",
			fmt.Sprintf("/collection/%s/nfts/%s/parent", collectionContractID, nftTokenID),
			fmt.Sprintf("/collection/%s/nfts/%s/root", collectionContractID, nftTokenID),
			fmt.Sprintf("/collection/%s/nfts/%s/children", collectionContractID, nftTokenID),
		},
		4,
	)
	targets = s.addQueriesWithWeights(
		targets,
		[]string{
			fmt.Sprintf("/collection/%s/tokentypes", collectionContractID),
			fmt.Sprintf("/collection/%s/collection", collectionContractID),
			fmt.Sprintf("/collection/%s/tokentypes/%s/count", collectionContractID, nftTokenType),
			fmt.Sprintf("/collection/%s/nfts/%s/count", collectionContractID, nftTokenID),
		},
		7,
	)
	return targets
}

func (s *TxAndQueryAllScenario) addQueriesWithWeights(targets []*vegeta.Target, urls []string,
	weight int) []*vegeta.Target {
	for i := 0; i < weight; i++ {
		for _, url := range urls {
			targets = append(targets, s.targetBuilder.MakeQueryTarget(url))
		}
	}
	return targets
}

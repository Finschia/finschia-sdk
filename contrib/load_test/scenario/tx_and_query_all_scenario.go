package scenario

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxAndQueryAllScenario struct {
	Info
}

func (s *TxAndQueryAllScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet, scenarioParams []string) ([]sdk.Msg, map[string]string, error) {
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

	numNFTPerUser := 13 * s.config.MsgsPerTxLoadTest
	mintNFTMsgs, err := GenerateMintNFTMsgs(masterAddress, hdWallet, s.config, collectionContractID, nftTokenType, numNFTPerUser)
	if err != nil {
		return nil, nil, err
	}

	msgs = append(msgs, grantTokenPermMsgs...)
	msgs = append(msgs, grantCollectionPermMsgs...)
	msgs = append(msgs, mintNFTMsgs...)

	return msgs, map[string]string{"token_contract_id": tokenContractID, "collection_contract_id": collectionContractID,
		"ft_token_id": ftTokenID, "nft_token_type": nftTokenType, "tx_hash": txHash,
		"num_nft_per_user": fmt.Sprintf("%d", numNFTPerUser)}, nil
}

func (s *TxAndQueryAllScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	account, err := s.linkService.GetAccount(keyWallet.Address().String())
	if err != nil {
		return nil, 0, err
	}

	numNFTPerUser, err := strconv.Atoi(s.stateParams["num_nft_per_user"])
	if err != nil {
		return nil, 0, err
	}
	nftTokenID := fmt.Sprintf("%s%08x", s.stateParams["nft_token_type"], numNFTPerUser*walletIndex+1)

	target, err := BuildTxTarget(s.Info, keyWallet, walletIndex, []string{"MsgSend", "MsgSend", "MsgIssue",
		"MsgMint", "MsgMint", "MsgTransfer", "MsgTransfer", "MsgModifyTokenName", "MsgModifyTokenURI", "MsgBurn",
		"MsgBurn", "MsgCreateCollection", "MsgApprove", "MsgIssueFT", "MsgMintFT", "MsgTransferFT", "MsgTransferFT",
		"MsgBurnFT", "MsgModifyCollection", "MsgIssueNFT", "MsgMintOneNFT", "MsgMintFiveNFT", "MsgAttach", "MsgDetach",
		"MsgTransferNFT", "MsgTransferNFT", "MsgMultiTransferNFT", "MsgMultiTransferNFT", "MsgBurnNFT"},
		s.scenarioParams)
	targets := []*vegeta.Target{target}
	targets = s.addQueryTargets(targets, account.Address, nftTokenID)

	return &targets, len(targets), err
}

func (s *TxAndQueryAllScenario) addQueryTargets(targets []*vegeta.Target, from sdk.AccAddress, nftTokenID string) []*vegeta.Target {
	tokenContractID := s.stateParams["token_contract_id"]
	collectionContractID := s.stateParams["collection_contract_id"]
	ftTokenID := s.stateParams["ft_token_id"]
	nftTokenType := s.stateParams["nft_token_type"]
	txHash := s.stateParams["tx_hash"]

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
			fmt.Sprintf("/collection/%s/nfts/%s", collectionContractID, nftTokenID),
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

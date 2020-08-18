package scenario

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxAllScenario struct {
	Info
}

func (s *TxAllScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet, scenarioParams []string) ([]sdk.Msg, map[string]string, error) {
	masterAddress := masterKeyWallet.Address()
	tokenContractID, _, err := IssueToken(s.linkService, s.txBuilder, masterKeyWallet)
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
		"ft_token_id": ftTokenID, "nft_token_type": nftTokenType,
		"num_nft_per_user": fmt.Sprintf("%d", numNFTPerUser)}, nil
}

func (s *TxAllScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	target, err := BuildTxTarget(s.Info, keyWallet, walletIndex, []string{"MsgSend", "MsgSend", "MsgIssue",
		"MsgMint", "MsgMint", "MsgTransfer", "MsgTransfer", "MsgModifyTokenName", "MsgModifyTokenURI", "MsgBurn",
		"MsgBurn", "MsgCreateCollection", "MsgApprove", "MsgIssueFT", "MsgMintFT", "MsgTransferFT", "MsgTransferFT",
		"MsgBurnFT", "MsgModifyCollection", "MsgIssueNFT", "MsgMintOneNFT", "MsgMintFiveNFT", "MsgAttach", "MsgDetach",
		"MsgTransferNFT", "MsgTransferNFT", "MsgMultiTransferNFT", "MsgMultiTransferNFT", "MsgBurnNFT"},
		s.scenarioParams)
	targets := []*vegeta.Target{target}

	return &targets, len(targets), err
}

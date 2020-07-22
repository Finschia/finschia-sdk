package scenario

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxCollectionScenario struct {
	Info
}

func (s *TxCollectionScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet, scenarioParams []string) ([]sdk.Msg, map[string]string, error) {
	masterAddress := masterKeyWallet.Address()
	contractID, err := CreateCollection(s.linkService, s.txBuilder, masterKeyWallet)
	if err != nil {
		return nil, nil, err
	}
	ftTokenID, err := IssueFT(s.linkService, s.txBuilder, masterKeyWallet, contractID)
	if err != nil {
		return nil, nil, err
	}
	nftTokenType, err := IssueNFT(s.linkService, s.txBuilder, masterKeyWallet, contractID)
	if err != nil {
		return nil, nil, err
	}

	msgs, err := GenerateRegisterAccountMsgs(masterAddress, hdWallet, s.config)
	if err != nil {
		return nil, nil, err
	}

	grantPermMsgs, err := GenerateGrantPermissionMsgs(masterAddress, hdWallet, s.config, contractID, "collection",
		[]string{"mint", "burn", "modify"})
	if err != nil {
		return nil, nil, err
	}
	numNFTPerUser := 2 * s.config.MsgsPerTxLoadTest
	mintNFTMsgs, err := GenerateMintNFTMsgs(masterAddress, hdWallet, s.config, contractID, nftTokenType, numNFTPerUser)
	if err != nil {
		return nil, nil, err
	}
	msgs = append(msgs, grantPermMsgs...)
	msgs = append(msgs, mintNFTMsgs...)

	return msgs, map[string]string{"collection_contract_id": contractID, "ft_token_id": ftTokenID,
		"nft_token_type": nftTokenType, "num_nft_per_user": fmt.Sprintf("%d", numNFTPerUser)}, nil
}

func (s *TxCollectionScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	target, err := BuildTxTarget(s.Info, keyWallet, walletIndex, []string{"MsgMintFT", "MsgTransferFT",
		"MsgModifyCollection", "MsgBurnFT", "MsgAttach", "MsgDetach", "MsgTransferNFT", "MsgBurnNFT"},
		s.scenarioParams)
	targets := []*vegeta.Target{target}
	return &targets, len(targets), err
}

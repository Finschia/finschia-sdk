package scenario

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxTransferNFTScenario struct {
	Info
}

func (s *TxTransferNFTScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet, scenarioParams []string) ([]sdk.Msg, map[string]string, error) {
	masterAddress := masterKeyWallet.Address()
	contractID, err := CreateCollection(s.linkService, s.txBuilder, masterKeyWallet)
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
	mintNFTMsgs, err := GenerateMintNFTMsgs(masterAddress, hdWallet, s.config, contractID, nftTokenType,
		s.config.MsgsPerTxLoadTest)
	if err != nil {
		return nil, nil, err
	}
	msgs = append(msgs, mintNFTMsgs...)

	return msgs, map[string]string{"collection_contract_id": contractID, "nft_token_type": nftTokenType,
		"num_nft_per_user": fmt.Sprintf("%d", s.config.MsgsPerTxLoadTest)}, nil
}

func (s *TxTransferNFTScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	target, err := BuildTxTarget(s.Info, keyWallet, walletIndex, []string{"MsgTransferNFT"}, s.scenarioParams)
	targets := []*vegeta.Target{target}
	return &targets, len(targets), err
}

package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxMintNFTScenario struct {
	Info
}

func (s *TxMintNFTScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
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
	grantPermMsgs, err := GenerateGrantPermissionMsgs(masterAddress, hdWallet, s.config, contractID, "collection",
		[]string{"mint"})
	if err != nil {
		return nil, nil, err
	}
	msgs = append(msgs, grantPermMsgs...)

	return msgs, map[string]string{"collection_contract_id": contractID, "nft_token_type": nftTokenType}, nil
}

func (s *TxMintNFTScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	target, err := BuildTxTarget(s.Info, keyWallet, walletIndex, []string{"MsgMintNFT"},
		s.scenarioParams)
	targets := []*vegeta.Target{target}
	return &targets, len(targets), err
}

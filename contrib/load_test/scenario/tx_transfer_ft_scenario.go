package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxTransferFTScenario struct {
	Info
}

func (s *TxTransferFTScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
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

	msgs, err := GenerateRegisterAccountMsgs(masterAddress, hdWallet, s.config)
	if err != nil {
		return nil, nil, err
	}
	mintFTMsgs, err := GenerateMintFTMsgs(masterAddress, hdWallet, s.config, contractID, ftTokenID)
	if err != nil {
		return nil, nil, err
	}
	msgs = append(msgs, mintFTMsgs...)

	return msgs, map[string]string{"collection_contract_id": contractID, "ft_token_id": ftTokenID}, nil
}

func (s *TxTransferFTScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	target, err := BuildTxTarget(s.Info, keyWallet, walletIndex, []string{"MsgTransferFT"}, s.scenarioParams)
	targets := []*vegeta.Target{target}
	return &targets, len(targets), err
}

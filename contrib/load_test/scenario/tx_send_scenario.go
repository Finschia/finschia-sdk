package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxSendScenario struct {
	Info
}

func (s *TxSendScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet, scenarioParams []string) ([]sdk.Msg, map[string]string, error) {
	msgs, err := GenerateRegisterAccountMsgs(masterKeyWallet.Address(), hdWallet, s.config)
	return msgs, nil, err
}

func (s *TxSendScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	target, err := BuildTxTarget(s.Info, keyWallet, walletIndex, []string{"MsgSend"}, s.scenarioParams)
	targets := []*vegeta.Target{target}
	return &targets, len(targets), err
}

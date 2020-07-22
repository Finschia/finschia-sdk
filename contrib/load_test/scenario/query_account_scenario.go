package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type QueryAccountScenario struct {
	Info
}

func (s *QueryAccountScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet, scenarioParams []string) ([]sdk.Msg, map[string]string, error) {
	msgs, err := GenerateRegisterAccountMsgs(masterKeyWallet.Address(), hdWallet, s.config)
	return msgs, nil, err
}

func (s *QueryAccountScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	targets := []*vegeta.Target{s.targetBuilder.MakeQueryTarget("/auth/accounts/" + keyWallet.Address().String())}
	return &targets, len(targets), nil
}

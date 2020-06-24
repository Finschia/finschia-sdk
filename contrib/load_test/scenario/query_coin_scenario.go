package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type QueryCoinScenario struct {
	linkService   *service.LinkService
	targetBuilder *TargetBuilder
	config        types.Config
}

func (s *QueryCoinScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet) ([]sdk.Msg, map[string]string, error) {
	return nil, nil, nil
}

func (s *QueryCoinScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	targets := []*vegeta.Target{s.targetBuilder.MakeQueryTarget("/coin/cony")}
	return &targets, len(targets), nil
}

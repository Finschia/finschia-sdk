package scenario

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type QueryBlockScenario struct {
	Info
}

func (s *QueryBlockScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet, scenarioParams []string) ([]sdk.Msg, map[string]string, error) {
	block, err := s.linkService.GetLatestBlock()
	if err != nil {
		return nil, nil, err
	}

	msgs, err := GenerateRegisterAccountMsgs(masterKeyWallet.Address(), hdWallet, s.config)
	if err != nil {
		return nil, nil, err
	}

	return msgs, map[string]string{"height": strconv.FormatInt(block.Block.Header.Height, 10)}, nil
}

func (s *QueryBlockScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	targets := []*vegeta.Target{s.targetBuilder.MakeQueryTarget(fmt.Sprintf("/blocks_with_tx_results/%s?fetchsize=%d",
		s.stateParams["height"], 3))}
	return &targets, len(targets), nil
}

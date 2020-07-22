package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type QuerySimulateScenario struct {
	Info
}

func (s *QuerySimulateScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet, scenarioParams []string) ([]sdk.Msg, map[string]string, error) {
	msgType := scenarioParams[0]
	var scenario Scenario
	switch msgType {
	case "MsgSend":
		scenario = &TxSendScenario{s.Info}
	case "MsgMintNFT":
		scenario = &TxMintNFTScenario{s.Info}
	case "MsgTransferFT":
		scenario = &TxTransferFTScenario{s.Info}
	case "MsgTransferNFT":
		scenario = &TxTransferNFTScenario{s.Info}
	default:
		return nil, nil, types.InvalidScenarioParameterError(msgType)
	}

	return scenario.GenerateStateSettingMsgs(masterKeyWallet, hdWallet, scenarioParams)
}

func (s *QuerySimulateScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	stdTx, err := BuildStdTx(s.Info, keyWallet, walletIndex, []string{s.scenarioParams[0]}, s.scenarioParams[1:])
	if err != nil {
		return nil, 0, err
	}

	target, err := s.targetBuilder.MakeQuerySimulateTarget(stdTx, service.BroadcastSync)
	targets := []*vegeta.Target{target}
	return &targets, len(targets), err
}

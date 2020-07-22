package scenario

import (
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Scenarios map[string]Scenario

type Scenario interface {
	GenerateStateSettingMsgs(*wallet.KeyWallet, *wallet.HDWallet, []string) ([]sdk.Msg, map[string]string, error)
	GenerateTarget(*wallet.KeyWallet, int) (*[]*vegeta.Target, int, error)
}

type Info struct {
	linkService    *service.LinkService
	targetBuilder  *TargetBuilder
	txBuilder      transaction.TxBuilderWithoutKeybase
	config         types.Config
	scenarioParams []string
	stateParams    map[string]string
}

func NewScenarios(config types.Config, params map[string]string, scenarioParams []string) Scenarios {
	scenarios := make(Scenarios)
	info := Info{
		linkService:    service.NewLinkService(&http.Client{}, app.MakeCodec(), config.TargetURL),
		targetBuilder:  NewTargetBuilder(app.MakeCodec(), config.TargetURL),
		txBuilder:      transaction.NewTxBuilder(uint64(config.MaxGasLoadTest)).WithChainID(config.ChainID),
		config:         config,
		scenarioParams: scenarioParams,
		stateParams:    params,
	}
	scenarios[types.QueryAccount] = &QueryAccountScenario{info}
	scenarios[types.QueryBlock] = &QueryBlockScenario{info}
	scenarios[types.QueryCoin] = &QueryCoinScenario{info}
	scenarios[types.QuerySimulate] = &QuerySimulateScenario{info}
	scenarios[types.TxSend] = &TxSendScenario{info}
	scenarios[types.TxEmpty] = &TxEmptyScenario{info}
	scenarios[types.TxMintNFT] = &TxMintNFTScenario{info}
	scenarios[types.TxTransferFT] = &TxTransferFTScenario{info}
	scenarios[types.TxTransferNFT] = &TxTransferNFTScenario{info}
	scenarios[types.TxToken] = &TxTokenScenario{info}
	scenarios[types.TxCollection] = &TxCollectionScenario{info}
	scenarios[types.TxAndQueryAll] = &TxAndQueryAllScenario{info}

	return scenarios
}

func GetPrepareBroadcastMode(scenario string) string {
	if scenario == types.QueryBlock {
		return service.BroadcastSync
	}
	return service.BroadcastBlock
}

func GetNumTargets(scenario string) int {
	if scenario == types.TxAndQueryAll {
		return 98
	}
	return 1
}

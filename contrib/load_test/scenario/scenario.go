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
	GenerateStateSettingMsgs(*wallet.KeyWallet, *wallet.HDWallet) ([]sdk.Msg, map[string]string, error)
	GenerateTarget(*wallet.KeyWallet, int) (*[]*vegeta.Target, int, error)
}

func NewScenarios(config types.Config, params map[string]string) Scenarios {
	scenarios := make(Scenarios)
	targetBuilder := NewTargetBuilder(app.MakeCodec(), config.TargetURL)
	txBuilder := transaction.NewTxBuilder(uint64(config.MaxGasLoadTest)).WithChainID(config.ChainID)
	linkService := service.NewLinkService(&http.Client{}, app.MakeCodec(), config.TargetURL)
	scenarios[types.QueryAccount] = &QueryAccountScenario{linkService, targetBuilder, config}
	scenarios[types.QueryBlock] = &QueryBlockScenario{linkService, targetBuilder, config, params}
	scenarios[types.QueryCoin] = &QueryCoinScenario{linkService, targetBuilder, config}
	scenarios[types.TxSend] = &TxSendScenario{linkService, targetBuilder, txBuilder, config}
	scenarios[types.TxEmpty] = &TxEmptyScenario{linkService, targetBuilder, txBuilder, config}
	scenarios[types.TxMintNFT] = &TxMintNFTScenario{linkService, targetBuilder, txBuilder, config, params}
	scenarios[types.TxMultiMintNFT] = &TxMultiMintNFTScenario{linkService, targetBuilder, txBuilder, config, params}
	scenarios[types.TxTransferFT] = &TxTransferFTScenario{linkService, targetBuilder, txBuilder, config, params}
	scenarios[types.TxTransferNFT] = &TxTransferNFTScenario{linkService, targetBuilder, txBuilder, config, params}
	scenarios[types.TxToken] = &TxTokenScenario{linkService, targetBuilder, txBuilder, config, params}
	scenarios[types.TxCollection] = &TxCollectionScenario{linkService, targetBuilder, txBuilder, config, params}
	scenarios[types.TxAndQueryAll] = &TxAndQueryAllScenario{linkService, targetBuilder, txBuilder, config, params}

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

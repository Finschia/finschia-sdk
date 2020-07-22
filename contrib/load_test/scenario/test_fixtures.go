package scenario

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	linktypes "github.com/line/link/types"
	"github.com/stretchr/testify/require"
)

func GivenTestEnvironments(t *testing.T, url string, scenarioType string, stateParams map[string]string,
	scenarioParams []string) (Scenario, *wallet.HDWallet, *wallet.KeyWallet) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet),
		linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Master Key Wallet
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	require.NoError(t, err)
	keyWallet, err := hdWallet.GetKeyWallet(1, 0)
	require.NoError(t, err)
	// Given Config
	config := types.Config{
		MsgsPerTxPrepare:  tests.TestMsgsPerTxPrepare,
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		TargetURL:         url,
		Duration:          tests.TestDuration,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
	}
	// Given Scenario
	scenarios := NewScenarios(config, stateParams, scenarioParams)

	return scenarios[scenarioType], hdWallet, keyWallet
}

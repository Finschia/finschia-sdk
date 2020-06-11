// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	sdktests "github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/line/link/app"
	clitest "github.com/line/link/cli_test"
	"github.com/line/link/contrib/load_test/cli"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	err := os.Chdir(fmt.Sprintf("%s/..", basepath))
	if err != nil {
		panic(err)
	}
}

const (
	keyFoo                = "foo"
	TestCoinName          = "stake"
	TestCoinAmount        = 50000000
	localhost             = "http://localhost"
	slavePort             = 8000
	lcdPort               = 1317
	TestDuration          = 10
	TestNumPrepareRequest = 10
	ExpectedAttackCount   = (TestDuration-tests.TestRampUpTime/2)*tests.TestTPS + TestDuration
)

var (
	masterAddress, _ = getMasterAddress()
	mutex            = &sync.Mutex{}
)

func TestLinkLoadTester(t *testing.T) {
	testCases := []struct {
		name       string
		targetType string
		isTx       bool
	}{
		{"QueryAccount", types.QueryAccount, false},
		{"TxSend", types.TxSend, true},
	}
	for i, tt := range testCases {
		tt := tt
		slavePort := slavePort + i
		lcdPort := lcdPort + i
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Log(tt.name)
			// Start LINK and LCD
			f := clitest.InitFixtures(t)
			proc := f.LDStart()
			restProc, err := f.RestServerStart(lcdPort, fmt.Sprintf("--node=%s:%s", localhost, f.Port))
			require.NoError(t, err)
			defer func() { require.NoError(t, proc.Stop(false)) }()
			defer func() { require.NoError(t, restProc.Stop(false)) }()
			defer f.Cleanup()
			// Run load generator
			cmd := cli.RunSlaveCmd()
			require.NoError(t, cmd.Flags().Set(cli.FlagPort, fmt.Sprintf("%d", slavePort)))
			go func() {
				require.NoError(t, cmd.RunE(cmd, nil))
			}()
			// Send enough coins to master account
			f.LogResult(f.TxSend(keyFoo, masterAddress, sdk.NewCoin(TestCoinName, sdk.NewInt(TestCoinAmount)), "-y"))
			// Given config
			mutex.Lock()
			assert.NoError(t, setConfig(tt.targetType, f.ChainID, lcdPort, slavePort))

			// When
			cmd = cli.PrepareCmd()
			assert.NoError(t, cmd.RunE(cmd, nil))
			cmd = cli.StartCmd()
			assert.NoError(t, cmd.RunE(cmd, nil))
			mutex.Unlock()
			sdktests.WaitForNextHeightTM(f.Port)

			// Then check the number of prepare txs
			if tt.isTx {
				txsPage := f.QueryTxs(1, 100, fmt.Sprintf("--tags='message.sender:%s'", masterAddress.String()))
				require.Equal(t, TestNumPrepareRequest, txsPage.Count)
				var preparedHeight int64
				for _, tx := range txsPage.Txs {
					require.Len(t, tx.Logs, tests.TestMsgsPerTxPrepare)
					require.Equal(t, uint32(0), tx.Code)
					preparedHeight = tx.Height
				}
				// Then check the number of generated txs
				latestHeight := f.QueryLatestBlock().Block.Height
				totalTxs := 0
				cdc := app.MakeCodec()
				for h := preparedHeight + 1; h <= latestHeight; h++ {
					txs := f.QueryBlockWithHeight(int(h)).Block.Txs
					// Then check the number of msgs per tx
					for _, txBytes := range txs {
						var tx auth.StdTx
						cdc.MustUnmarshalBinaryLengthPrefixed(txBytes, &tx)
						require.Len(t, tx.Msgs, tests.TestMsgsPerTxLoadTest)
					}
					totalTxs += len(txs)
				}
				require.InDelta(t, ExpectedAttackCount, totalTxs, tests.TestMaxAttackDifference)
			}
		})
	}
}

func getMasterAddress() (sdk.AccAddress, error) {
	hdWallet, err := wallet.NewHDWallet(tests.TestMasterMnemonic)
	if err != nil {
		return nil, err
	}
	masterKeyWallet, err := hdWallet.GetKeyWallet(0, 0)
	if err != nil {
		return nil, err
	}
	return masterKeyWallet.Address(), nil
}

func setConfig(targetType, chainID string, lcdPort, slavePort int) error {
	viper.Set(cli.FlagMasterMnemonic, tests.TestMasterMnemonic)
	viper.Set(cli.FlagTargetURL, fmt.Sprintf("%s:%d", localhost, lcdPort))
	viper.Set(cli.FlagChainID, chainID)
	viper.Set(cli.FlagCoinName, TestCoinName)
	viper.Set(cli.FlagMaxWorkers, tests.TestMaxWorkers)
	viper.Set(cli.FlagMsgsPerTxPrepare, tests.TestMsgsPerTxPrepare)
	viper.Set(cli.FlagMsgsPerTxLoadTest, tests.TestMsgsPerTxLoadTest)
	viper.Set(cli.FlagTPS, tests.TestTPS)
	viper.Set(cli.FlagDuration, TestDuration)
	viper.Set(cli.FlagRampUpTime, tests.TestRampUpTime)

	slavesMap := make(map[string]types.Slave)
	slaveURL := fmt.Sprintf("%s:%d", localhost, slavePort)
	slavesMap["slave1"] = types.NewSlave(slaveURL, tests.TestMnemonic, targetType)
	bytes, err := json.Marshal(slavesMap)
	if err != nil {
		return err
	}
	viper.Set(cli.FlagSlaves, string(bytes))
	return nil
}

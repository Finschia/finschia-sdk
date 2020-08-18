// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/line/link/app"
	clitest "github.com/line/link/cli_test"
	"github.com/line/link/contrib/load_test/cli"
	"github.com/line/link/contrib/load_test/master"
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
	keyFoo         = "foo"
	TestCoinName   = "stake"
	TestCoinAmount = 50000000
	localhost      = "http://localhost"
	slavePort      = 8000
	lcdPort        = 1317
	TestDuration   = 10
	TestTPS        = 50
)

var (
	masterAddress, _ = getMasterAddress()
	mutex            = &sync.Mutex{}
)

func TestLinkLoadTester(t *testing.T) {
	testCases := []struct {
		name                 string
		scenario             string
		scenarioParams       []string
		isTx                 bool
		tps                  int
		numPrepareTx         int
		numSingleMsgTx       int
		numMsgsPerTxLoadTest int
	}{
		{
			"QueryAccount",
			types.QueryAccount,
			[]string{},
			false,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration, tests.TestMsgsPerTxPrepare),
			0,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"QueryBlock",
			types.QueryBlock,
			[]string{},
			false,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration, tests.TestMsgsPerTxPrepare),
			0,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"QueryCoin",
			types.QueryCoin,
			[]string{},
			false,
			TestTPS,
			0,
			0,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"QuerySimulate_MsgSend",
			types.QuerySimulate,
			[]string{"MsgSend"},
			false,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration, tests.TestMsgsPerTxPrepare),
			0,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"QuerySimulate_MsgMintNFT",
			types.QuerySimulate,
			[]string{"MsgMintNFT", "1"},
			false,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration*2, tests.TestMsgsPerTxPrepare) + 2,
			2,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"QuerySimulate_MsgTransferFT",
			types.QuerySimulate,
			[]string{"MsgTransferFT"},
			false,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration*2, tests.TestMsgsPerTxPrepare) + 2,
			2,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"QuerySimulate_MsgTransferNFT",
			types.QuerySimulate,
			[]string{"MsgTransferNFT"},
			false,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration*(1+tests.TestMsgsPerTxLoadTest), tests.TestMsgsPerTxPrepare) + 2,
			2,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"QueryAll",
			types.QueryAll,
			[]string{},
			false,
			1,
			tests.GetNumPrepareTx(1*TestDuration*(8+13*tests.TestMsgsPerTxLoadTest), tests.TestMsgsPerTxPrepare) + 4,
			4,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"TxSend",
			types.TxSend,
			[]string{},
			true,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration, tests.TestMsgsPerTxPrepare),
			0,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"TxEmpty",
			types.TxEmpty,
			[]string{},
			true,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration, tests.TestMsgsPerTxPrepare),
			0,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"TxMintNFT",
			types.TxMintNFT,
			[]string{"1"},
			true,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration*2, tests.TestMsgsPerTxPrepare) + 2,
			2,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"TxMultiMintNFT",
			types.TxMintNFT,
			[]string{"5"},
			true,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration*2, tests.TestMsgsPerTxPrepare) + 2,
			2,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"TxTransferFT",
			types.TxTransferFT,
			[]string{},
			true,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration*2, tests.TestMsgsPerTxPrepare) + 2,
			2,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"TxTransferNFT",
			types.TxTransferNFT,
			[]string{},
			true,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration*(1+tests.TestMsgsPerTxLoadTest), tests.TestMsgsPerTxPrepare) + 2,
			2,
			tests.TestMsgsPerTxLoadTest,
		},
		{
			"TxToken",
			types.TxToken,
			[]string{},
			true,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration*4, tests.TestMsgsPerTxPrepare) + 1,
			1,
			5 * tests.TestMsgsPerTxLoadTest,
		},
		{
			"TxCollection",
			types.TxCollection,
			[]string{},
			true,
			TestTPS,
			tests.GetNumPrepareTx(TestTPS*TestDuration*(4+2*tests.TestMsgsPerTxLoadTest), tests.TestMsgsPerTxPrepare) + 3,
			3,
			8 * tests.TestMsgsPerTxLoadTest,
		},
		{
			"TxAll",
			types.TxAll,
			[]string{},
			true,
			1,
			tests.GetNumPrepareTx(1*TestDuration*(8+13*tests.TestMsgsPerTxLoadTest), tests.TestMsgsPerTxPrepare) + 4,
			4,
			29 * tests.TestMsgsPerTxLoadTest,
		},
	}
	for i, tt := range testCases {
		tt := tt
		slavePort := slavePort + i
		lcdPort := lcdPort + i
		t.Run(tt.name, func(t *testing.T) {
			defer tests.RemoveFile(master.ParamsFileName)
			defer tests.RemoveFile("TPS.png")
			defer tests.RemoveFile("Latency.png")
			defer tests.RemoveFile("result_data.json")
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
			log.Println(tt.name)
			requireNoErrorWithMutex(t, setConfig(tt.scenario, f.ChainID, lcdPort, slavePort, tt.tps, tt.scenarioParams))
			// Given buffer that can capture stdout
			origin, w, outC := captureStdout()

			// When
			cmd = cli.PrepareCmd()
			requireNoErrorWithMutex(t, cmd.RunE(cmd, nil))
			cmd = cli.StartCmd()
			requireNoErrorWithMutex(t, cmd.RunE(cmd, nil))
			time.Sleep(10 * time.Second)
			cmd = cli.ReportCmd()
			requireNoErrorWithMutex(t, cmd.RunE(cmd, nil))

			recoverStdout(origin, w)
			out := <-*outC
			fmt.Println(out)
			// Then there is no missing tx
			r, _ := regexp.Compile("(?:Num Missing Txs: )([0-9]+)")
			requireEqualWithMutex(t, "0", r.FindStringSubmatch(out)[1])

			// Then there is no failed tx in blocks
			r, _ = regexp.Compile("(?:Num Failed Tx Logs: )([0-9]+)")
			numFailedTxLogs := r.FindAllStringSubmatch(out, -1)
			if len(numFailedTxLogs) == 2 {
				requireEqualWithMutex(t, "0", numFailedTxLogs[1][1])
			}
			mutex.Unlock()

			// Then check the number of prepare txs
			txsPage := f.QueryTxs(1, 100, fmt.Sprintf("--tags='message.sender:%s'", masterAddress.String()))
			require.Equal(t, tt.numPrepareTx, txsPage.Count)
			multiMsgsTxs := txsPage.Txs[tt.numSingleMsgTx:]
			var preparedHeight int64
			for i, tx := range multiMsgsTxs {
				if i < len(multiMsgsTxs)-1 {
					require.Len(t, tx.Logs, tests.TestMsgsPerTxPrepare)
				}
				require.Equal(t, uint32(0), tx.Code)
				preparedHeight = tx.Height
			}
			if tt.isTx {
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
						require.Len(t, tx.Msgs, tt.numMsgsPerTxLoadTest)
					}
					totalTxs += len(txs)
				}
				ExpectedAttackCount := (TestDuration-tests.TestRampUpTime/2)*tt.tps + TestDuration
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

func setConfig(scenario, chainID string, lcdPort, slavePort, tps int, scenarioParams []string) error {
	viper.Set(cli.FlagMasterMnemonic, tests.TestMasterMnemonic)
	viper.Set(cli.FlagPrepareTargetURL, fmt.Sprintf("%s:%d", localhost, lcdPort))
	viper.Set(cli.FlagLoadTargetURL, fmt.Sprintf("%s:%d", localhost, lcdPort))
	viper.Set(cli.FlagChainID, chainID)
	viper.Set(cli.FlagCoinName, TestCoinName)
	viper.Set(cli.FlagMaxWorkers, tests.TestMaxWorkers)
	viper.Set(cli.FlagMsgsPerTxPrepare, tests.TestMsgsPerTxPrepare)
	viper.Set(cli.FlagMaxGasPrepare, tests.TestMaxGasPrepare)
	viper.Set(cli.FlagMsgsPerTxLoadTest, tests.TestMsgsPerTxLoadTest)
	viper.Set(cli.FlagMaxGasLoadTest, tests.TestMaxGasLoadTest)
	viper.Set(cli.FlagTPS, tps)
	viper.Set(cli.FlagDuration, TestDuration)
	viper.Set(cli.FlagRampUpTime, tests.TestRampUpTime)
	viper.Set(cli.FlagOutputDir, ".")
	viper.Set(cli.FlagTPSThreshold, -1)
	viper.Set(cli.FlagLatencyThreshold, -1)

	slavesMap := make(map[string]types.Slave)
	slaveURL := fmt.Sprintf("%s:%d", localhost, slavePort)
	slavesMap["slave1"] = types.NewSlave(slaveURL, tests.TestMnemonic, scenario, scenarioParams)
	bytes, err := json.Marshal(slavesMap)
	if err != nil {
		return err
	}
	viper.Set(cli.FlagSlaves, string(bytes))
	return nil
}

func requireNoErrorWithMutex(t *testing.T, err error) {
	if assert.NoError(t, err) {
		return
	}
	defer mutex.Unlock()
	t.FailNow()
}

func requireEqualWithMutex(t *testing.T, expected, actual interface{}) {
	if assert.Equal(t, expected, actual) {
		return
	}
	defer mutex.Unlock()
	t.FailNow()
}

func captureStdout() (*os.File, *os.File, *chan string) {
	originStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	return originStdout, w, &outC
}

func recoverStdout(originStdout *os.File, w *os.File) {
	w.Close()
	os.Stdout = originStdout
}

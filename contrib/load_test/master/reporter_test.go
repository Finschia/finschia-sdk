package master

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	linktypes "github.com/line/link/types"
	"github.com/stretchr/testify/require"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestReporter_Report(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given config
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		RampUpTime:        tests.TestRampUpTime,
		TargetURL:         server.URL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		MaxWorkers:        tests.TestMaxWorkers,
	}
	// Given Slaves
	slaves := []types.Slave{
		types.NewSlave(server.URL, tests.TestMnemonic, types.TxSend, []string{}),
		types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend, []string{}),
	}
	// Given Results
	controller := NewController(slaves, config, nil)
	require.NoError(t, controller.StartLoadTest())
	// Given Reporter
	startHeight := int64(1)
	endHeight := int64(3)
	thresholds := types.Thresholds{
		Latency: -1,
		TPS:     -1,
	}
	reporter := NewReporter(controller.Results, slaves, config, startHeight, endHeight, thresholds)

	// When
	require.NoError(t, reporter.Report(""))

	// Then
	require.Len(t, reporter.metrics.blockMetrics.tps, int(endHeight-startHeight+1))
	require.Len(t, reporter.metrics.blockMetrics.height, int(endHeight-startHeight+1))
	require.Len(t, reporter.metrics.blockMetrics.blockInterval, int(endHeight-startHeight+1))
	require.Len(t, reporter.metrics.blockMetrics.numTxs, int(endHeight-startHeight+1))
	require.Equal(t, uint64(len(slaves)), reporter.metrics.Requests)
	require.Len(t, reporter.metrics.latencies, len(slaves))
	require.Len(t, reporter.metrics.timeStamps, len(slaves))
	require.Equal(t, float64(1), reporter.metrics.Success)
}

func TestReporter_getBlockMetrics(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given config
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		RampUpTime:        tests.TestRampUpTime,
		TargetURL:         server.URL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		MaxWorkers:        tests.TestMaxWorkers,
	}
	// Given Slaves
	slaves := []types.Slave{
		types.NewSlave(server.URL, tests.TestMnemonic, types.QueryAccount, []string{}),
		types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend, []string{}),
	}
	// Given Reporter
	r := &Reporter{
		metrics:     NewMetrics(),
		linkService: service.NewLinkService(&http.Client{}, app.MakeCodec(), config.TargetURL),
		slaves:      slaves,
		config:      config,
		startHeight: 3,
		endHeight:   10,
	}

	// When
	require.NoError(t, r.getBlockMetrics())

	// Then
	require.Equal(t, 1, mock.GetCallCounter(server.URL).QueryBlockCallCount)
	require.Equal(t, int(r.endHeight-r.startHeight+2), mock.GetCallCounter(server.URL).QueryBlocksWithTxResultsCallCount)
	require.Len(t, r.metrics.blockMetrics.height, int(r.endHeight-r.startHeight+1))
	require.Len(t, r.metrics.blockMetrics.numTxs, int(r.endHeight-r.startHeight+1))
	require.Len(t, r.metrics.blockMetrics.blockInterval, int(r.endHeight-r.startHeight+1))
	require.Len(t, r.metrics.blockMetrics.tps, int(r.endHeight-r.startHeight+1))
}

func TestReporter_checkMissingTx(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given config
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		RampUpTime:        tests.TestRampUpTime,
		TargetURL:         server.URL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		MaxWorkers:        tests.TestMaxWorkers,
	}
	// Given Slaves
	slaves := []types.Slave{
		types.NewSlave(server.URL, tests.TestMnemonic, types.QueryAccount, []string{}),
		types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend, []string{}),
	}
	// Given Reporter
	m := NewMetrics()
	m.transferredTxs = []string{
		"0119A8EE45E43090478BF7CBB593FE3EF83A9C8E96A4477383A251CE0A6BD38C",
		"011AD934A6253FFE119D9C7E0B76DB783A8FB87BFA39F531FB8BFC6AEC6BF4E0",
		"011B07538DA3AE5AA1C2ACD83ED36DA1B64F0DC8D21D3517224EE4383CF463A7",
		"011BB3B34F4B59B4E1143151BA0726A8C1A52D471ED1181A3940012472FE8253",
		"011C5CD5200C93C06A0C392D89DD0A2D978FBB58BA7193D58580B9C614D35D91",
		"011EB319B0CA6923D4AEE54E82CCAF71991D2E52B51FB71C6FAFBED33E21D66D",
	}

	r := &Reporter{
		metrics:     m,
		linkService: service.NewLinkService(&http.Client{}, app.MakeCodec(), server.URL),
		slaves:      slaves,
		config:      config,
		startHeight: 3,
		endHeight:   10,
	}

	// When
	require.NoError(t, r.checkMissingTx())

	// Then
	require.Equal(t, 6, mock.GetCallCounter(server.URL).QueryTxCallCount)
	require.Len(t, r.metrics.missingTxs, 0)
}

func TestReporter_printReport(t *testing.T) {
	// Given config
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		RampUpTime:        tests.TestRampUpTime,
		TargetURL:         tests.TestTargetURL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		MaxWorkers:        tests.TestMaxWorkers,
	}
	// Given Slaves
	slaves := []types.Slave{
		types.NewSlave(tests.TestTargetURL, tests.TestMnemonic, types.TxSend, []string{}),
		types.NewSlave(tests.TestTargetURL, tests.TestMnemonic2, types.TxSend, []string{}),
	}
	// Given Reporter
	m := &Metrics{
		Metrics: vegeta.Metrics{
			Latencies: vegeta.LatencyMetrics{
				Total: 30 * time.Millisecond,
				Mean:  3 * time.Millisecond,
				P50:   2 * time.Millisecond,
				P90:   3 * time.Millisecond,
				P95:   4 * time.Millisecond,
				P99:   5 * time.Millisecond,
				Max:   6 * time.Millisecond,
				Min:   1 * time.Millisecond,
			},
			BytesIn: vegeta.ByteMetrics{
				Total: 1050,
				Mean:  105,
			},
			BytesOut: vegeta.ByteMetrics{
				Total: 10290,
				Mean:  1029,
			},
			Earliest:    time.Unix(0, 0),
			Latest:      time.Unix(1, 0),
			End:         time.Unix(2, 0),
			Duration:    30 * time.Millisecond,
			Wait:        10 * time.Millisecond,
			Requests:    10,
			Rate:        300,
			Throughput:  210,
			Success:     7,
			StatusCodes: map[string]int{"200": 7, "500": 3},
			Errors:      []string{"Error Log", "Error Log", "Error Log", "Error Log2", "Error Log3"},
		},
		latencies: []float64{float64(1 * time.Millisecond), float64(2 * time.Millisecond),
			float64(3 * time.Millisecond), float64(4 * time.Millisecond), float64(5 * time.Millisecond),
			float64(6 * time.Millisecond), float64(3 * time.Millisecond), float64(4 * time.Millisecond),
			float64(3 * time.Millisecond), float64(4 * time.Millisecond)},
		timeStamps: []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0), time.Unix(4, 0),
			time.Unix(5, 0), time.Unix(6, 0), time.Unix(7, 0), time.Unix(8, 0), time.Unix(9, 0), time.Unix(10, 0)},
		txSize:     181,
		maxUsedGas: 70000,
		minUsedGas: 50000,
		blockMetrics: BlockMetrics{
			height:        []float64{1, 2, 3, 4, 5, 6},
			numTxs:        []float64{100, 120, 130, 120, 90, 120},
			blockInterval: []float64{5, 3, 2, 1, 3, 4},
			blockGas:      []float64{100000, 200000, 300000, 400000, 500000, 200000},
			round:         []float64{0, 1, 0, 0, 2, 1},
			tps:           []float64{20, 40, 75, 120, 30, 30},
			numFailedTxs:  5,
			failedTxLogs:  []string{"Error Log", "Error Log", "Error Log", "Error Log2", "Error Log2"},
		},
		failedTxLogs: []string{"Error Log", "Error Log", "Error Log", "Error Log2", "Error Log2"},
		missingTxs: []string{
			"0119A8EE45E43090478BF7CBB593FE3EF83A9C8E96A4477383A251CE0A6BD38C",
			"011AD934A6253FFE119D9C7E0B76DB783A8FB87BFA39F531FB8BFC6AEC6BF4E0",
			"011B07538DA3AE5AA1C2ACD83ED36DA1B64F0DC8D21D3517224EE4383CF463A7",
		},
	}
	reporter := &Reporter{
		metrics:     m,
		linkService: service.NewLinkService(&http.Client{}, app.MakeCodec(), config.TargetURL),
		slaves:      slaves,
		config:      config,
		startHeight: 1,
		endHeight:   6,
	}

	var tests = []struct {
		reportFunc     func(io.Writer) error
		expectedReport string
	}{
		{
			reporter.printTxReport,
			`----------- Tx Report

[Throughput]         
Throughput(total)    300.00/sec
Throughput(success)  210.00/sec

[Latencies]    
Average        3ms
Std.Dev        1433720.88
Max            6ms
99 Percentile  5ms
95 Percentile  4ms
90 Percentile  3ms
Median         2ms
Min            1ms

[Bytes]            
Total Bytes In     1050 bytes
Average Bytes In   105.00 bytes
Total Bytes Out    10290 bytes
Average Bytes Out  1029.00 bytes

[Summary]       
Total Requests  10
Success         700.00%
Duration        [total, attack, wait]  40ms, 30ms, 10ms
Status Codes    [code:count]           200:7  500:3  
Num Error Set: 5
Error Set:
Error Log (3 logs)
Error Log2
Error Log3

Num Failed Tx Logs: 5
Failed Tx Logs:
Error Log (3 logs)
Error Log2 (2 logs)

Num Missing Txs: 3
Missing Txs:
0119A8EE45E43090478BF7CBB593FE3EF83A9C8E96A4477383A251CE0A6BD38C
011AD934A6253FFE119D9C7E0B76DB783A8FB87BFA39F531FB8BFC6AEC6BF4E0
011B07538DA3AE5AA1C2ACD83ED36DA1B64F0DC8D21D3517224EE4383CF463A7
`,
		},
		{
			reporter.printBlockReport,
			`
----------- Block Report

[TPS]    
Average  52.50/sec
Std.Dev  38.18
MPS      157.50/sec

[Block Interval]  
Average           3.00
Std.Dev           1.41

[Num Txs Per Block]  
Average              113.33
Std.Dev              15.06

[Failed Txs]     
Total Failed Tx  5

Num Failed Tx Logs: 5
Failed Tx Logs:
Error Log (3 logs)
Error Log2 (2 logs)

[Block Details]
+--------------+---------+----------------+----------------+-------+
| BLOCK HEIGHT | NUM TXS | BLOCK USED GAS | BLOCK INTERVAL | ROUND |
+--------------+---------+----------------+----------------+-------+
|            1 |     100 |         100000 |              5 |     0 |
|            2 |     120 |         200000 |              3 |     1 |
|            3 |     130 |         300000 |              2 |     0 |
|            4 |     120 |         400000 |              1 |     0 |
|            5 |      90 |         500000 |              3 |     2 |
|            6 |     120 |         200000 |              4 |     1 |
+--------------+---------+----------------+----------------+-------+
`,
		},
		{
			reporter.printMsgReport,
			`
----------- Msg Report

+---------+-----------------+-----------------+--------------+--------------+------------------+
|   MSG   | NUM MSGS PER TX | TX SIZE (BYTES) | MAX USED GAS | MIN USED GAS | AVERAGE USED GAS |
+---------+-----------------+-----------------+--------------+--------------+------------------+
| tx_send |               3 |             181 |        70000 |        50000 |      2500.000000 |
+---------+-----------------+-----------------+--------------+--------------+------------------+
`,
		},
	}
	for _, tt := range tests {
		var b bytes.Buffer
		require.NoError(t, tt.reportFunc(&b))

		require.Equal(t, tt.expectedReport, b.String())
	}
}

func TestReporter_checkThresholds(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given config
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		RampUpTime:        tests.TestRampUpTime,
		TargetURL:         server.URL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		MaxWorkers:        tests.TestMaxWorkers,
	}
	// Given Slaves
	slaves := []types.Slave{
		types.NewSlave(server.URL, tests.TestMnemonic, types.QueryAccount, []string{}),
		types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend, []string{}),
	}
	// Given Reporter
	r := &Reporter{
		metrics:     NewMetrics(),
		linkService: service.NewLinkService(&http.Client{}, app.MakeCodec(), config.TargetURL),
		slaves:      slaves,
		config:      config,
		startHeight: 3,
		endHeight:   10,
	}
	r.metrics.Throughput = 250
	r.metrics.Latencies.Mean = 300
	r.metrics.blockMetrics.tps = []float64{15.44}

	var tests = []struct {
		latencyThreshold    time.Duration
		tpsThreshold        int
		throughputThreshold int
		errCheckFunc        func(require.TestingT, error, ...interface{})
	}{
		{-1, -1, -1, require.NoError},
		{300, 15, 250, require.NoError},
		{300, 15, 251, require.Error},
		{299, 15, 250, require.Error},
		{300, 16, 250, require.Error},
	}

	for i, tt := range tests {
		t.Logf("Test %d", i)
		// When
		r.thresholds = types.Thresholds{
			Latency:    tt.latencyThreshold,
			TPS:        tt.tpsThreshold,
			Throughput: tt.throughputThreshold,
		}

		// Then
		tt.errCheckFunc(t, r.checkThresholds())
	}
}

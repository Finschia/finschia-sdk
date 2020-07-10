package master

import (
	"time"

	"github.com/line/link/app"
	"github.com/line/link/x/account/client/types"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type BlockMetrics struct {
	height        []float64
	numTxs        []float64
	blockInterval []float64
	blockGas      []float64
	round         []float64
	tps           []float64
	numFailedTxs  int
	failedTxLogs  []string
}

type Metrics struct {
	vegeta.Metrics
	latencies  []float64
	timeStamps []time.Time

	txSize         int
	maxUsedGas     int64
	minUsedGas     int64
	blockMetrics   BlockMetrics
	failedTxLogs   []string
	transferredTxs []string
	processedTxs   map[string]bool
	missingTxs     []string
}

func NewMetrics() *Metrics {
	return &Metrics{
		Metrics:      vegeta.Metrics{},
		latencies:    make([]float64, 0, 512),
		timeStamps:   make([]time.Time, 0, 512),
		blockMetrics: BlockMetrics{},
		processedTxs: make(map[string]bool),
	}
}

func (m *Metrics) AddResults(results []vegeta.Result) {
	for _, res := range results {
		res := res
		m.Add(&res)
	}
}

func (m *Metrics) Add(r *vegeta.Result) {
	m.Metrics.Add(r)
	m.latencies = append(m.latencies, r.Latency.Seconds())
	m.timeStamps = append(m.timeStamps, r.Timestamp)

	var res types.TxResponse
	if err := app.MakeCodec().UnmarshalJSON(r.Body, &res); err == nil {
		if res.Code != 0 {
			m.failedTxLogs = append(m.failedTxLogs, string(r.Body))
		} else if res.TxHash != "" {
			m.transferredTxs = append(m.transferredTxs, res.TxHash)
		}
	}
}

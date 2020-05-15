package master

import (
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type BlockMetrics struct {
	height        []float64
	numTxs        []float64
	blockInterval []float64
	tps           []float64
}

type Metrics struct {
	vegeta.Metrics
	latencies    []float64
	timeStamps   []time.Time
	blockMetrics BlockMetrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		Metrics:      vegeta.Metrics{},
		latencies:    make([]float64, 0, 512),
		timeStamps:   make([]time.Time, 0, 512),
		blockMetrics: BlockMetrics{},
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
}

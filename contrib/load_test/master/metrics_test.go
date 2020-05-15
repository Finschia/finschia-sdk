package master

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	TestAttack  = "TestAttack"
	TestLatency = time.Second
)

var TestTimestamp = time.Now()

func TestMetrics_Add(t *testing.T) {
	metrics := NewMetrics()
	res := &vegeta.Result{
		Attack:    TestAttack,
		Timestamp: TestTimestamp,
		Latency:   TestLatency,
	}

	metrics.Add(res)

	require.Equal(t, uint64(1), metrics.Metrics.Requests)
	require.Equal(t, TestTimestamp, metrics.timeStamps[0])
	require.Equal(t, TestLatency.Seconds(), metrics.latencies[0])
}

func TestMetrics_AddResults(t *testing.T) {
	metrics := NewMetrics()
	res := []vegeta.Result{
		{
			Attack:    TestAttack,
			Timestamp: TestTimestamp,
			Latency:   TestLatency,
		},
		{
			Attack:    TestAttack,
			Timestamp: TestTimestamp,
			Latency:   TestLatency,
		},
	}

	metrics.AddResults(res)

	require.Equal(t, uint64(2), metrics.Metrics.Requests)
	require.Len(t, metrics.latencies, 2)
	require.Len(t, metrics.timeStamps, 2)
	require.Equal(t, TestTimestamp, metrics.timeStamps[0])
	require.Equal(t, TestLatency.Seconds(), metrics.latencies[0])
}

package iavl

import (
	"testing"

	"github.com/go-kit/kit/metrics/discard"
	"github.com/stretchr/testify/require"
)

func TestPrometheusMetrics(t *testing.T) {
	metrics := PrometheusMetrics("test")
	require.NotEqual(t, metrics.IAVLCacheHits, discard.NewGauge())
	require.NotEqual(t, metrics.IAVLCacheMisses, discard.NewGauge())
	require.NotEqual(t, metrics.IAVLCacheEntries, discard.NewGauge())
	require.NotEqual(t, metrics.IAVLCacheBytes, discard.NewGauge())
}

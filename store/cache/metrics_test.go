package cache

import (
	"testing"

	"github.com/go-kit/kit/metrics/discard"
	"github.com/stretchr/testify/require"
)

func TestPrometheusMetrics(t *testing.T) {
	metrics := PrometheusMetrics("test")
	require.NotEqual(t, metrics.InterBlockCacheHits, discard.NewCounter())
	require.NotEqual(t, metrics.InterBlockCacheMisses, discard.NewCounter())
	require.NotEqual(t, metrics.InterBlockCacheEntries, discard.NewGauge())
	require.NotEqual(t, metrics.InterBlockCacheBytes, discard.NewGauge())
}

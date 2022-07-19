package baseapp

import (
	"testing"

	"github.com/go-kit/kit/metrics/discard"
	"github.com/stretchr/testify/require"
)

func TestMetricsProvider(t *testing.T) {
	p1, p2 := MetricsProvider(true)
	c1 := p1()
	c2 := p2()
	require.NotEqual(t, c1.InterBlockCacheHits, discard.NewCounter())
	require.NotEqual(t, c2.IAVLCacheHits, discard.NewGauge())

	p1, p2 = MetricsProvider(false)
	c1 = p1()
	c2 = p2()
	require.Equal(t, c1.InterBlockCacheHits, discard.NewCounter())
	require.Equal(t, c2.IAVLCacheHits, discard.NewGauge())
}

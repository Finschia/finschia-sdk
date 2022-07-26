package baseapp

import (
	"testing"

	"github.com/go-kit/kit/metrics/discard"
	"github.com/stretchr/testify/require"
)

func TestMetricsProvider(t *testing.T) {
	p1 := MetricsProvider(true)
	c1 := p1()
	require.NotEqual(t, c1.InterBlockCacheHits, discard.NewCounter())

	p1 = MetricsProvider(false)
	c1 = p1()
	require.Equal(t, c1.InterBlockCacheHits, discard.NewCounter())
}

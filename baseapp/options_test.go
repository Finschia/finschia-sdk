package baseapp

import (
	"testing"

	"github.com/go-kit/kit/metrics/discard"
	"github.com/stretchr/testify/require"
)

func TestMetricsProvider(t *testing.T) {
	p := MetricsProvider(true)
	c := p()
	require.NotEqual(t, c.InterBlockCacheHits, discard.NewCounter())

	p = MetricsProvider(false)
	c = p()
	require.Equal(t, c.InterBlockCacheHits, discard.NewCounter())
}

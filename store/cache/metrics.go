package cache

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this
	// package.
	MetricsSubsystem = "inter_block_cache"
)

// Metrics contains metrics exposed by this package.
type Metrics struct {
	InterBlockCacheHits    metrics.Counter
	InterBlockCacheMisses  metrics.Counter
	InterBlockCacheEntries metrics.Gauge
	InterBlockCacheBytes   metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
// Optionally, labels can be provided along with their values ("foo",
// "fooValue").
func PrometheusMetrics(namespace string, labelsAndValues ...string) *Metrics {
	labels := []string{}
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}
	return &Metrics{
		InterBlockCacheHits: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "hits",
			Help:      "Cache hits of the inter block cache",
		}, labels).With(labelsAndValues...),
		InterBlockCacheMisses: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "misses",
			Help:      "Cache misses of the inter block cache",
		}, labels).With(labelsAndValues...),
		InterBlockCacheEntries: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "entries",
			Help:      "Cache entry count of the inter block cache",
		}, labels).With(labelsAndValues...),
		InterBlockCacheBytes: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "bytes_size",
			Help:      "Cache bytes size of the inter block cache",
		}, labels).With(labelsAndValues...),
	}
}

// NopMetrics returns no-op Metrics.
func NopMetrics() *Metrics {
	return &Metrics{
		InterBlockCacheHits:    discard.NewCounter(),
		InterBlockCacheMisses:  discard.NewCounter(),
		InterBlockCacheEntries: discard.NewGauge(),
		InterBlockCacheBytes:   discard.NewGauge(),
	}
}

type MetricsProvider func() *Metrics

// PrometheusMetricsProvider returns PrometheusMetrics for each store
func PrometheusMetricsProvider(namespace string, labelsAndValues ...string) func() *Metrics {
	return func() *Metrics {
		return PrometheusMetrics(namespace, labelsAndValues...)
	}
}

// NopMetricsProvider returns NopMetrics for each store
func NopMetricsProvider() func() *Metrics {
	return NopMetrics
}

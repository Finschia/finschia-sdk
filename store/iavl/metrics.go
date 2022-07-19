package iavl

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this
	// package.
	MetricsSubsystem = "iavl_cache"
)

// Metrics contains metrics exposed by this package.
type Metrics struct {
	IAVLCacheHits    metrics.Gauge
	IAVLCacheMisses  metrics.Gauge
	IAVLCacheEntries metrics.Gauge
	IAVLCacheBytes   metrics.Gauge
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
		IAVLCacheHits: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "hits",
			Help:      "Cache hit count of the iavl cache",
		}, labels).With(labelsAndValues...),
		IAVLCacheMisses: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "misses",
			Help:      "Cache miss count of the iavl cache",
		}, labels).With(labelsAndValues...),
		IAVLCacheEntries: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "entries",
			Help:      "Cache entry count of the iavl cache",
		}, labels).With(labelsAndValues...),
		IAVLCacheBytes: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "bytes_size",
			Help:      "Cache bytes size of the iavl cache",
		}, labels).With(labelsAndValues...),
	}
}

// NopMetrics returns no-op Metrics.
func NopMetrics() *Metrics {
	return &Metrics{
		IAVLCacheHits:    discard.NewGauge(),
		IAVLCacheMisses:  discard.NewGauge(),
		IAVLCacheEntries: discard.NewGauge(),
		IAVLCacheBytes:   discard.NewGauge(),
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
	return func() *Metrics {
		return NopMetrics()
	}
}

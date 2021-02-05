package baseapp

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this
	// package.
	MetricsSubsystem = "baseapp"
)

// Metrics contains metrics exposed by this package.
type Metrics struct {
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
// Optionally, labels can be provided along with their values ("foo",
// "fooValue").
func PrometheusMetrics(namespace string, labelsAndValues ...string) *Metrics {
	return &Metrics{}
}

// NopMetrics returns no-op Metrics.
func NopMetrics() *Metrics {
	return &Metrics{}
}

func GeneratePrometheusMetrics(prometheus bool) *Metrics {
	if prometheus {
		return PrometheusMetrics("app")
	}
	return NopMetrics()
}

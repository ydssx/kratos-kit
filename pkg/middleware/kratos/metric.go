package kratos

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func MetricServer() middleware.Middleware {
	var (
		_metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "server",
			Subsystem: "requests",
			Name:      "duration_sec",
			Help:      "server requests duration(sec).",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1},
		}, []string{"kind", "operation"})

		_metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "client",
			Subsystem: "requests",
			Name:      "code_total",
			Help:      "The total number of processed requests",
		}, []string{"kind", "operation", "code", "reason"})
	)

	prometheus.MustRegister(_metricSeconds, _metricRequests)

	return metrics.Server(
		// metrics.WithSeconds(prometheus.NewHistogram(_metricSeconds)),
		// metrics.WithRequests(prometheus.NewCounter(_metricRequests)),
	)
}

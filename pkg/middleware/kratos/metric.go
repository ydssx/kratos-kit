package kratos

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func MetricServer() middleware.Middleware {
	meterProvider := sdkmetric.NewMeterProvider()
	_metricRequests, _ := metrics.DefaultRequestsCounter(meterProvider.Meter("server_requests_code_total"), "server_requests_code_total")
	_metricSeconds, _ := metrics.DefaultSecondsHistogram(meterProvider.Meter("server_requests_duration_seconds"), "server_requests_duration_seconds")

	return metrics.Server(
		metrics.WithRequests(_metricRequests),
		metrics.WithSeconds(_metricSeconds),
	)
}

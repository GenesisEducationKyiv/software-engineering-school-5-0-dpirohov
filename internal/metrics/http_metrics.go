package metrics

import "github.com/prometheus/client_golang/prometheus"

type HTTPMetrics struct {
	httpRequestDuration *prometheus.HistogramVec
}

func NewHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_requests_duration_seconds",
				Help:    "Histogram of response time for handler",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "route", "status"},
		),
	}
}

func (m *HTTPMetrics) Register(reg prometheus.Registerer) {
	reg.MustRegister(
		m.httpRequestDuration,
	)
}

func (m *HTTPMetrics) ObserveRequestDuration(method, route, status string, duration float64) {
	m.httpRequestDuration.WithLabelValues(method, route, status).Observe(duration)
}

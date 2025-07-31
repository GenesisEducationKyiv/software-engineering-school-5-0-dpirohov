package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type CacheMetrics struct {
	cacheResult         *prometheus.CounterVec
	lockWaitDurationSec prometheus.Histogram
}

func NewCacheMetrics() *CacheMetrics {
	return &CacheMetrics{
		cacheResult: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_total",
				Help: "Total number of cache accesses, labeled by result",
			},
			[]string{"result"},
		),
		lockWaitDurationSec: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "cache_lock_wait_seconds",
			Help:    "Time spent waiting for cache lock",
			Buckets: prometheus.DefBuckets,
		}),
	}
}

func (m *CacheMetrics) Register(reg prometheus.Registerer) {
	reg.MustRegister(
		m.cacheResult,
		m.lockWaitDurationSec,
	)
}

func (m *CacheMetrics) IncCacheHit() {
	m.cacheResult.WithLabelValues("hit").Inc()
}

func (m *CacheMetrics) IncCacheMiss() {
	m.cacheResult.WithLabelValues("miss").Inc()
}

func (m *CacheMetrics) ObserveLockWaitDuration(seconds float64) {
	m.lockWaitDurationSec.Observe(seconds)
}

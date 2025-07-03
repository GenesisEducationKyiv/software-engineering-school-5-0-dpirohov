package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type CacheMetrics struct {
	cacheHit            prometheus.Counter
	cacheMiss           prometheus.Counter
	lockWait            prometheus.Counter
	lockWaitDurationSec prometheus.Histogram
}

func NewCacheMetrics() *CacheMetrics {
	return &CacheMetrics{
		cacheHit: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_hit_total",
			Help: "Number of cache hits",
		}),
		cacheMiss: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_miss_total",
			Help: "Number of cache misses",
		}),
		lockWait: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_lock_wait_total",
			Help: "Number of times request waited for a cache lock",
		}),
		lockWaitDurationSec: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "cache_lock_wait_seconds",
			Help:    "Time spent waiting for cache lock",
			Buckets: prometheus.DefBuckets,
		}),
	}
}

func (m *CacheMetrics) Register(reg prometheus.Registerer) {
	reg.MustRegister(
		m.cacheHit,
		m.cacheMiss,
		m.lockWait,
		m.lockWaitDurationSec,
	)
}

func (m *CacheMetrics) IncCacheHit() {
	m.cacheHit.Inc()
}

func (m *CacheMetrics) IncCacheMiss() {
	m.cacheMiss.Inc()
}

func (m *CacheMetrics) IncLockWait() {
	m.lockWait.Inc()
}

func (m *CacheMetrics) ObserveLockWaitDuration(seconds float64) {
	m.lockWaitDurationSec.Observe(seconds)
}

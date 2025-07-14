package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"weatherApi/internal/metrics"

	"weatherApi/internal/dto"

	"github.com/redis/go-redis/v9"
)

type CacheRepoInterface interface {
	Get(ctx context.Context, city string) (*dto.WeatherResponse, error)
	Set(ctx context.Context, city string, data *dto.WeatherResponse) error
	AcquireLock(ctx context.Context, city string) (bool, error)
	WaitForUnlock(ctx context.Context, city string) (*dto.WeatherResponse, error)
	ReleaseLock(ctx context.Context, city string) error
}

var ErrCacheIsEmpty = errors.New("weather cache is empty")

type Repository struct {
	client       *redis.Client
	cacheTTL     time.Duration
	lockTTL      time.Duration
	lockRetryDur time.Duration
	lockMaxWait  time.Duration
	metrics      *metrics.CacheMetrics
}

type RepositoryOptions struct {
	Client       *redis.Client
	CacheTTL     time.Duration
	LockTTL      time.Duration
	LockRetryDur time.Duration
	LockMaxWait  time.Duration
	Metrics      *metrics.CacheMetrics
}

func NewWeatherRepository(options *RepositoryOptions) *Repository {
	return &Repository{
		client:       options.Client,
		cacheTTL:     options.CacheTTL,
		lockTTL:      options.LockTTL,
		lockRetryDur: options.LockRetryDur,
		lockMaxWait:  options.LockMaxWait,
		metrics:      options.Metrics,
	}
}

func (r *Repository) getCacheKey(city string) string {
	return fmt.Sprintf("weather:city:%s", city)
}

func (r *Repository) getLockKey(city string) string {
	return fmt.Sprintf("weather:lock:%s", city)
}

func (r *Repository) Get(ctx context.Context, city string) (*dto.WeatherResponse, error) {
	key := r.getCacheKey(city)

	data, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		r.metrics.IncCacheMiss()
		return nil, ErrCacheIsEmpty
	} else if err != nil {
		return nil, err
	}
	r.metrics.IncCacheHit()
	var res dto.WeatherResponse
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *Repository) Set(ctx context.Context, city string, data *dto.WeatherResponse) error {
	key := r.getCacheKey(city)

	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, raw, r.cacheTTL).Err()
}

func (r *Repository) AcquireLock(ctx context.Context, city string) (bool, error) {
	lockKey := r.getLockKey(city)
	ok, err := r.client.SetNX(ctx, lockKey, "1", r.lockTTL).Result()
	return ok, err
}

func (r *Repository) WaitForUnlock(ctx context.Context, city string) (*dto.WeatherResponse, error) {
	key := r.getCacheKey(city)
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(r.lockRetryDur):
			data, err := r.client.Get(ctx, key).Result()
			if err == nil {
				r.metrics.ObserveLockWaitDuration(time.Since(start).Seconds())
				var res dto.WeatherResponse
				if err := json.Unmarshal([]byte(data), &res); err == nil {
					return &res, nil
				}
			} else if !errors.Is(err, redis.Nil) {
				r.metrics.ObserveLockWaitDuration(time.Since(start).Seconds())
				return nil, err
			}

			if time.Since(start) > r.lockMaxWait {
				r.metrics.ObserveLockWaitDuration(time.Since(start).Seconds())
				return nil, errors.New("timeout waiting for cache fill")
			}
		}
	}
}

func (r *Repository) ReleaseLock(ctx context.Context, city string) error {
	lockKey := r.getLockKey(city)
	return r.client.Del(ctx, lockKey).Err()
}

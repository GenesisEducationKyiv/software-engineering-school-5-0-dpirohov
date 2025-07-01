package weather

import (
	"context"
	"fmt"

	"sync"

	"weatherApi/internal/dto"
)

type MockCacheRepo struct {
	mu       sync.Mutex
	data     map[string]*dto.WeatherResponse
	locks    map[string]bool
	lockCond map[string]*sync.Cond
}

func NewMockCacheRepo() *MockCacheRepo {
	return &MockCacheRepo{
		data:     make(map[string]*dto.WeatherResponse),
		locks:    make(map[string]bool),
		lockCond: make(map[string]*sync.Cond),
	}
}

func (m *MockCacheRepo) AcquireLock(ctx context.Context, city string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.locks[city] {
		return false, nil
	}
	m.locks[city] = true
	if _, exists := m.lockCond[city]; !exists {
		m.lockCond[city] = sync.NewCond(&m.mu)
	}
	return true, nil
}

func (m *MockCacheRepo) ReleaseLock(ctx context.Context, city string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.locks[city] = false
	if cond, exists := m.lockCond[city]; exists {
		cond.Broadcast()
	}
	return nil
}

func (m *MockCacheRepo) WaitForUnlock(ctx context.Context, city string) (*dto.WeatherResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cond, exists := m.lockCond[city]
	if !exists {
		cond = sync.NewCond(&m.mu)
		m.lockCond[city] = cond
	}

	for m.locks[city] {
		done := make(chan struct{})
		go func() {
			cond.Wait()
			close(done)
		}()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-done:
		}
	}

	val, ok := m.data[city]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return val, nil
}

func (m *MockCacheRepo) Set(ctx context.Context, city string, data *dto.WeatherResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[city] = data
	return nil
}

func (m *MockCacheRepo) Get(ctx context.Context, city string) (*dto.WeatherResponse, error) {
	return m.WaitForUnlock(ctx, city)
}

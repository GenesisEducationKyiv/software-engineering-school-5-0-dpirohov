package subscription

import "context"

type MockSubscriptionRepository struct {
	FindOneOrNoneFn func(query any, args ...any) (*SubscriptionModel, error)
	CreateOneFn     func(entity *SubscriptionModel) error
	UpdateFn        func(entity *SubscriptionModel) error
	DeleteFn        func(entity *SubscriptionModel) error
}

func (m *MockSubscriptionRepository) FindOneOrNone(_ context.Context, q any, args ...any) (*SubscriptionModel, error) {
	return m.FindOneOrNoneFn(q, args...)
}

func (m *MockSubscriptionRepository) CreateOne(_ context.Context, e *SubscriptionModel) error {
	return m.CreateOneFn(e)
}

func (m *MockSubscriptionRepository) Update(_ context.Context, e *SubscriptionModel) error {
	return m.UpdateFn(e)
}

func (m *MockSubscriptionRepository) Delete(_ context.Context, e *SubscriptionModel) error {
	return m.DeleteFn(e)
}

func (m *MockSubscriptionRepository) FindOneOrCreate(
	context context.Context,
	conditions map[string]any,
	model *SubscriptionModel,
) (*SubscriptionModel, error) {
	return &SubscriptionModel{}, nil
}

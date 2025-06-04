package subscription

type MockSubscriptionRepository struct {
	FindOneOrNoneFn func(query any, args ...any) (*SubscriptionModel, error)
	CreateOneFn     func(entity *SubscriptionModel) error
	UpdateFn        func(entity *SubscriptionModel) error
	DeleteFn        func(entity *SubscriptionModel) error
}

func (m *MockSubscriptionRepository) FindOneOrNone(q any, args ...any) (*SubscriptionModel, error) {
	return m.FindOneOrNoneFn(q, args...)
}

func (m *MockSubscriptionRepository) CreateOne(e *SubscriptionModel) error {
	return m.CreateOneFn(e)
}

func (m *MockSubscriptionRepository) Update(e *SubscriptionModel) error {
	return m.UpdateFn(e)
}

func (m *MockSubscriptionRepository) Delete(e *SubscriptionModel) error {
	return m.DeleteFn(e)
}

func (m *MockSubscriptionRepository) FindOneOrCreate(
	map[string]any,
	*SubscriptionModel,
) (*SubscriptionModel, error) {
	return nil, nil
}

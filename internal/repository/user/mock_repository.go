package user

import "context"

type MockUserRepository struct {
	FindOneOrCreateFn func(conditions map[string]any, entity *UserModel) (*UserModel, error)
}

func (m *MockUserRepository) FindOneOrNone(ctx context.Context, q any, args ...any) (*UserModel, error) {
	return &UserModel{}, nil
}

func (m *MockUserRepository) CreateOne(ctx context.Context, e *UserModel) error {
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, e *UserModel) error {
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, e *UserModel) error {
	return nil
}

func (m *MockUserRepository) FindOneOrCreate(ctx context.Context, c map[string]any, e *UserModel) (*UserModel, error) {
	return m.FindOneOrCreateFn(c, e)
}

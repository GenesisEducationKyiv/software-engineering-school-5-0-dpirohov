package user

type MockUserRepository struct {
	FindOneOrCreateFn func(conditions map[string]any, entity *UserModel) (*UserModel, error)
}

func (m *MockUserRepository) FindOneOrNone(q any, args ...any) (*UserModel, error) {
	return &UserModel{}, nil
}

func (m *MockUserRepository) CreateOne(e *UserModel) error {
	return nil
}

func (m *MockUserRepository) Update(e *UserModel) error {
	return nil
}

func (m *MockUserRepository) Delete(e *UserModel) error {
	return nil
}

func (m *MockUserRepository) FindOneOrCreate(c map[string]any, e *UserModel) (*UserModel, error) {
	return m.FindOneOrCreateFn(c, e)
}

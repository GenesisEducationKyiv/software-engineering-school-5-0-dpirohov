package user

import (
	"weatherApi/internal/repository/base"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	FindOneOrNone(query any, args ...any) (*UserModel, error)
	CreateOne(entity *UserModel) error
	Update(entity *UserModel) error
	Delete(entity *UserModel) error
	FindOneOrCreate(conditions map[string]any, entity *UserModel) (*UserModel, error)
}

type UserRepository struct {
	*base.BaseRepository[UserModel]
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{
		BaseRepository: base.NewRepository[UserModel](db),
	}
}
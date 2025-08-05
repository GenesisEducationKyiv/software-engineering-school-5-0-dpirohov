package user

import (
	"context"
	"weatherApi/internal/repository/base"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	FindOneOrNone(ctx context.Context, query any, args ...any) (*UserModel, error)
	CreateOne(ctx context.Context, entity *UserModel) error
	Update(ctx context.Context, entity *UserModel) error
	Delete(ctx context.Context, entity *UserModel) error
	FindOneOrCreate(ctx context.Context, conditions map[string]any, entity *UserModel) (*UserModel, error)
}

type UserRepository struct {
	*base.BaseRepository[UserModel]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: base.NewRepository[UserModel](db),
	}
}

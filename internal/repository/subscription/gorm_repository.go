package subscription

import (
	"context"
	"weatherApi/internal/repository/base"

	"gorm.io/gorm"
)

type SubscriptionRepositoryInterface interface {
	FindOneOrNone(ctx context.Context, query any, args ...any) (*SubscriptionModel, error)
	FindOneOrCreate(
		ctx context.Context,
		conditions map[string]any,
		entity *SubscriptionModel,
	) (*SubscriptionModel, error)
	CreateOne(ctx context.Context, entity *SubscriptionModel) error
	Update(ctx context.Context, entity *SubscriptionModel) error
	Delete(ctx context.Context, entity *SubscriptionModel) error
}

type SubscriptionRepository struct {
	*base.BaseRepository[SubscriptionModel]
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepositoryInterface {
	return &SubscriptionRepository{
		BaseRepository: base.NewRepository[SubscriptionModel](db),
	}
}

package subscription

import (
	"weatherApi/internal/repository/base"

	"gorm.io/gorm"
)

type SubscriptionRepositoryInterface interface {
	FindOneOrNone(query any, args ...any) (*SubscriptionModel, error)
	FindOneOrCreate(conditions map[string]any, entity *SubscriptionModel) (*SubscriptionModel, error)
	CreateOne(entity *SubscriptionModel) error
	Update(entity *SubscriptionModel) error
	Delete(entity *SubscriptionModel) error
}

type SubscriptionRepository struct {
	*base.BaseRepository[SubscriptionModel]
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepositoryInterface {
	return &SubscriptionRepository{
		BaseRepository: base.NewRepository[SubscriptionModel](db),
	}
}

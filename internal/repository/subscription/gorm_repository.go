package subscription

import (
	"context"
	"weatherApi/internal/common/constants"
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
	FindAllSubscriptionsByFrequency(ctx context.Context, frequency constants.Frequency) ([]SubscriptionModel, error)
}

type SubscriptionRepository struct {
	*base.BaseRepository[SubscriptionModel]
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		BaseRepository: base.NewRepository[SubscriptionModel](db),
	}
}

func (r *SubscriptionRepository) FindAllSubscriptionsByFrequency(
	ctx context.Context,
	frequency constants.Frequency,
) ([]SubscriptionModel, error) {
	var entities []SubscriptionModel

	result := r.DB.WithContext(ctx).
		Preload("User").
		Where("frequency = ? AND is_confirmed = ?", frequency, true).
		Find(&entities)

	return entities, result.Error
}

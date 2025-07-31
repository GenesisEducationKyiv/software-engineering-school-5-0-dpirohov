package subscription

import (
	"context"
	"weatherApi/internal/common/constants"
	"weatherApi/internal/repository/base"

	"gorm.io/gorm"
)

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

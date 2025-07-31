package base

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

var ErrNotFound = errors.New("record not found")

func NewRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, query any, args ...any) ([]T, error) {
	var entities []T

	db := r.DB.WithContext(ctx)
	if query != nil {
		db = db.Where(query, args...)
	}

	result := db.Find(&entities)
	return entities, result.Error
}

func (r *BaseRepository[T]) FindOneOrNone(ctx context.Context, query any, args ...any) (*T, error) {
	var entity T
	result := r.DB.WithContext(ctx).Where(query, args...).First(&entity)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (r *BaseRepository[T]) CreateOne(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) FindOneOrCreate(ctx context.Context, conditions map[string]any, entity *T) (*T, error) {
	err := r.DB.WithContext(ctx).Where(conditions).FirstOrCreate(entity).Error
	return entity, err
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Delete(entity).Error
}

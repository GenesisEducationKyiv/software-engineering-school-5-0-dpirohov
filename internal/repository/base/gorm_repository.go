package base

import (
	"errors"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

var ErrNotFound = errors.New("record not found")

func NewRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		DB: db,
	}
}

func (r *BaseRepository[T]) FindAll() ([]T, error) {
	var entities []T
	result := r.DB.Find(&entities)
	return entities, result.Error
}

func (r *BaseRepository[T]) FindOneOrNone(query any, args ...any) (*T, error) {
	var entity T
	result := r.DB.Where(query, args...).First(&entity)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &entity, nil
}

func (r *BaseRepository[T]) CreateOne(entity *T) error {
	result := r.DB.Create(entity)
	return result.Error
}

func (r *BaseRepository[T]) FindOneOrCreate(conditions map[string]any, entity *T) (*T, error) {
	err := r.DB.Where(conditions).FirstOrCreate(entity).Error
	return entity, err
}

func (r *BaseRepository[T]) Update(entity *T) error {
	return r.DB.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(entity *T) error {
	return r.DB.Delete(entity).Error
}

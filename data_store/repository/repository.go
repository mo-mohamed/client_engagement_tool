package repository

import (
	"gorm.io/gorm"
)

type repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *repository[T] {
	return &repository[T]{
		db: db,
	}
}

func (r *repository[T]) Add(entity *T) error {
	return r.db.Create(&entity).Error
}

func (r *repository[T]) GetById(id int) (*T, error) {
	var entity T
	err := r.db.Model(&entity).Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *repository[T]) Exists(id int) (*T, bool) {
	var entity T
	err := r.db.Model(&entity).Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, false
	}
	return &entity, true
}

func (r *repository[T]) Get(params *T) *T {
	var entity T
	r.db.Where(&params).FirstOrInit(&entity)
	return &entity
}

func (r *repository[T]) GetAll() (*[]T, error) {
	var entities []T
	err := r.db.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *repository[T]) Where(params *T) ([]T, error) {
	var entities []T
	err := r.db.Where(&params).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository[T]) Update(entity *T) error {
	return r.db.Save(&entity).Error
}

func (r repository[T]) UpdateAll(entities *[]T) error {
	return r.db.Save(&entities).Error
}

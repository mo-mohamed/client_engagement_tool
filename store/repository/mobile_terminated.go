package repository

import (
	"context"
	"customer_engagement/store/models"

	"gorm.io/gorm"
)

type mtRepository struct{ db *gorm.DB }

func NewMTRepository(db *gorm.DB) *mtRepository {
	return &mtRepository{db: db}
}

func (repo *mtRepository) CreateMT(ctx context.Context, mt *models.MT) (*models.MT, error) {
	err := repo.db.WithContext(ctx).Create(mt).Error
	return mt, err
}

func (repo *mtRepository) IsProcessed(ctx context.Context, mt *models.MT) bool {
	query := "SELECT processed FROM mobile_terminated WHERE group_id = ? AND profile_id = ? AND broadcast_id = ?"
	var processed bool
	repo.db.WithContext(ctx).Raw(query, mt.GroupId, mt.ProfileId, mt.BroadcastId).Scan(processed)
	return processed
}

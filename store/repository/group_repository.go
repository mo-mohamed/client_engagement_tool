package repository

import (
	"context"
	"customer_engagement/store/interfaces"
	"customer_engagement/store/models"
	"database/sql"

	"gorm.io/gorm"
)

type groupRepository struct{ db *gorm.DB }

func NewGroupRepo(db *gorm.DB) interfaces.IGroupRepository {
	return &groupRepository{db: db}
}

// Get Group by id.
func (repo *groupRepository) GetGroup(ctx context.Context, id int) (*models.GroupStore, error) {
	var dbGroup *models.GroupStore = &models.GroupStore{ID: id}
	err := repo.db.WithContext(ctx).First(dbGroup).Error
	return dbGroup, err
}

// Create a Group
func (repo *groupRepository) CreateGroup(ctx context.Context, group *models.GroupStore) (*models.GroupStore, error) {
	err := repo.db.WithContext(ctx).Create(group).Error
	return group, err
}

// Updates a Group
func (repo *groupRepository) UpdateGroup(ctx context.Context, group *models.GroupStore) (*models.GroupStore, error) {
	err := repo.db.WithContext(ctx).Model(group).Updates(models.GroupStore{
		Name:      group.Name,
		DeletedAt: group.DeletedAt}).Error

	return group, err
}

func (repo *groupRepository) Exists(ctx context.Context, groupID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM group where id = ?`
	err := repo.db.WithContext(ctx).Exec(query, groupID).Scan(&exists).Error
	if err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

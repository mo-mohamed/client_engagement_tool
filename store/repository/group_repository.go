package repository

import (
	"context"
	"customer_engagement/store/interfaces"
	"customer_engagement/store/models"
	"time"

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

// Checks existence of a group by id
func (repo *groupRepository) Exists(ctx context.Context, groupID int) (bool, error) {
	var exists bool = false
	err := repo.db.Model(&models.GroupStore{}).Select("COUNT(*) > 0").Where("id = ?", groupID).Find(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Counts how many profiles to process per group that recorded before a specific insertion date
func (repo *groupRepository) CountNumberOfProfilesToProcess(ctx context.Context, grId int, dateEnqueued time.Time) int {
	query := `
		SELECT COUNT(*) FROM group_profile gp
		INNER JOIN profile p on gp.profile_id = p.id
		WHERE gp.group_id = ? AND gp.created_at <= ?;
	`
	var count int
	repo.db.Raw(query, grId, dateEnqueued).Scan(&count)
	return count
}

package dbRepository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type GroupStore struct {
	ID        int        `gorm:"column:id; primaryKey"`
	Name      *string    `gorm:"column:name"`
	CreatedAt time.Time  `gorm:"column:created_at; autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at; autoCreateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	// Profiles  []ProfileStore `gorm:"many2many:group_profile"`
}

func (GroupStore) TableName() string { return "group" }

func NewGroupRepo(db *gorm.DB) *dbGroupRepository {
	return &dbGroupRepository{db: db}
}

type IGroupRepository interface {
	GetGroup(ctx context.Context, id int) (*GroupStore, error)
	CreateGroup(ctx context.Context, group *GroupStore) (*GroupStore, error)
	UpdateGroup(ctx context.Context, group *GroupStore) (*GroupStore, error)
}

type dbGroupRepository struct{ db *gorm.DB }

// Get Group by id.
func (repo *dbGroupRepository) GetGroup(ctx context.Context, id int) (*GroupStore, error) {
	var dbGroup *GroupStore = &GroupStore{ID: id}
	err := repo.db.WithContext(ctx).First(dbGroup).Error
	return dbGroup, err
}

// Create a Group
func (repo *dbGroupRepository) CreateGroup(ctx context.Context, group *GroupStore) (*GroupStore, error) {
	err := repo.db.WithContext(ctx).Create(group).Error
	return group, err
}

// Updates a Group
func (repo *dbGroupRepository) UpdateGroup(ctx context.Context, group *GroupStore) (*GroupStore, error) {
	err := repo.db.WithContext(ctx).Model(group).Updates(GroupStore{
		Name:      group.Name,
		DeletedAt: group.DeletedAt}).Error

	return group, err
}

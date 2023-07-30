package repository

import (
	"context"
	"time"

	"customer_engagement/store/interfaces"
	"customer_engagement/store/models"

	"gorm.io/gorm"
)

type profileRepo struct{ db *gorm.DB }

func NewProfileRepo(db *gorm.DB) interfaces.IProfileRepo {
	return &profileRepo{db: db}
}

// Persists a new profile.
func (repo *profileRepo) CreateProfile(ctx context.Context, profile *models.ProfileStore) (*models.ProfileStore, error) {
	err := repo.db.WithContext(ctx).Create(profile).Error
	return profile, err
}

// Finds a profile by id.
func (repo *profileRepo) GetProfile(ctx context.Context, id int) (*models.ProfileStore, error) {
	profile := &models.ProfileStore{ID: id}
	err := repo.db.WithContext(ctx).First(profile).Error
	return profile, err
}

// Updates a profile.
func (repo *profileRepo) UpdateProfile(ctx context.Context, profile *models.ProfileStore) (*models.ProfileStore, error) {
	err := repo.db.WithContext(ctx).Model(profile).Updates(models.ProfileStore{
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		MDN:       profile.MDN}).Error

	return profile, err
}

// Gets profiles paginated.
func (repo *profileRepo) GetProfilesPaginated(ctx context.Context, limit int, offset int) ([]models.ProfileStore, error) {
	var result []models.ProfileStore
	err := repo.db.WithContext(ctx).Order("created_at ASC").Limit(limit).Offset(offset).Find(&result).Error
	return result, err
}

func (repo *profileRepo) AddProfileToGroup(ctx context.Context, profileId, groupId int) error {
	query := `INSERT INTO group_profile (profile_id, group_id) VALUES (?, ?)`
	return repo.db.WithContext(ctx).Exec(query, profileId, groupId).Error
}

func (repo *profileRepo) GetGroupProfilesPaginated(ctx context.Context, grId, limit, offset int, enqueueDate time.Time) ([]models.ProfileStore, error) {
	var result []models.ProfileStore
	query := `
	SELECT p.* FROM profile p INNER JOIN group_profile gp ON p.id = gp.profile_id
	WHERE gp.group_id = ?
	AND gp.created_at <= ?
	ORDER BY p.created_at ASC
	LIMIT ?, ?
	`
	err := repo.db.WithContext(ctx).Raw(query, grId, enqueueDate, limit, offset).Scan(&result).Error
	return result, err
}

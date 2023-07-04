package dbRepository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type IProfileRepo interface {
	CreateProfile(ctx context.Context, profile *ProfileStore) (*ProfileStore, error)
	GetProfile(ctx context.Context, id int) (*ProfileStore, error)
	UpdateProfile(ctx context.Context, profile *ProfileStore) (*ProfileStore, error)
	GetProfilesPaginated(ctx context.Context, limit int, offset int) ([]ProfileStore, error)
}

type dbProfileRepo struct{ db *gorm.DB }

type ProfileStore struct {
	ID        int        `gorm:"column:id; primaryKey"`
	FirstName *string    `gorm:"column:first_name"`
	LastName  *string    `gorm:"column:last_name"`
	MDN       string     `gorm:"column:mdn"`
	CreatedAt time.Time  `gorm:"column:created_at; autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at; autoCreateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	// Groups    []GroupStore `gorm:"many2many:group_profile"`
}

func (ProfileStore) TableName() string { return "profile" }

// Returns a new profile repository instance.
func NewProfileRepo(db *gorm.DB) *dbProfileRepo {
	return &dbProfileRepo{db: db}
}

// Persists a new profile.
func (repo *dbProfileRepo) CreateProfile(ctx context.Context, profile *ProfileStore) (*ProfileStore, error) {
	err := repo.db.WithContext(ctx).Create(profile).Error
	return profile, err
}

// Finds a profile by id.
func (repo *dbProfileRepo) GetProfile(ctx context.Context, id int) (*ProfileStore, error) {
	profile := &ProfileStore{ID: id}
	err := repo.db.WithContext(ctx).First(profile).Error
	return profile, err
}

// Updates a profile.
func (repo *dbProfileRepo) UpdateProfile(ctx context.Context, profile *ProfileStore) (*ProfileStore, error) {
	err := repo.db.Model(profile).Updates(ProfileStore{
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		MDN:       profile.MDN}).Error

	return profile, err
}

// Gets profiles paginated.
func (repo *dbProfileRepo) GetProfilesPaginated(ctx context.Context, limit int, offset int) ([]ProfileStore, error) {
	var result []ProfileStore
	err := repo.db.Order("created_at ASC").Limit(limit).Offset(offset).Find(&result).Error
	return result, err
}

// Gets profiles paginated for a Group.
func (repo *dbProfileRepo) GetGroupProfilesPaginated(ctx context.Context, groupId, limit, offset int) ([]ProfileStore, error) {
	var result []ProfileStore
	query := `
		SELECT * FROM profile p
		INNER JOIN group_profile gp
		ON p.id = gp.profile_id
		WHERE gp.group_id = ?
		ORDER BY p.created_at ASC
		LIMIT ?
		OFFSET ?
		`
	err := repo.db.Raw(query, groupId, limit, offset).Scan(result).Error
	return result, err
}

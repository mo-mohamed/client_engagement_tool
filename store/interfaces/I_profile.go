package interfaces

import (
	"context"
	"customer_engagement/store/models"
	"time"
)

type IProfileRepo interface {
	CreateProfile(ctx context.Context, profile *models.ProfileStore) (*models.ProfileStore, error)
	GetProfile(ctx context.Context, id int) (*models.ProfileStore, error)
	UpdateProfile(ctx context.Context, profile *models.ProfileStore) (*models.ProfileStore, error)
	GetProfilesPaginated(ctx context.Context, limit int, offset int) ([]models.ProfileStore, error)
	GetGroupProfilesPaginated(ctx context.Context, grId, limit, offset int, enqueueDate time.Time) ([]models.ProfileStore, error)
	AddProfileToGroup(ctx context.Context, profileId, groupId int) error
}

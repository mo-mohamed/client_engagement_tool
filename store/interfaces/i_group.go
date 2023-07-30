package interfaces

import (
	"context"
	"customer_engagement/store/models"
	"time"
)

type IGroupRepository interface {
	GetGroup(ctx context.Context, id int) (*models.GroupStore, error)
	CreateGroup(ctx context.Context, group *models.GroupStore) (*models.GroupStore, error)
	UpdateGroup(ctx context.Context, group *models.GroupStore) (*models.GroupStore, error)
	Exists(ctx context.Context, groupID int) (bool, error)
	CountNumberOfProfilesToProcess(ctx context.Context, grId int, dateEnqueued time.Time) int
}

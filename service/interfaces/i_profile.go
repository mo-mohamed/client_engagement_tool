package interfaces

import (
	"context"
	domains "customer_engagement/service/models"
)

type IProfileService interface {
	Create(ctx context.Context, profile *domains.Profile) (*domains.Profile, error)
	Get(ctx context.Context, id int) (*domains.Profile, error)
	Update(ctx context.Context, profile *domains.Profile) (*domains.Profile, error)
	GetPaginated(ctx context.Context, limit int, offset int) ([]*domains.Profile, error)
	AttachToGroup(ctx context.Context, profileId, groupId int) error
}

package interfaces

import (
	"context"
	domains "customer_engagement/service/models"
)

type IGroupService interface {
	Get(ctx context.Context, id int) (*domains.Group, error)
	Create(ctx context.Context, group *domains.Group) (*domains.Group, error)
	Update(ctx context.Context, group *domains.Group) (*domains.Group, error)
	Exists(ctx context.Context, groupID int) (bool, error)
}

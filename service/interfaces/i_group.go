package interfaces

import (
	"context"
	domains "customer_engagement/service/models"
)

type IGroupService interface {
	GetGroup(ctx context.Context, id int) (*domains.Group, error)
	CreateGroup(ctx context.Context, group *domains.Group) (*domains.Group, error)
	UpdateGroup(ctx context.Context, group *domains.Group) (*domains.Group, error)
	Exists(ctx context.Context, groupID int) (bool, error)
}

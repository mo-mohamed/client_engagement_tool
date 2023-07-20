package service

import (
	"context"
	storeLayer "customer_engagement/store"
	store "customer_engagement/store/models"
)

type IGroupService interface {
	GetGroup(ctx context.Context, id int) (*store.GroupStore, error)
	CreateGroup(ctx context.Context, group *store.GroupStore) (*store.GroupStore, error)
	UpdateGroup(ctx context.Context, group *store.GroupStore) (*store.GroupStore, error)
	Exists(ctx context.Context, groupID int) (bool, error)
}
type GroupService struct {
	store *storeLayer.Store
}

func NewGroupService(store *storeLayer.Store) IGroupService {
	return &GroupService{store: store}
}

// Get Group by id.
func (g *GroupService) GetGroup(ctx context.Context, id int) (*store.GroupStore, error) {
	return g.store.Group.GetGroup(ctx, id)
}

// Create a Group
func (g *GroupService) CreateGroup(ctx context.Context, group *store.GroupStore) (*store.GroupStore, error) {
	return g.store.Group.CreateGroup(ctx, group)
}

// Updates a Group
func (g *GroupService) UpdateGroup(ctx context.Context, group *store.GroupStore) (*store.GroupStore, error) {
	return g.store.Group.UpdateGroup(ctx, group)
}

// Updates a Group
func (g *GroupService) Exists(ctx context.Context, groupID int) (bool, error) {
	return g.store.Group.Exists(ctx, groupID)
}

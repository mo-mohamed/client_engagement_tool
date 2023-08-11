package service

import (
	"context"
	interfaces "customer_engagement/service/interfaces"
	domains "customer_engagement/service/models"
	storeLayer "customer_engagement/store"
	storeModels "customer_engagement/store/models"
)

type GroupService struct {
	store *storeLayer.Store
}

func NewGroupService(store *storeLayer.Store) interfaces.IGroupService {
	return &GroupService{store: store}
}

// Get Group by id.
func (gs *GroupService) Get(ctx context.Context, id int) (*domains.Group, error) {
	g, err := gs.store.Group.GetGroup(ctx, id)
	if err != nil {
		return nil, err
	}
	return gs.toDomain(g), err
}

// Create a Group
func (gs *GroupService) Create(ctx context.Context, group *domains.Group) (*domains.Group, error) {
	g, err := gs.store.Group.CreateGroup(ctx, gs.toDatabaseEntity(group))
	if err != nil {
		return nil, err
	}
	return gs.toDomain(g), err
}

// Updates a Group
func (gs *GroupService) Update(ctx context.Context, group *domains.Group) (*domains.Group, error) {
	g, err := gs.store.Group.UpdateGroup(ctx, gs.toDatabaseEntity(group))
	if err != nil {
		return nil, err
	}

	return gs.toDomain(g), err
}

// Checks existence of a group by id
func (gs *GroupService) Exists(ctx context.Context, groupID int) (bool, error) {
	return gs.store.Group.Exists(ctx, groupID)
}

func (*GroupService) toDomain(g *storeModels.GroupStore) *domains.Group {
	return &domains.Group{
		ID:        g.ID,
		Name:      g.Name,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
		DeletedAt: g.DeletedAt,
	}
}

func (*GroupService) toDatabaseEntity(g *domains.Group) *storeModels.GroupStore {
	return &storeModels.GroupStore{
		ID:        g.ID,
		Name:      g.Name,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
		DeletedAt: g.DeletedAt,
	}
}

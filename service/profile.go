package service

import (
	"context"
	storeLayer "customer_engagement/store"
	store "customer_engagement/store/models"
)

type ProfileService struct {
	store *storeLayer.Store
}

func NewProfileService(store *storeLayer.Store) IProfileService {
	return &ProfileService{store: store}
}

type IProfileService interface {
	Create(ctx context.Context, profile *store.ProfileStore) (*store.ProfileStore, error)
	Get(ctx context.Context, id int) (*store.ProfileStore, error)
	Update(ctx context.Context, profile *store.ProfileStore) (*store.ProfileStore, error)
	GetPaginated(ctx context.Context, limit int, offset int) ([]store.ProfileStore, error)
	AttachToGroup(ctx context.Context, profileId, groupId int) error
}

func (ps *ProfileService) Create(ctx context.Context, profile *store.ProfileStore) (*store.ProfileStore, error) {
	return ps.store.Profile.CreateProfile(ctx, profile)
}

func (ps *ProfileService) Get(ctx context.Context, id int) (*store.ProfileStore, error) {
	return ps.store.Profile.GetProfile(ctx, id)
}

func (ps *ProfileService) Update(ctx context.Context, profile *store.ProfileStore) (*store.ProfileStore, error) {
	return ps.store.Profile.UpdateProfile(ctx, profile)
}

func (ps *ProfileService) GetPaginated(ctx context.Context, limit int, offset int) ([]store.ProfileStore, error) {
	return ps.store.Profile.GetProfilesPaginated(ctx, limit, offset)
}

func (ps *ProfileService) AttachToGroup(ctx context.Context, profileId, groupId int) error {
	return ps.store.Profile.AddProfileToGroup(ctx, profileId, groupId)
}

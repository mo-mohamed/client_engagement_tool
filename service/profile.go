package service

import (
	"context"
	interfaces "customer_engagement/service/interfaces"
	domains "customer_engagement/service/models"
	storeLayer "customer_engagement/store"
	store "customer_engagement/store/models"
)

type ProfileService struct {
	store *storeLayer.Store
}

func NewProfileService(store *storeLayer.Store) interfaces.IProfileService {
	return &ProfileService{store: store}
}

func (ps *ProfileService) Create(ctx context.Context, profile *domains.Profile) (*domains.Profile, error) {
	p, err := ps.store.Profile.CreateProfile(ctx, ps.toDatabaseEntity(profile))
	if err != nil {
		return nil, err
	}
	return ps.toDomain(p), err
}

func (ps *ProfileService) Get(ctx context.Context, id int) (*domains.Profile, error) {
	p, err := ps.store.Profile.GetProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	return ps.toDomain(p), nil
}

func (ps *ProfileService) Update(ctx context.Context, profile *domains.Profile) (*domains.Profile, error) {
	p, err := ps.store.Profile.UpdateProfile(ctx, ps.toDatabaseEntity(profile))
	if err != nil {
		return nil, err
	}
	return ps.toDomain(p), err
}

func (ps *ProfileService) GetPaginated(ctx context.Context, limit int, offset int) ([]*domains.Profile, error) {
	profiles, err := ps.store.Profile.GetProfilesPaginated(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	domainProfiles := make([]*domains.Profile, 0)
	for _, p := range profiles {
		domainProfiles = append(domainProfiles, ps.toDomain(&p))
	}
	return domainProfiles, nil
}

func (ps *ProfileService) AttachToGroup(ctx context.Context, profileId, groupId int) error {
	return ps.store.Profile.AddProfileToGroup(ctx, profileId, groupId)
}

func (ps *ProfileService) toDomain(profile *store.ProfileStore) *domains.Profile {
	return &domains.Profile{
		ID:        profile.ID,
		FirstName: profile.FirstName,
		MDN:       *profile.LastName,
		LastName:  profile.LastName,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
		DeletedAt: profile.DeletedAt,
	}
}

func (ps *ProfileService) toDatabaseEntity(profile *domains.Profile) *store.ProfileStore {
	return &store.ProfileStore{
		ID:        profile.ID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		MDN:       profile.MDN,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
		DeletedAt: profile.DeletedAt,
	}
}

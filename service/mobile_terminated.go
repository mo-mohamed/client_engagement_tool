package service

import (
	"context"
	domains "customer_engagement/service/models"
	storeLayer "customer_engagement/store"
	storeModels "customer_engagement/store/models"
)

type MTService struct {
	store storeLayer.Store
}

func NewMtService(store storeLayer.Store) *MTService {
	return &MTService{store: store}
}

func (mts *MTService) create(ctx context.Context, mt domains.MT) (*domains.MT, error) {
	dbMt, error := mts.store.MT.CreateMT(ctx, mts.toDatabaseEntity(&mt))
	return mts.toDomain(dbMt), error
}

func (mts *MTService) IsProcessed(ctx context.Context, mt domains.MT) bool {
	return mts.store.MT.IsProcessed(ctx, mts.toDatabaseEntity(&mt))
}

func (mts *MTService) toDomain(mt *storeModels.MT) *domains.MT {
	return &domains.MT{
		ID:          mt.ID,
		ProfileId:   mt.ProfileId,
		GroupId:     mt.GroupId,
		BroadcastId: mt.BroadcastId,
		Processed:   mt.Processed,
		CreatedAt:   mt.CreatedAt,
		UpdatedAt:   mt.UpdatedAt,
		DeletedAt:   mt.DeletedAt,
	}
}

func (mts *MTService) toDatabaseEntity(mt *domains.MT) *storeModels.MT {
	return &storeModels.MT{
		ID:          mt.ID,
		ProfileId:   mt.ProfileId,
		GroupId:     mt.GroupId,
		BroadcastId: mt.BroadcastId,
		Processed:   mt.Processed,
		CreatedAt:   mt.CreatedAt,
		UpdatedAt:   mt.UpdatedAt,
		DeletedAt:   mt.DeletedAt,
	}
}

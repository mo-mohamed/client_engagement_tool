package viewModels

import (
	"customer_engagement/data_store/models"
	db_models "customer_engagement/data_store/models"
	"time"
)

type GroupProfile struct {
	ID        int        `json:"id"`
	ProfileId int        `json:"profile_id"`
	GroupId   int        `json:"group_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (g GroupProfile) ToDatabaseEntity() db_models.GroupProfile {
	return db_models.GroupProfile{
		ProfileId: g.ProfileId,
		GroupId:   g.GroupId,
	}
}

func (g GroupProfile) FromDatabaseEntity(dbgp models.GroupProfile) GroupProfile {
	return GroupProfile{
		ProfileId: dbgp.ProfileId,
		GroupId:   dbgp.GroupId,
		CreatedAt: dbgp.CreatedAt,
		UpdatedAt: dbgp.UpdatedAt,
		DeletedAt: dbgp.DeletedAt,
		ID:        dbgp.ID,
	}
}

package viewModels

import (
	db_models "customer_engagement/data_store/models"
	"time"
)

type Group struct {
	ID        int        `json:"id"`
	Name      *string    `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Active    bool       `json:"active"`
}

func (g Group) ToDTO() db_models.Group {
	return db_models.Group{
		Name: g.Name,
	}
}

func (g Group) FromDTO(dbGroup db_models.Group) Group {
	return_group := Group{
		Name:      dbGroup.Name,
		CreatedAt: dbGroup.CreatedAt,
		UpdatedAt: dbGroup.UpdatedAt,
		DeletedAt: dbGroup.DeletedAt,
		ID:        dbGroup.ID,
	}

	if dbGroup.DeletedAt == nil {
		return_group.Active = true
	} else {
		return_group.Active = false
	}

	return return_group
}

package viewModels

import (
	db_models "customer_engagement/data_store/models"
	"time"
)

type Profile struct {
	ID        int        `json:"id"`
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	MDN       *string    `json:"mdn"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Active    bool       `json:"active"`
	GroupID   *int       `json:"group_id"`
}

func (p Profile) ToDTO() db_models.Profile {
	return db_models.Profile{
		FirstName: p.FirstName,
		LastName:  p.LastName,
		MDN:       *p.MDN,
		GroupID:   p.GroupID,
	}
}

func (Profile) FromDTO(dbProfile db_models.Profile) Profile {
	returnProfile := Profile{
		FirstName: dbProfile.FirstName,
		LastName:  dbProfile.LastName,
		MDN:       &dbProfile.MDN,
		CreatedAt: dbProfile.CreatedAt,
		UpdatedAt: dbProfile.UpdatedAt,
		DeletedAt: dbProfile.DeletedAt,
		ID:        dbProfile.ID,
		GroupID:   dbProfile.GroupID,
	}

	if dbProfile.DeletedAt == nil {
		returnProfile.Active = true
	} else {
		returnProfile.Active = false
	}

	return returnProfile
}

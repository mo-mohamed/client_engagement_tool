/*
Holds profiles related to/from database objects convertions and API's input validations
*/

package viewModels

import (
	serviceModels "customer_engagement/service/models"
	"time"

	"gopkg.in/go-playground/validator.v9"
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
}

func (p Profile) ToDomain() serviceModels.Profile {
	return serviceModels.Profile{
		FirstName: p.FirstName,
		LastName:  p.LastName,
		MDN:       *p.MDN,
	}
}

func (Profile) FromService(p serviceModels.Profile) Profile {
	returnProfile := Profile{
		FirstName: p.FirstName,
		LastName:  p.LastName,
		MDN:       &p.MDN,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
		ID:        p.ID,
	}

	if p.DeletedAt == nil {
		returnProfile.Active = true
	} else {
		returnProfile.Active = false
	}

	return returnProfile
}

func (p *Profile) Validate() (bool, []*ValidationError) {
	validation := validator.New()
	err := validation.Struct(p)
	if err != nil {
		return false, convertErrors(err.(validator.ValidationErrors))
	}
	return true, nil
}

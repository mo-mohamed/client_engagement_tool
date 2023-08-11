/*
Holds groups related to/from database objects convertions and API's input validations
*/
package viewModels

import (
	"time"

	serviceModels "customer_engagement/service/models"

	"gopkg.in/go-playground/validator.v9"
)

type Group struct {
	ID        int        `json:"id"`
	Name      *string    `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Active    bool       `json:"active"`
}

func (g Group) ToDomain() serviceModels.Group {
	return serviceModels.Group{
		Name: g.Name,
	}
}

func (g Group) FromService(group serviceModels.Group) Group {
	return_group := Group{
		Name:      group.Name,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
		DeletedAt: group.DeletedAt,
		ID:        group.ID,
	}

	if group.DeletedAt == nil {
		return_group.Active = true
	} else {
		return_group.Active = false
	}

	return return_group
}

func (g *Group) Validate() (bool, []*ValidationError) {
	validation := validator.New()
	err := validation.Struct(g)
	if err != nil {
		return false, convertErrors(err.(validator.ValidationErrors))
	}
	return true, nil
}

type BroadcastRequest struct {
	GroupId     int    `json:"group_id" validate:"required"`
	MessageBody string `json:"message_body" validate:"required"`
}

func (bcr *BroadcastRequest) Validate() (bool, []*ValidationError) {
	validation := validator.New()
	err := validation.Struct(bcr)

	if err != nil {
		return false, convertErrors(err.(validator.ValidationErrors))
	}
	return true, nil
}

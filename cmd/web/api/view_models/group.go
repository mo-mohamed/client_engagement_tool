/*
Holds groups related to/from database objects convertions and API's input validations
*/
package viewModels

import (
	"time"

	db_models "customer_engagement/data_store/models"

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

func (g Group) ToDatabaseEntity() db_models.Group {
	return db_models.Group{
		Name: g.Name,
	}
}

func (g Group) FromDatabaseEntity(dbGroup db_models.Group) Group {
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

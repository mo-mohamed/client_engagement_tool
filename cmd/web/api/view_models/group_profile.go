package viewModels

import (
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

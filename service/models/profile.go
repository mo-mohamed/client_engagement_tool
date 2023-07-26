package models

import "time"

type Profile struct {
	ID        int
	FirstName *string
	LastName  *string
	MDN       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

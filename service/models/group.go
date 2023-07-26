package models

import "time"

type Group struct {
	ID        int
	Name      *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

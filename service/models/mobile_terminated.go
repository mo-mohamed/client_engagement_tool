package models

import "time"

type MT struct {
	ID          int
	ProfileId   int
	GroupId     int
	BroadcastId string
	Processed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

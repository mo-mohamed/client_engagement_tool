package models

import "time"

type MT struct {
	ID          int        `gorm:"column:id; primaryKey"`
	ProfileId   int        `gorm:"column:profile_id"`
	GroupId     int        `gorm:"column:group_id"`
	BroadcastId string     `gorm:"column:broadcast_id"`
	Processed   bool       `gorm:"column:processed"`
	CreatedAt   time.Time  `gorm:"column:created_at; autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at; autoCreateTime"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (MT) TableName() string { return "mobile_terminated" }

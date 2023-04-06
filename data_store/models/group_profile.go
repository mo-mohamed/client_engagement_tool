package models

import "time"

type GroupProfile struct {
	ID        int        `gorm:"column:id; primaryKey"`
	ProfileId int        `gorm:"column:profile_id"`
	GroupUd   int        `gorm:"column:group_id"`
	CreatedAt time.Time  `gorm:"column:created_at; autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at; autoCreateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

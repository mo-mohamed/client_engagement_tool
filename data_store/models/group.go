package models

import "time"

type Group struct {
	ID        int        `gorm:"column:id; primaryKey"`
	Name      *string    `gorm:"column:name"`
	CreatedAt time.Time  `gorm:"column:created_at; autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at; autoCreateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	Profiles  []Profile
}

func (Group) TableName() string {
	return "group"
}

package models

import "time"

type GroupStore struct {
	ID        int            `gorm:"column:id; primaryKey"`
	Name      *string        `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at; autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at; autoCreateTime"`
	DeletedAt *time.Time     `gorm:"column:deleted_at"`
	Profiles  []ProfileStore `gorm:"many2many:group_profile"`
}

func (GroupStore) TableName() string { return "group" }

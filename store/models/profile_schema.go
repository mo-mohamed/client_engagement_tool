package models

import "time"

type ProfileStore struct {
	ID        int          `gorm:"column:id; primaryKey"`
	FirstName *string      `gorm:"column:first_name"`
	LastName  *string      `gorm:"column:last_name"`
	MDN       string       `gorm:"column:mdn"`
	CreatedAt time.Time    `gorm:"column:created_at; autoCreateTime"`
	UpdatedAt time.Time    `gorm:"column:updated_at; autoCreateTime"`
	DeletedAt *time.Time   `gorm:"column:deleted_at"`
	Groups    []GroupStore `gorm:"many2many:group_profile"`
}

func (ProfileStore) TableName() string { return "profile" }

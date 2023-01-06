package models

import "time"

type Profile struct {
	ID        int        `gorm:"column:id; primaryKey"`
	FirstName *string    `gorm:"column:first_name"`
	LastName  *string    `gorm:"column:last_name"`
	MDN       string     `gorm:"column:mdn"`
	CreatedAt time.Time  `gorm:"column:created_at; autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at; autoCreateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	GroupID   *int       `gorm:"column:group_id"`
	Group     Group
}

func (Profile) TableName() string {
	return "profile"
}

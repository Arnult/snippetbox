package models

import (
	"gorm.io/gorm"
	"time"
)

type Snippets struct {
	ID        uint      `gorm:"primary_key;auto_increment"`
	CreatedAt time.Time `gorm:"index:idx_snippets_created"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index:idx_snippets_deleted"`
	Title     string         `gorm:"type:varchar(100);not null"`
	Content   string         `gorm:"type:text"`
	Expires   time.Time
}

func (Snippets) TableName() string {
	return "snippets"
}

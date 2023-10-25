package models

import (
	"gorm.io/gorm"
	"time"
)

type Snippet struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index:idx_snippets_created"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index:idx_snippets_deleted"`
	Title     string         `gorm:"type:varchar(100);not null"`
	Content   string         `gorm:"type:text"`
	Expires   time.Time
}

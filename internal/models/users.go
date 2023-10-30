package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Name           string     `gorm:"type:varchar(255);not null"`
	Email          string     `gorm:"type:varchar(255);not null;unique"`
	HashedPassword string     `gorm:"type:char(60);not null"`
	Snippets       []Snippets `gorm:"many2many:users_snippets;"`
}

func (Users) TableName() string {
	return "users"
}

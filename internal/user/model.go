package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `json:"name,omitempty" gorm:"not null"`
	Familia      string         `json:"familia,omitempty" gorm:"not null"`
	Role         string         `json:"role" gorm:"not null"`
	Password     string         `json:"password" gorm:"-"`
	PasswordHash string         `json:"-" gorm:"not null"`
	Username     string         `json:"username,omitempty" gorm:"not null,unique"`
}

package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint           `gorm:"primarykey"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `json:"name,omitempty" gorm:"not null"`
	Familia      string         `json:"familia,omitempty" gorm:"not null"`
	Role         string         `json:"role" gorm:"not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	Username     string         `json:"username,omitempty" gorm:"not null,unique"`
}

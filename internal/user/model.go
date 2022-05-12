package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `json:"name,omitempty" gorm:"not null"`
	Familia      string `json:"familia,omitempty" gorm:"not null"`
	Role         string `json:"role" gorm:"not null"`
	Password     string `json:"password" gorm:"-"`
	PasswordHash string `json:"-" gorm:"not null"`
	Username     string `json:"username,omitempty" gorm:"not null,unique"`
}

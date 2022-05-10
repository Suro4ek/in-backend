package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `json:"name,omitempty"`
	Familia      string `json:"familia,omitempty"`
	Role         string `json:"role"`
	Password     string `json:"password" gorm:"-"`
	PasswordHash string `json:"-"`
	Username     string `json:"username,omitempty"`
}

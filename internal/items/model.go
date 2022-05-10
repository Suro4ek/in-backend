package items

import (
	"gorm.io/gorm"
	"in-backend/internal/user"
)

type Item struct {
	gorm.Model
	ProductName  string    `json:"productName,omitempty"`
	SerialNumber string    `json:"serialNumber,omitempty"`
	Name         string    `json:"name,omitempty"`
	OwnerID      *int      `json:"owner,omitempty'"`
	Owner        user.User `json:"-'" gorm:"OnDelete:SET NULL;"`
}

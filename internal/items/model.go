package items

import (
	"gorm.io/gorm"
	"in-backend/internal/user"
	"time"
)

type Item struct {
	ID           uint           `gorm:"primarykey"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	ProductName  string         `json:"productName,omitempty"`
	SerialNumber string         `json:"serialNumber,omitempty"`
	Name         string         `json:"name,omitempty"`
	OwnerID      *int           `json:"owner,omitempty'"`
	Owner        user.User      `json:"-'" gorm:"OnDelete:SET NULL;"`
}

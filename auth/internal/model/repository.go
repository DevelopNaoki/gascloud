package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Common struct {
	ID        UUID `gorm:"primaryKey;type:binary(16);<-:create"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// If you include the Common structure and create a new function BeforeCreate,
// this process is overridden, so the following process must be added to the newly created function.
func (c *Common) BeforeCreate(db *gorm.DB) (err error) {
	c.ID = UUID(uuid.Must(uuid.NewV7()))
	return
}

type Account struct {
	Common
	Name        string `gorm:"unique"`
	Passwd      string `gorm:"not null"`
	Description string
	IsActive    bool `gorm:"default:true"`
}

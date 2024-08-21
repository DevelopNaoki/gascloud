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

type Role struct {
	Common
	Name        string `gorm:"unique"`
	Description string
	IsActive    bool `gorm:"default:true"`
}

type RoleBind struct {
	Common
	Account UUID `gorm:"not null;type:binary(16)"`
	Role    UUID `gorm:"not null;type:binary(16)"`
}

type Permission struct {
	Common
	Service UUID   `gorm:"not null;type:binary(16)"`
	Action  string `gorm:"not null"`
}

type PermissionBind struct {
	Common
	Role    UUID `gorm:"not null;type:binary(16)"`
	Service UUID `gorm:"not null;type:binary(16)"`
}

type ServiceCatalog struct {
	Common
	Name        string `gorm:"unique"`
	Endpoint    string `gorm:"unique"`
	Description string
	IsActive    bool `gorm:"default:true"`
}

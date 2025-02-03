package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Common struct {
	ID        UUID `gorm:"primaryKey;type:varchar(36);<-:create"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// If you include the Common structure and create a new function BeforeCreate,
// this process is overridden, so the following process must be added to the newly created function.
func (c *Common) BeforeCreate(db *gorm.DB) (err error) {
	uid := uuid.Must(uuid.NewV7())
	c.ID = UUID(uid)
	return
}

type Account struct {
	Common
	Name        string `gorm:"not null;unique"`
	Passwd      string `gorm:"not null"`
	MailAddr    string
	Description string
	Role        UUID `gorm:"type:varchar(36)"`
	IsActive    bool `gorm:"default:true"`
}

type GroupBind struct {
	Common
	Account UUID `gorm:"type:varchar(36)"`
	Group   UUID `gorm:"type:varchar(36)"`
}

type Group struct {
	Common
	Name        string `gorm:"not null;unique"`
	Description string
	IsActive    bool `gorm:"default:true"`
}

type Session struct {
	Common
	Account   UUID   `gorm:"type:varchar(36)"`
	Token     string `gorm:"type:varchar(36);unique"`
	ExpiredAt time.Time
}

type Role struct {
	Common
	Name        string `gorm:"not null;unique"`
	Description string
}

type ServiceCatalog struct {
	Common
	Name        string `gorm:"unique"`
	Endpoint    string `gorm:"unique"`
	Description string
	IsActive    bool `gorm:"default:true"`
}

type ServiceToken struct {
	Common
	Service     UUID `gorm:"type:varchar(36)"`
	Description string
	IsActive    bool `gorm:"default:true"`
}

type Categoly struct {
	Common
	Name string `gorm:"unique"`
}

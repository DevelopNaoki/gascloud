package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Common struct {
	ID        UUID      `gorm:"primaryKey;type:varchar(36);<-:create" json:"id"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
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
	Name        string `gorm:"not null;unique" json:"account"`
	Passwd      string `gorm:"not null" json:"-"`
	MailAddr    string `json:"mail_address"`
	Description string `json:"description"`
	Role        UUID   `gorm:"type:varchar(36)" json:"role_id"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

type GroupBind struct {
	Common
	Account UUID `gorm:"type:varchar(36)" json:"account_id"`
	Group   UUID `gorm:"type:varchar(36)" json:"group_id"`
}

type Group struct {
	Common
	Name        string `gorm:"not null;unique" json:"group"`
	Description string `json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

type Session struct {
	Common
	Account   UUID      `gorm:"type:varchar(36)" json:"account_id"`
	Token     string    `gorm:"type:varchar(36);unique" json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type Role struct {
	Common
	Name        string `gorm:"not null;unique" json:"role"`
	Description string `json:"description"`
}

type ServiceCatalog struct {
	Common
	Name        string `gorm:"unique" json:"service_catalog"`
	Endpoint    string `gorm:"unique" json:"endpoint"`
	Description string `json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

type ServiceToken struct {
	Common
	Service     UUID   `gorm:"type:varchar(36)" json:"service_id"`
	Description string `json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

type Categoly struct {
	Common
	Name string `gorm:"unique" json:"categoly"`
}

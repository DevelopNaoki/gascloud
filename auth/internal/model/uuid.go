package model

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type UUID uuid.UUID

func (u *UUID) GormDataType() string {
	return "binary(16)"
}

func (u *UUID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "binary"
}

func (u *UUID) Scan(value any) (err error) {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannnot scan uuid")
	}
	parseByte, err := uuid.FromBytes(bytes)
	*u = UUID(parseByte)
	return
}

func (u UUID) Value() (bytes driver.Value, err error) {
	bytes, err = uuid.UUID(u).MarshalBinary()
	return
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

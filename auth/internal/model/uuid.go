package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type UUID uuid.UUID

func (u *UUID) GormDataType() string {
	return "varchar(36)"
}

func (u *UUID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "varchar(36)"
}

func (u *UUID) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		parsedUUID, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*u = UUID(parsedUUID)
	case []byte:
		parsedUUID, err := uuid.Parse(string(v))
		if err != nil {
			return err
		}
		*u = UUID(parsedUUID)
	default:
		return fmt.Errorf("cannot scan uuid from type %T", value)
	}
	return nil
}

func (u UUID) Value() (str driver.Value, err error) {
	return u.String(), nil
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func (u *UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *UUID) UnmarshalJSON(b []byte) error {
	uid, err := uuid.Parse(string(b))
	if err != nil {
		return err
	}
	*u = UUID(uid)
	return nil
}

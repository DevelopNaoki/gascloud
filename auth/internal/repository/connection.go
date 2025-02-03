package repository

import (
	"strconv"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectionDB(c model.DBConfig) (db *gorm.DB, err error) {
	dsn := c.User + ":" + c.Passwd + "@tcp(" + c.Host + ":" + strconv.Itoa(c.Port) + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	switch c.Driver {
	case "mysql", "mariadb":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return db, err
		}
	}

	db.AutoMigrate(
		&model.Account{},
		&model.GroupBind{},
		&model.Group{},
		&model.Session{},
		&model.Role{},
		&model.ServiceCatalog{},
		&model.ServiceToken{},
		&model.Categoly{},
	)

	err = initialData(db)
	if err != nil {
		return db, err
	}

	return db, nil
}

package repository

import (
	"fmt"
	"strconv"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectionDB(conf model.DBConfig) (db *gorm.DB, err error) {
	dsn := conf.User + ":" + conf.Passwd + "@tcp(" + conf.Host + ":" + strconv.Itoa(conf.Port) + ")/" + conf.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	switch conf.Driver {
	case "mysql", "mariadb":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			PrepareStmt: true,
		})
		if err != nil {
			return db, err
		}
	default:
		return db, fmt.Errorf("%s is not supported", conf.Driver)
	}

	if db == nil {
		return db, fmt.Errorf("failed open database")
	}
	err = db.AutoMigrate(
		&model.Account{},
		&model.GroupBind{},
		&model.Group{},
		&model.Session{},
		&model.Role{},
		&model.ServiceCatalog{},
		&model.ServiceToken{},
		&model.Categoly{},
	)
	if err != nil {
		return db, err
	}

	err = initialData(db)
	if err != nil {
		return db, err
	}

	return db, nil
}

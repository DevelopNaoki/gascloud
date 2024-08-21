package repository

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectionDB(c model.DBConfig) (db *gorm.DB, err error) {
	dsn := c.User + ":" + c.Pass + "@tcp(" + c.Host + ":" + strconv.Itoa(c.Port) + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	switch c.Driver {
	case "mysql", "mariadb":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return db, err
		}
	}
	db.AutoMigrate(&model.Account{}, &model.Role{}, &model.RoleBind{}, &model.Permission{}, &model.PermissionBind{}, &model.ServiceCatalog{})

	// Register Default User
	account := &model.Account{
		Name: "admin",
	}
	err = db.First(&account).Error
	if err == nil {
		return db, nil
	}

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(uint64(time.Now().UnixNano()))
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	hash, err := bcrypt.GenerateFromPassword(([]byte(string(b))), bcrypt.DefaultCost)
	if err != nil {
		return db, err
	}

	fmt.Printf("default administrator: admin:%s\n", string([]byte(string(b))))
	account = &model.Account{
		Name:        "admin",
		Passwd:      string(hash),
		Description: "Auto-registered administrator",
		IsAdmin:     true,
	}
	result := db.Create(&account)
	if result.Error != nil {
		return db, result.Error
	}

	return db, nil
}

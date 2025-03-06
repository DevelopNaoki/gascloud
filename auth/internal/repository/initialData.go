package repository

import (
	"fmt"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

func initialData(db *gorm.DB) (err error) {
	err = initialRole(db)
	if err != nil {
		return err
	}

	err = initialUser(db)
	if err != nil {
		return err
	}

	return nil
}

func initialUser(db *gorm.DB) error {
	// search admin user
	// if admin user exist, skip this process
	err := db.First(&model.Account{Name: "admin"}).Error
	if err == nil {
		return nil
	}

	// search admin role
	role := &model.Role{
		Name: "admin",
	}
	err = db.First(&role).Error
	if err != nil {
		return err
	}

	// create random password for length 10
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(uint64(time.Now().UnixNano()))
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	passwd := string(b)
	hash, err := bcrypt.GenerateFromPassword(([]byte(passwd)), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	fmt.Printf("default account: admin:%s\n", string([]byte(passwd)))
	account := &model.Account{
		Name:        "admin",
		Passwd:      string(hash),
		Role:        role.ID,
		Description: "Auto-Registered Administrator",
	}
	result := db.Create(&account)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func initialRole(db *gorm.DB) error {
	roles := []model.Role{
		{
			Name:        "admin",
			Description: "superuser account",
		},
		{
			Name:        "member",
			Description: "general accounts",
		},
		{
			Name:        "reader",
			Description: "readonly account",
		},
	}
	for _, role := range roles {
		err := db.Take(&role).Error
		if err == nil {
			continue
		}
		result := db.Create(&role)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

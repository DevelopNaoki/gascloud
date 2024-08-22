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
	// Register Default User
	err = initialUserData(db)
	if err != nil {
		return err
	}

	return nil
}

func initialUserData(db *gorm.DB) error {
	err := db.First(&model.Account{Name: "admin"}).Error
	if err == nil {
		return nil
	}

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(uint64(time.Now().UnixNano()))
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	hash, err := bcrypt.GenerateFromPassword(([]byte(string(b))), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	fmt.Printf("default account: admin:%s\n", string([]byte(string(b))))
	account := &model.Account{
		Name:        "admin",
		Passwd:      string(hash),
		Description: "Auto-registered administrator",
	}
	result := db.Create(&account)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

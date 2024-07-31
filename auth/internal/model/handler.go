package model

import "gorm.io/gorm"

type Handler struct {
	DB *gorm.DB
}

package handler

import (
	"gorm.io/gorm"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/cache"
)

type Handler struct {
	DB        *gorm.DB
	Cache     cache.Cache
	ExpiredAt time.Duration
}

func NewHandler(db *gorm.DB, cache cache.Cache, expiredAt time.Duration) *Handler {
	return &Handler{
		DB:        db,
		Cache:     cache,
		ExpiredAt: expiredAt,
	}
}

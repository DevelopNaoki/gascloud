package model

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UID      string `json:"uid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type ErrMsg struct {
	Message string `json:"message"`
}

package handler

import (
	"github.com/golang-jwt/jwt"
)

var JwtKey []byte = []byte("secret-jwt-key")

type Claims struct {
	Name string `json:"username"`
	jwt.StandardClaims
}

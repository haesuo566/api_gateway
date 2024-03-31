package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/novel/auth/config"
)

type JwtUtil struct {
}

const (
	ACCESS_TOKEN = iota
	REFRESH_TOKEN
)

var (
	instance IJwtUtil
	secret   string
)

func New() IJwtUtil {
	if instance == nil {
		secret = config.Getenv("JWT_SECRET")
		instance = &JwtUtil{}
	}
	return instance
}

func (j *JwtUtil) GenerateToken(tokenType int) (string, error) {
	sub := "access_token"
	iat := time.Now()
	exp := time.Now().Add(time.Minute * 30) // 30 minutes
	if tokenType == REFRESH_TOKEN {
		sub = "refresh_token"
		exp = time.Now().Add(time.Hour * 24 * 7) // week
	}

	claim := jwt.MapClaims{
		"sub": sub,
		"iat": iat,
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(secret)
}

func (j *JwtUtil) ValidateToken() {

}

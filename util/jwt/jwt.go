package jwt

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/novel/auth/config"
)

type JwtUtil struct {
}

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

func (j *JwtUtil) GenerateAccessToken(email string) (string, error) {
	claim := jwt.MapClaims{
		"sub":   "access_token",
		"iat":   time.Now(),
		"exp":   time.Now().Add(time.Minute * 30), // 30 minutes
		"email": email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(secret))
}

func (j *JwtUtil) GenerateRefreshToken() (string, error) {
	claim := jwt.MapClaims{
		"sub": "refresh_token",
		"iat": time.Now(),
		"exp": time.Now().Add(time.Hour * 24 * 7), // 1 week
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(secret))
}

func (j *JwtUtil) ValidateToken(token string) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// hash method check
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("not equals method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		log.Println(err)
		return err
	}

	// 나중에 claim 얻는거라든지 뭐 그런거 pointer map으로 받는걸로 추가 가능함
	return nil
}

package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/novel/api-gateway/config"
	"github.com/novel/api-gateway/entity/user"
)

type ResponseToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var secret string

func init() {
	config.LoadEnv()
	secret = os.Getenv("JWT_SECRET")
}

func GenerateResponseToken(user *user.User) (*ResponseToken, error) {
	accessToken, err := GenerateAccessToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	responseToken := &ResponseToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return responseToken, nil
}

func GenerateAccessToken() (string, error) {
	claim := jwt.MapClaims{
		"iat": time.Now(),
		"exp": time.Now().Add(time.Hour * 30),
		"sub": "access_token",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken() (string, error) {
	claim := jwt.MapClaims{
		"iat": time.Now(),
		"exp": time.Now().Add(time.Hour * 24 * 7),
		"sub": "refresh_token",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(secret))
}

func ValidateToken(token string) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token is invalidate")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	return nil
}

package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/novel/auth/config"
	"github.com/novel/auth/entity/user"
)

type ResponseToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateResponseToken(user *user.User) (*ResponseToken, error) {
	accessToken, err := GenerateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateToken(nil)
	if err != nil {
		return nil, err
	}

	responseToken := &ResponseToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return responseToken, nil
}

func GenerateToken(user *user.User) (string, error) {
	claim := jwt.MapClaims{
		"iat": time.Now(),
	}

	if user == nil {
		claim["exp"] = time.Now().Add(time.Hour * 24 * 7)
		claim["sub"] = "refresh_token"
	} else {
		claim["exp"] = time.Now().Add(time.Hour * 30)
		claim["sub"] = "access_token"
		claim["email"] = user.Email
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	config.LoadEnv()
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

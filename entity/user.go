package entity

import (
	"time"
)

const (
	GOOGLE = iota
	NAVER
	KAKAO
	NOVEL
)

type User struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Credential   *string   `json:"credential"`
	AccessToken  *string   `json:"oauth_access_token"`
	RefreshToken *string   `json:"oauth_refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Provider     int       `json:"provider"`
}

package user

import (
	"time"
)

const (
	GOOGLE = iota
	NAVER
	KAKAO
	NOVEL
)

// 연관테이블 JPA 마냥 넣어서 하면 좋을 듯
// pk 를 id, user, email 로 하면 name 을 중복으로 사용할 수 있긴 함
type User struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Credential *string   `json:"credential"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Provider   int       `json:"provider"`
}

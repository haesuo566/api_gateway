package user

import (
	"crypto/sha256"
	"fmt"
	"time"
)

const (
	GOOGLE = "GOOGLE"
	NAVER  = "NAVER"
	// KAKAO = "KAKAO"
	NOVEL = "NOVEL"
)

// 연관테이블 JPA 마냥 넣어서 하면 좋을 듯
// pk 를 id, user, email 로 하면 name 을 중복으로 사용할 수 있긴 함
type User struct {
	Id         int       `json:"id"`
	Email      string    `json:"email"`
	Credential *string   `json:"credential"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Provider   string    `json:"provider"`
}

func (user *User) HashCode() string {
	var credential string = ""
	if user.Credential != nil {
		credential = *user.Credential
	}

	hashString := fmt.Sprintf("%s|%s|%s|%s", user.Name, credential, user.Email, user.Provider)
	hash := sha256.New()
	hash.Write([]byte(hashString))

	return string(hash.Sum(nil))
}

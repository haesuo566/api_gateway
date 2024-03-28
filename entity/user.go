package entity

const (
	GOOGLE = iota
	NAVER
	KAKAO
)

type User struct {
	Id       string
	Platform int
}

func (u *User) Method() {

}

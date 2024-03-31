package entity

const (
	GOOGLE = iota
	NAVER
	KAKAO
	NOVEL
)

type User struct {
	Name string
}

func (u *User) Method() {

}

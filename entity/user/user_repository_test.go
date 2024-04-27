package user

import (
	"testing"

	"github.com/novel/auth/db"
)

func TestSave(t *testing.T) {
	db := db.New()
	repository := NewRepository(db)

	credential := "asdsadsadasd"
	user := &User{
		Name:       "asdasd",
		Email:      "asdsadsaddsdsa",
		Credential: &credential,
		Provider:   0,
	}

	if _, err := repository.Save(user); err != nil {
		t.Error(err)
	}
}

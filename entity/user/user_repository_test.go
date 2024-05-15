package user

import (
	"testing"
	"time"

	"github.com/novel/auth/db"
)

func TestSave(t *testing.T) {
	db := db.New()
	repository := NewRepository(db)

	credential := "dbstjdwls0129@@@"
	user := &User{
		Name:       "haesuo",
		Email:      "haesuo566@gmail.com",
		Credential: &credential,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Provider:   GOOGLE,
	}

	if _, err := repository.Save(user); err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	db := db.New()
	repository := NewRepository(db)

	user := &User{
		Name:      "haesuo",
		Email:     "haesuo566@gmail.com",
		UpdatedAt: time.Now().Local(),
	}

	if _, err := repository.Update(user, nil); err != nil {
		t.Error(err)
	}
}

func TestFindByEmail(t *testing.T) {
	db := db.New()
	repository := NewRepository(db)

	user := &User{
		Email:    "haesuo566@gmail.com",
		Provider: GOOGLE,
	}
	if _, err := repository.FindByEmailAndProvider(user, nil); err != nil {
		t.Error(err)
	}
}

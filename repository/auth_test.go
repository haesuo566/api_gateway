package repository

import (
	"testing"

	"github.com/novel/auth/entity"
)

func TestFindById(t *testing.T) {
	repo := NewAuthRepository()
	t.Run("execute query from user table", func(t *testing.T) {
		testId := "hso@trinitysoft.co.kr"
		user, err := repo.FindById(testId)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(user.Email)
	})
}

func TestSave(t *testing.T) {
	repo := NewAuthRepository()

	user := &entity.User{
		Name:       "김옥지",
		Email:      "김옥지@naver.com",
		Credential: "김옥지@@@",
		Provider:   1,
	}
	if err := repo.Save(user); err != nil {
		t.Fatal(err)
	}
}

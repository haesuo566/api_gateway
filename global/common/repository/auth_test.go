package repository

import (
	"fmt"
	"testing"

	"github.com/novel/auth/global/common/entity"
)

func TestFindById(t *testing.T) {
	repo := NewAuthRepository()
	t.Run("execute query from user table", func(t *testing.T) {
		testId := "hso@titysoft.co.kr"
		user, err := repo.FindById(testId)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(user.Email)
	})
}

func TestSave(t *testing.T) {
	repo := NewAuthRepository()

	access := "asddaadsasdasds"
	refresh := "asddaadsasdasasd"

	testUser := &entity.User{
		Name:         "김a아",
		Email:        "haoaaa1@trinitysoft.co.kr",
		AccessToken:  &access,
		RefreshToken: &refresh,
		Provider:     1,
	}

	user, err := repo.Save(testUser)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(*user)
}

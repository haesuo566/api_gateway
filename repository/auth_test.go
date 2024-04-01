package repository

import (
	"testing"
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

package google

import "github.com/novel/auth/entity"

type IGoogleRepository interface {
	FindById(id string) (*entity.User, error)
	Save(user *entity.User) (*entity.User, error)
}

type GoogleRepository struct {
}

var repositoryInstance IGoogleRepository = nil

func NewRepository() IGoogleRepository {
	if repositoryInstance == nil {
		repositoryInstance = &GoogleRepository{}
	}
	return repositoryInstance
}

func (g *GoogleRepository) FindById(id string) (*entity.User, error) {
	return nil, nil
}

func (g *GoogleRepository) Save(user *entity.User) (*entity.User, error) {
	return nil, nil
}

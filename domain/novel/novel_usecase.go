package novel

import "github.com/novel/auth/entity/user"

type INovelUsecase interface {
}

type NovelUsecase struct {
	userRepository user.IUserRepository
}

var usecaseInstance INovelUsecase = nil

func newUsecase(userRepository user.IUserRepository) INovelUsecase {
	if usecaseInstance == nil {
		usecaseInstance = &NovelUsecase{
			userRepository: userRepository,
		}
	}

	return usecaseInstance
}

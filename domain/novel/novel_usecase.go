package novel

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strings"

	"github.com/novel/auth/config"
	"github.com/novel/auth/entity/user"
)

type iNovelUsecase interface {
	login(string, string) (*user.User, error)
	singup(string, string, string) error
}

type novelUsecase struct {
	userRepository user.IUserRepository
}

var usecaseInstance iNovelUsecase = nil

func newUsecase(userRepository user.IUserRepository) iNovelUsecase {
	if usecaseInstance == nil {
		usecaseInstance = &novelUsecase{
			userRepository: userRepository,
		}
	}

	return usecaseInstance
}

func (n *novelUsecase) login(email, password string) (*user.User, error) {
	hashedPassword, err := hashString(password)
	if err != nil {
		return nil, err
	}

	user, err := n.userRepository.FindByEmail(email, nil)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(*user.Credential, hashedPassword) {
		return nil, errors.New("password is incorrect")
	}

	return user, nil
}

func (n *novelUsecase) singup(username, email, password string) error {
	hashedPassword, err := hashString(password)
	if err != nil {
		return err
	}

	user := &user.User{
		Name:       username,
		Email:      email,
		Credential: &hashedPassword,
		Provider:   1,
	}

	if _, err := n.userRepository.Save(user); err != nil {
		return err
	}

	return nil
}

func hashString(str string) (string, error) {
	config.LoadEnv()
	salt := os.Getenv("SALT")

	hash := sha256.New()
	if _, err := hash.Write([]byte(salt + str)); err != nil {
		return "", err
	}

	md := hash.Sum(nil)
	hashedPassword := hex.EncodeToString(md)

	return hashedPassword, nil
}

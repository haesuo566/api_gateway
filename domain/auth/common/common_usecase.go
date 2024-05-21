package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strings"

	"github.com/novel/api-gateway/config"
	"github.com/novel/api-gateway/entity/user"
)

type iCommonUsecase interface {
	login(string, string) (*user.User, error)
	singup(string, string, string) error
}

type commonUsecase struct {
	userRepository user.IUserRepository
}

var usecaseInstance iCommonUsecase = nil

func newUsecase(userRepository user.IUserRepository) iCommonUsecase {
	if usecaseInstance == nil {
		usecaseInstance = &commonUsecase{
			userRepository: userRepository,
		}
	}

	return usecaseInstance
}

func (n *commonUsecase) login(email, password string) (*user.User, error) {
	hashedPassword, err := hashString(password)
	if err != nil {
		return nil, err
	}

	user, err := n.userRepository.FindByEmailAndProvider(&user.User{
		Email:    email,
		Provider: user.NOVEL,
	}, nil)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(*user.Credential, hashedPassword) {
		return nil, errors.New("password is incorrect")
	}

	return user, nil
}

func (n *commonUsecase) singup(username, email, password string) error {
	hashedPassword, err := hashString(password)
	if err != nil {
		return err
	}

	user := &user.User{
		Name:       username,
		Email:      email,
		Credential: &hashedPassword,
		Provider:   user.NOVEL,
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

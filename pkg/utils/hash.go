package pkgutils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCannotGenerateHash = errors.New("cannot generate hash from password")
)

//go:generate mockery --name PasswordHashing --filename hash.go --output ./mocks
type PasswordHashing interface {
	Hash(password string) (string, error)
	CompareHash(hashedPassword, password string) bool
}

type passwordHashing struct {
}

func NewPasswordHashing() PasswordHashing {
	return &passwordHashing{}
}

func (p *passwordHashing) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrCannotGenerateHash
	}

	return string(hashBytes), nil
}

func (p *passwordHashing) CompareHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

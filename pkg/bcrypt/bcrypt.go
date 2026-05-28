package bcrypt

import (
	"github.com/BangNopall/hology8-be/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type BcryptInterface interface {
	Hash(plain string) (string, error)
	Compare(password, hashed string) bool
}

type BcryptStruct struct{}

var Bcrypt = getBcrypt()

func getBcrypt() BcryptInterface {
	return &BcryptStruct{}
}

func (b *BcryptStruct) Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[BCRYPT][Hash] failed to hash password")

		return "", err
	}

	return string(bytes), nil
}

func (b *BcryptStruct) Compare(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))

	return err == nil
}

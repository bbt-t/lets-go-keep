package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/bbt-t/lets-go-keep/internal/entity"
)

// GenerateRandom generates random bytes for encrypting.
func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// PasswordHash encryption.
func PasswordHash(credentials entity.UserCredentials) string {
	sha := sha256.New()
	sha.Write([]byte(credentials.Login + credentials.Password))

	return hex.EncodeToString(sha.Sum(nil))
}

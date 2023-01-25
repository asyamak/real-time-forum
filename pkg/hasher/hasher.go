package hasher

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

type Hasher interface {
	HashPassword(password string) string
}

type HasherService struct {
	salt string
}

func NewHasher(salt string) (*HasherService, error) {
	if salt == "" {
		return nil, errors.New("password salt is empty")
	}

	return &HasherService{salt: salt}, nil
}

func (h *HasherService) HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}

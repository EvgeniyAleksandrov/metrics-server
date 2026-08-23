package hash

import (
	"crypto/hmac"
	"crypto/sha256"
)

type SHA256 struct {
	key []byte
}

func NewSHA256(key string) *SHA256 {
	if key == "" {
		panic("Set empty key for sha256 hasher")
	}

	return &SHA256{
		key: []byte(key),
	}
}

func (s *SHA256) MakeHash(data []byte) []byte {
	hasher := hmac.New(sha256.New, s.key)
	hasher.Write(data)
	return hasher.Sum(nil)
}

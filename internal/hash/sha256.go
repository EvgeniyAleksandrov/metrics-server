package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
)

type SHA256 struct {
	hasher hash.Hash
}

func NewSHA256(key string) *SHA256 {
	if key == "" {
		panic("Set empty key for sha256 haser")
	}

	return &SHA256{
		hasher: hmac.New(sha256.New, []byte(key)),
	}
}

func (s *SHA256) MakeHash(data []byte) []byte {
	return s.hasher.Sum(data)
}

package helper

import (
	"crypto"
	"fmt"
)

type hash struct {
}

type Hash interface {
	MakeHash() string
	HashCompare(hashedValue string, comparator string) bool
}

func NewHash() *hash {
	return &hash{}
}

func (h *hash) MakeHash(value string) string {
	hashing := crypto.MD5.New()
	hashing.Write([]byte(value))
	hashed := hashing.Sum(nil)
	return fmt.Sprintf("%x", hashed)
}

func (h *hash) HashCompare(hashedValue string, comparator string) bool {
	return hashedValue == h.MakeHash(comparator)
}

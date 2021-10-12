package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
)

type Content struct {
	Size  int
	Hash  string
	Type  string
	Bytes []byte
}

func GenerateRandomContent(size int) (*Content, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	hash := fmt.Sprintf("%x", md5.Sum(bytes))
	return &Content{
		Size:  size,
		Hash:  hash,
		Type:  "application/octet-stream",
		Bytes: bytes,
	}, nil
}

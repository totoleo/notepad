package biz

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCipher(t *testing.T) {
	var key = make([]byte, 32)
	io.ReadFull(rand.Reader, key)

	original := []byte("hello world")
	encrypted, err := Encrypt(key, original)

	assert.NoError(t, err)

	planText, err := Decrypt(key, encrypted)

	assert.NoError(t, err)

	assert.Equal(t, original, planText)
}

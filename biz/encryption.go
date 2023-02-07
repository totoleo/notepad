package biz

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Decrypt(key, data []byte) ([]byte, error) {
	dataCap := base64.StdEncoding.DecodedLen(len(data))
	if dataCap <= 0 {
		return nil, nil
	}
	input := make([]byte, dataCap)
	n, err := base64.StdEncoding.Decode(input, data)
	if err != nil {
		return nil, err
	}
	data = input[:n]
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}
	c, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < c.NonceSize() {
		return nil, nil
	}

	nonce, encryptData := data[:c.NonceSize()], data[c.NonceSize():]

	return c.Open(nil, nonce, encryptData, nil)
}

func Encrypt(key []byte, plainData []byte) ([]byte, error) {

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}
	cipher, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	cipherLength := len(plainData) + cipher.Overhead()
	nonce := make([]byte, cipher.NonceSize(), cipher.NonceSize()+cipherLength)
	_, err = io.ReadFull(rand.Reader, nonce[:cipher.NonceSize()])
	if err != nil {
		return nil, err
	}
	nonce = cipher.Seal(nonce, nonce, plainData, nil)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(nonce)))
	base64.StdEncoding.Encode(dst, nonce)
	return dst, nil
}

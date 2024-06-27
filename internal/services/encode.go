package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type EncodingService interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

type EncodeService struct {
	encryptionKey string
}

func NewEncodeService(key string) *EncodeService {
	return &EncodeService{encryptionKey: key}
}

func (es EncodeService) Encrypt(data string) (string, error) {
	block, _ := aes.NewCipher([]byte(es.encryptionKey))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (es EncodeService) Decrypt(data string) (string, error) {
	key := []byte(es.encryptionKey)
	ciphertext, _ := base64.URLEncoding.DecodeString(data)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

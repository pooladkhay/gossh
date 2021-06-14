package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

func encryptPass(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)
	hashedKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashedKey[:])
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func decryptPass(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return "", err
	}

	hashedKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashedKey[:])
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

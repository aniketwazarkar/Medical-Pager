package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"medical-pager/utils"
)

// Encrypt encrypts plain text using AES-GCM
func Encrypt(plainText string) (string, error) {
	keyString := utils.GetEnv("ENCRYPTION_KEY", "32byte_supersecret_encryptionkey!")
	key := []byte(keyString)

	if len(key) != 32 {
		return "", errors.New("encryption key must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts encrypted text using AES-GCM
func Decrypt(encryptedHex string) (string, error) {
	keyString := utils.GetEnv("ENCRYPTION_KEY", "32byte_supersecret_encryptionkey!")
	key := []byte(keyString)

	if len(key) != 32 {
		return "", errors.New("encryption key must be 32 bytes")
	}

	enc, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

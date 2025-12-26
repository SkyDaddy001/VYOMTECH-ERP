package credentials

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

// EncryptionManager handles encryption/decryption of credentials
type EncryptionManager struct {
	masterKey string // AES-256 key (32 bytes when decoded)
}

// NewEncryptionManager creates a new encryption manager
func NewEncryptionManager(masterKey string) (*EncryptionManager, error) {
	if len(masterKey) == 0 {
		return nil, errors.New("master key cannot be empty")
	}
	return &EncryptionManager{
		masterKey: masterKey,
	}, nil
}

// Encrypt encrypts a credential object and returns base64 encoded result
func (em *EncryptionManager) Encrypt(credential interface{}) (string, error) {
	// Marshal credential to JSON
	jsonData, err := json.Marshal(credential)
	if err != nil {
		return "", err
	}

	// Decode master key
	decodedKey, err := base64.StdEncoding.DecodeString(em.masterKey)
	if err != nil {
		return "", errors.New("invalid master key format")
	}

	// Create cipher
	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", err
	}

	// Create GCM
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create nonce
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt
	ciphertext := aead.Seal(nonce, nonce, jsonData, nil)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a credential from base64 and returns unmarshalled object
func (em *EncryptionManager) Decrypt(encryptedData string, credential interface{}) error {
	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return errors.New("invalid encrypted data format")
	}

	// Decode master key
	decodedKey, err := base64.StdEncoding.DecodeString(em.masterKey)
	if err != nil {
		return errors.New("invalid master key format")
	}

	// Create cipher
	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return err
	}

	// Create GCM
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Extract nonce and ciphertext
	nonceSize := aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("decryption failed")
	}

	// Unmarshal JSON
	err = json.Unmarshal(plaintext, credential)
	if err != nil {
		return err
	}

	return nil
}

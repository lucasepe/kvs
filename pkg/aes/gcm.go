package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// GcmEncrypt implements the AES encryption with Galois/Counter Mode (AES-GCM)
func GcmEncrypt(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// generate a random nonce (makes encryption stronger)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted := gcm.Seal(nil, nonce, plainText, nil)
	// we need nonce for decryption so we put it at the beginning
	// of encrypted text
	return append(nonce, encrypted...), nil
}

// GcmDecrypt implements the AES decryption with Galois/Counter Mode (AES-GCM)
func GcmDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("Invalid data")
	}

	// extract random nonce we added to the beginning of the file
	nonce := ciphertext[:gcm.NonceSize()]
	encrypted := ciphertext[gcm.NonceSize():]

	return gcm.Open(nil, nonce, encrypted, nil)
}

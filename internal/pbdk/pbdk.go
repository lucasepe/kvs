package pbdk

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"io"
	mr "math/rand"

	"golang.org/x/crypto/pbkdf2"
)

// NewSalt create pseudo random binary data (salt)
// that is used as an additional input to derive the
// encryption key.
func NewSalt(secret []byte) ([]byte, error) {
	var seed int64
	binary.Read(bytes.NewBuffer(secret), binary.BigEndian, &seed)

	mr.Seed(seed)

	salt := make([]byte, 16)
	if _, err := mr.Read(salt); err != nil {
		return nil, err
	}

	return salt, nil
}

// NewEncryptionKey generates a random 256-bit key for Encrypt() and
// Decrypt(). It panics if the source of randomness fails.
func NewEncryptionKey() ([]byte, error) {
	res := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, res); err != nil {
		return nil, err
	}
	return res, nil
}

// DeriveKey derives a key as PBKDF2 from the specified secret string.
// HMAC is SHA 256, salt a random 16 bytes array, 2048 iterations.
// The returned key length is 32 bytes
func DeriveKey(secret []byte) ([]byte, error) {
	salt, err := NewSalt(secret)
	if err != nil {
		return nil, err
	}

	res := pbkdf2.Key(secret, salt, 2048, 32, sha256.New)
	return res, nil
}

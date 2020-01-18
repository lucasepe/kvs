package pbdk

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math/rand"

	"golang.org/x/crypto/pbkdf2"
)

func DeriveKey(secret []byte) ([]byte, error) {
	var seed int64
	if err := binary.Read(bytes.NewBuffer(secret), binary.BigEndian, &seed); err != nil {
		return nil, err
	}

	rand.Seed(seed)

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}
	//fmt.Printf("salt: %s\n", base64.StdEncoding.EncodeToString(salt))

	key := pbkdf2.Key([]byte(secret[:]), salt, 2048, 32, sha256.New)

	return key, nil
}

package signatureutil

import (
	"crypto/sha512"
	"encoding/hex"
)

func EncryptWithSalt(origin string) string {
	hash := sha512.New()
	hash.Write([]byte(origin))
	sha := hex.EncodeToString(hash.Sum([]byte("aA5L61J3bd#11k")))
	return sha
}

func Encrypt(origin string) string {
	hash := sha512.New()
	hash.Write([]byte(origin))
	sha := hex.EncodeToString(hash.Sum(nil))
	return sha
}

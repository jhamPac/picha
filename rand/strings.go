package rand

import (
	"crypto/rand"
	"encoding/base64"
)

// RememberTokenBytes is set to 32 bytes
const RememberTokenBytes = 32

// Bytes generates n random bytes
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String generates a byte slice of size nBytes and then return an encoded string
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken generates a remember token of a predefined size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}

// NBytes return the number of bytes used in any string generated by String or RememberToken
func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}

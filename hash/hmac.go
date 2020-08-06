package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
)

// HMAC is a custom wrapper for hmac
type HMAC struct {
	hmac hash.Hash
}

// NewHMAC instatiates a HMAC type with the provided key and sha256 as the function
func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))
	return HMAC{
		hmac: h,
	}
}

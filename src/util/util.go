package util

import (
	"crypto/sha512"
	"encoding/base64"
)

func GetSHA512(val string) string {
	// Compute the SHA-512 hash
	hash := sha512.New()
	hash.Write([]byte(val))

	// Get the final hashed output
	hashBytes := hash.Sum(nil)

	return base64.StdEncoding.EncodeToString(hashBytes)
}

func Ternary(cond bool, trueval any, falseval any) any {
	if cond {
		return trueval
	} else {
		return falseval
	}
}

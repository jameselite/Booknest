package helpers

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashString(str string) string {

	hash := sha256.Sum256([]byte(str))
	return base64.StdEncoding.EncodeToString(hash[:])
}
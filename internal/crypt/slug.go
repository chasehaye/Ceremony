package crypt


import (
	"crypto/rand"
    "encoding/hex"
)


func GenerateSlug() (string, error) {
    b := make([]byte, 6)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}
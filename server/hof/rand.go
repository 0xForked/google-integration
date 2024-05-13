package hof

import (
	"crypto/rand"
	"math/big"
)

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Generate a random string of the specified length
	randomString := make([]byte, length)
	for i := range randomString {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		randomString[i] = charset[n.Int64()]
	}
	return string(randomString), nil
}

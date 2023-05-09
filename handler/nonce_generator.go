package handler

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

func generateRandomString(length int) []byte {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, length)

	for idx := 0; idx < length; idx++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return []byte("")
		}

		ret[idx] = letters[num.Int64()]
	}

	return ret
}

func generateNonceString(length int) string {
	randomValue := generateRandomString(length)

	return base64.URLEncoding.EncodeToString(randomValue)
}

func defaultNonceGenerator() []byte {
	return []byte(generateNonceString(40))
}

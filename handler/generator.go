
package handler

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)


func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}

func generateNonceString() string {
        n := 50
	b := GenerateRandomString(n)
	return base64.URLEncoding.EncodeToString([]byte(b))
}

package shortener

import (
	"math/rand"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 7

func GenerateCode() string {
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(code)
}

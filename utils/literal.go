package utils

import "math/rand"

var char = "abcdefghiklmnopqrstuvxywzABCDEFGHIKLMNOPQRSTUVXYWZ"
var charLen = len(char)

func GenerateRandomString(n int) string {
	str := []byte{}
	for i := 0; i < n; i++ {
		str = append(str, char[rand.Intn(charLen)])
	}
	return string(str)
}

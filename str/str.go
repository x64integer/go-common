package str

import (
	"encoding/base64"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// UUID will return cryptographically-secure unique id
func UUID() string {
	return uuid.NewV4().String()
}

// Random - generate random string using masking with source
//
// Credits to: https://medium.com/@kpbird/golang-generate-fixed-size-random-string-dd6dbd5e63c0
func Random(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	l := len(letterBytes)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < l {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// Base64URLEncode will base64 URL encode given input
func Base64URLEncode(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

// Base64URLDecode will decode base64 URL encoded given input
func Base64URLDecode(input string) (string, error) {
	res, err := base64.URLEncoding.DecodeString(input)

	return string(res), err
}

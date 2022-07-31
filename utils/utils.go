package utils

import (
	"crypto/sha256"
	"math/big"
	"net/url"
)

func GenerateShortLink(originalLink string) string {

	sha256 := sha256.New()
	sha256.Write([]byte(originalLink))
	hashBytes := sha256.Sum(nil)

	base64String := toBase62(hashBytes)

	return base64String[:8]

}

func toBase62(b []byte) string {
	var i big.Int
	i.SetBytes(b)
	return i.Text(62)
}

func CheckIfURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

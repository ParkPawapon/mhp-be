package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

func RandomDigits(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length")
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result[i] = byte('0' + n.Int64())
	}
	return string(result), nil
}

func RandomRefCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length")
	}
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf)[:length], nil
}

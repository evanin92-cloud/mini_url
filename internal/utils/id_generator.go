package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const Base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const DefaultLength = 7

func GenerateRandomID(length int) (string, error) {
	if length <= 0 {
		length = DefaultLength
	}

	var builder strings.Builder
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(Base62Chars))))
		if err != nil {
			return "", err
		}
		builder.WriteByte(Base62Chars[num.Int64()])
	}
	return builder.String(), nil
}

func EncodeToBase62(num int64) string {
	if num == 0 {
		return string(Base62Chars[0])
	}

	var result strings.Builder
	base := int64(len(Base62Chars))

	for num > 0 {
		remainder := num % base
		result.WriteByte(Base62Chars[remainder])
		num = num / base
	}

	return reverseString(result.String())
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
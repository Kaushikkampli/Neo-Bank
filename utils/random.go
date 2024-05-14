package utils

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomString generates a random integer between min & max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var result strings.Builder
	for i := 0; i < n; i++ {
		randomIndesx := rand.Intn(len(alphabet))
		result.WriteByte(alphabet[randomIndesx])
	}
	return result.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, INR}

	return currencies[rand.Intn(len(currencies))]
}

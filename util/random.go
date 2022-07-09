package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghujklmnopqrstuvwxyz"

// This init function is called everytime the package is first used.
func init()  {
	rand.Seed(time.Now().UnixNano())
}

// Generate a random int between min and max
func RandomInt(min, max int64) int64{
	return min + rand.Int63n(max-min+1)
}

// Generate a random string of n characters
func RandomString(n int) string{
	var sb strings.Builder
	k := len(alphabet)

	for i:=0; i < n; i++{
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Get owner name of length 6.
func RandomOwner() string{
	return RandomString(6)
}

func RandomMoney() int64{
	return RandomInt(0, 1000)
}

func RandomCurrency()string{
	currencies := []string{USD, EUR, CAD}

	return currencies[rand.Intn(len(currencies))]
}
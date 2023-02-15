package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init(){
	rand.Seed(time.Now().UnixNano())
}

// generates a random integer between min and max
func RandomInt(min, max int64) int64{

	// rand.Int63n will return a random int between 0 and (max-min), when you add min, it will return 
	// int between min and max
	return min + rand.Int63n(max - min + 1) 
}

// generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	// build string of length n
	for i:=0; i < n; i++{
		c:= alphabet[rand.Intn(k)] // get position of char 
		sb.WriteByte(c) // write char c to string builder
	}

	return sb.String()
}

// generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// generates a random amount of money
func RandomMoney() int64{
	return RandomInt(0, 1000)
}

// generates a random currency code
func RandomCurrency() string{
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)

	return currencies[rand.Intn(n)] // return currency from the index of the slice
}
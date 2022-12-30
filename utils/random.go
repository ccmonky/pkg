package pkg

import (
	"math/rand"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	src         = rand.NewSource(time.Now().UnixNano())
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Random return a integer with in [min, max)
func Random(min int, max int) int {
	return rand.Intn(max-min) + min
}

// WeightedChoice used for weighted random selection.
// If weights is in descending order, runtime is further reduced by about 20%, suitable for dynamic weights.
func WeightedChoice(weights []float32) int {
	var sum float32
	for _, num := range weights {
		sum += num
	}
	rnd := rand.Float32() * sum
	for i, w := range weights {
		rnd -= w
		if rnd < 0 {
			return i
		}
	}
	return 0
}

// Shuffle returns a random distributed slice
func Shuffle(slice []int) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// RandStringRunes - generate random string using random int
func RandStringRunes(n int) string {
	b := make([]rune, n)
	l := len(letterRunes)
	for i := range b {
		b[i] = letterRunes[rand.Intn(l)]
	}
	return string(b)
}

// RandStringBytes - generate random string using bytes
func RandStringBytes(n int) string {
	b := make([]byte, n)
	l := len(letterBytes)
	for i := range b {
		b[i] = letterBytes[rand.Intn(l)]
	}
	return string(b)
}

//RandStringBytesRmndr - generate random string using Remainder
func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	l := len(letterBytes)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(l)]
	}
	return string(b)
}

//RandStringBytesMask - generate random string using masking
func RandStringBytesMask(n int) string {
	b := make([]byte, n)
	l := len(letterBytes)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < l {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

//RandStringBytesMaskImpr - generate random string using masking improved
func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	l := len(letterBytes)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
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

//RandStringBytesMaskImprSrc - generate random string using masking with source
func RandStringBytesMaskImprSrc(n int) string {
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

//RandASCIIBytes - A helper function create and fill a slice of length n with characters from a-zA-Z0-9_-. It panics if there are any problems getting random bytes.
func RandASCIIBytes(n int) []byte {
	output := make([]byte, n)

	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)

	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}

	l := len(letterBytes)
	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])

		// random % 64
		randomPos := random % uint8(l)

		// put into output
		output[pos] = letterBytes[randomPos]
	}

	return output
}

//RandomString - Generate a random string of A-Z chars with len = l
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return string(bytes)
}

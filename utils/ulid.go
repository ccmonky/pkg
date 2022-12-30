package pkg

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropyPool = sync.Pool{
	New: func() interface{} {
		seed := time.Now().UnixNano()
		source := rand.NewSource(seed)
		return rand.New(source)
	},
}

// Ulid 生成Ulid返回string
func Ulid() (string, error) {
	entropy := entropyPool.Get().(*rand.Rand)
	defer entropyPool.Put(entropy)
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// MustUlid 生成Ulid返回string，失败panic
func MustUlid() string {
	s, err := Ulid()
	if err != nil {
		panic(err)
	}
	return s
}

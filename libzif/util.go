package zif

import (
	"crypto/rand"
	"math/big"
)

func CryptoRandBytes(size int) ([]byte, error) {
	buf := make([]byte, size)
	_, err := rand.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func CryptoRandInt(min, max int64) int64 {
	num, err := rand.Int(rand.Reader, big.NewInt(max-min))

	if err != nil {
		panic(err)
	}

	return num.Int64() + min
}

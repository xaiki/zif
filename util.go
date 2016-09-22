package zif

import (
	"crypto/rand"
	"math/big"
	"time"
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

func NewLimiter(rate time.Duration, burst int) (chan time.Time, *time.Ticker) {
	tick := time.NewTicker(rate)
	throttle := make(chan time.Time, burst)

	go func() {
		for t := range tick.C {
			select {
			case throttle <- t:
			default:
			}
		}
	}()

	return throttle, tick
}

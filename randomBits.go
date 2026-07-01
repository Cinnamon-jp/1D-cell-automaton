package main

import (
	"errors"
	"math/rand/v2"
)

// randomBits は、ランダムなビット列を生成する関数
func randomBits(r *rand.Rand, length int) ([]byte, error) {
	var bits []byte

	if length < 0 {
		return nil, errors.New("length must be positive value")
	}
	bits = make([]byte, length)
	for i := range bits {
		bits[i] = byte(r.IntN(2))
	}

	return bits, nil
}

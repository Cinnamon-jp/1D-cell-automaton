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

// updateBits は、アップデートルールに従ってビット列を更新する関数
func updateBits(oldBits []byte, rule byte) ([]byte, error) {
	bitLength := len(oldBits)

	newBits := make([]byte, bitLength)
	var target3Bits byte

	for i := range oldBits {
		if i != 0 {
			target3Bits = oldBits[i-1]
		} else {
			target3Bits = oldBits[bitLength-1]
		}

		target3Bits = target3Bits<<1 | oldBits[i]

		if i != bitLength-1 {
			target3Bits = target3Bits<<1 | oldBits[i+1]
		} else {
			target3Bits = target3Bits<<1 | oldBits[0]
		}

		// ruleの中から target3Bits に対応するビットを取り出し、次の状態とする
		newBits[i] = (rule >> target3Bits) & 1
	}

	return newBits, nil
}

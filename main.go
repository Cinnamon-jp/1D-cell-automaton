package main

import (
	"errors"
	"math/rand/v2"
)

func main() {
}

func run() {
	const SEED_1 uint64 = 85671395
	const SEED_2 uint64 = 18429642

	r := rand.New(rand.NewPCG(SEED_1, SEED_2))
	
}

// randomBits は、ランダムな100ビットを生成する関数
func randomBits(r *rand.Rand) (bits []byte) {
	bits = make([]byte, 100)
	for i, _ := range bits {
		bits[i] = byte(r.IntN(2))
	}
	return
}


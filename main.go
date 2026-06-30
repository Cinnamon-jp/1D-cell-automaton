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

// updateBits は、アップデートルールに従ってビット列を更新する関数
func updateBits(oldBits []byte, rule byte) (newBits []byte, err error) {
	newBits = make([]byte, 100)
	var target3Bits byte

	// ビット列の長さチェック
	if len(oldBits) != 100 {
		return nil, errors.New("length of oldBits must be 100")
	}
	
	for i, _ := range oldBits {
		if i != 0 {
			target3Bits = oldBits[i-1]
		} else {
			target3Bits = oldBits[99]
		}
		
		target3Bits = target3Bits << 1 | oldBits[i]

		if i != 99 {
			target3Bits = target3Bits << 1 | oldBits[i+1]
		} else {
			target3Bits = target3Bits << 1 | oldBits[0]
		}

		// ruleの中から target3Bits に対応するビットを取り出し、次の状態とする
		newBits[i] = (rule >> target3Bits) & 1
	}

	return newBits, nil
}
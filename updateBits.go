package main

import (
	"errors"
	"fmt"
)

// updateBits は、アップデートルールに従ってビット列を更新する関数
func updateBits(oldBits []byte, rule byte, length int) ([]byte, error) {
	if length < 0 {
		return nil, errors.New("length must be positive value")
	}
	newBits := make([]byte, length)
	var target3Bits byte

	// ビット列の長さチェック
	if len(oldBits) != length {
		return nil, fmt.Errorf("length of oldBits must be %d", length)
	}

	for i := range oldBits {
		if i != 0 {
			target3Bits = oldBits[i-1]
		} else {
			target3Bits = oldBits[length-1]
		}

		target3Bits = target3Bits<<1 | oldBits[i]

		if i != length-1 {
			target3Bits = target3Bits<<1 | oldBits[i+1]
		} else {
			target3Bits = target3Bits<<1 | oldBits[0]
		}

		// ruleの中から target3Bits に対応するビットを取り出し、次の状態とする
		newBits[i] = (rule >> target3Bits) & 1
	}

	return newBits, nil
}

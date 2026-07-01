package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
)

func main() {
}

func run() (error) {
	var initMode int = -1
	for {
		initMode, _ = scanInput("初期値 { 自動生成(0) / 手動(1) } : ")
		if initMode == 0 || initMode == 1 {
			break
		}
		fmt.Println("入力が不正です")
	}

	var init100Bits = make([]byte, 100)

	if initMode == 0 {
		const SEED_1 uint64 = 42
		const SEED_2 uint64 = 18429642

		r := rand.New(rand.NewPCG(SEED_1, SEED_2))

		init100Bits = randomBits(r)

		fmt.Println("ランダムな100ビットを生成しました")
	} else {
		buf, err := os.ReadFile("./init100Bits")
		if err != nil {
			return err
		}

		if len(buf) != 100 {
			return errors.New("length of init100Bits must be 100")
		}
		
		for i, v := range buf {
			if v != '0' && v != '1' {
				return errors.New("value of init100Bits must be 0 or 1")
			}
			init100Bits[i] = v - '0'
		}

		fmt.Println("./init100Bits を読み取りました")
	}

	
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

// scanInput は、プロンプトを表示し、入力を文字列として取得する関数
func scanInput(prompt string) (input int, err error) {
	fmt.Print(prompt)
	_, err = fmt.Scan(&input)

	return input, err
}
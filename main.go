package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	var initMode int = -1
	var ruleNumber int = -1
	var byteLength int = 100
	var err error

	for initMode != 0 && initMode != 1 {
		initMode, _ = scanInput("初期値 { 自動生成(0) / 手動(1) } : ")
		if initMode != 0 && initMode != 1 {
			fmt.Println("入力が不正です")
		}
	}

	for ruleNumber < 0 || ruleNumber > 255 {
		ruleNumber, err = scanInput("ルール {0 ~ 255} : ")
		if err != nil || ruleNumber < 0 || ruleNumber > 255 {
			fmt.Println("入力が不正です")
		}
	}

	for byteLength <= 0 || byteLength > 1000 {
		byteLength, err = scanInput("ビット長 {1 ~ 1000} (デフォルト 100): ")
		if byteLength == "" {

		}
		if err != nil || byteLength <= 0 || byteLength > 1000 {
			fmt.Println("入力が不正です")
		}
	}

	var initBits = make([]byte, byteLength)

	if initMode == 0 {
		const SEED_1 uint64 = 42
		const SEED_2 uint64 = 18429642

		r := rand.New(rand.NewPCG(SEED_1, SEED_2))

		initBits, err = randomBits(r, byteLength)
		if err != nil {
			return err
		}

		fmt.Println("ランダムなビット列を生成しました")
	} else {
		buf, err := os.ReadFile("./initBits")
		if err != nil {
			return err
		}

		if len(buf) != byteLength {
			return fmt.Errorf("length of initBits must be %d", byteLength)
		}

		for i, v := range buf {
			if v != '0' && v != '1' {
				return errors.New("value of initBits must be 0 or 1")
			}
			initBits[i] = v - '0'
		}

		fmt.Println("./initBits を読み取りました")
	}

	return nil
}

// randomBits は、ランダムなビット列を生成する関数
func randomBits(r *rand.Rand, length int) (bits []byte, err error) {
	if length < 0 {
		return nil, errors.New("length must be positive value")
	}
	bits = make([]byte, length)
	for i, _ := range bits {
		bits[i] = byte(r.IntN(2))
	}
	return
}

// updateBits は、アップデートルールに従ってビット列を更新する関数
func updateBits(oldBits []byte, rule byte, length int) (newBits []byte, err error) {
	if length < 0 {
		return nil, errors.New("length must be positive value")
	}
	newBits = make([]byte, length)
	var target3Bits byte

	// ビット列の長さチェック
	if len(oldBits) != length {
		return nil, fmt.Errorf("length of oldBits must be %d", length)
	}

	for i, _ := range oldBits {
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

// scanInput は、プロンプトを表示し、入力を文字列として取得する関数
func scanInput(prompt string) (input int, err error) {
	fmt.Print(prompt)

	var strInput string
	_, err = fmt.Scan(&strInput)
	if err != nil {
		return -1, err
	}

	input, err = strconv.Atoi(strInput)

	return input, err
}

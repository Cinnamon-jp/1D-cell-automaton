package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	// ハイパーパラメータ設定
	// --------------------------------------------------
	// ビット長
	bitLength, err := getVal("ビット長 (デフォルト 100) : ", 3, 1000, 100)
	if err != nil {
		return err
	}

	// 初期化方法
	initMode, err := getVal("初期化方法 {自動 0 / 手動 1} (デフォルト 自動): ", 0, 1, 0)
	if err != nil {
		return err
	}

	// ルールナンバー
	rules, err := getVal("ルールナンバー: ", 0, 255, -1)
	if err != nil {
		return err
	}

	// 状態の初期化
	// --------------------------------------------------

	var initBits []byte

	if initMode == 0 {
		const SEED_1 uint64 = 42
		const SEED_2 uint64 = 85653950

		r := rand.New(rand.NewPCG(SEED_1, SEED_2))

		initBits, err = randomBits(r, bitLength)
	} else {
		initBitsStr, err := os.ReadFile("initBits")
		if err != nil {
			return err
		}

		// トリム
		initBitsStr = bytes.TrimSpace(initBitsStr)

		// 文字数が指定と一致しているかチェック
		if len(initBitsStr) != bitLength {
			return fmt.Errorf("./initBits の文字数 (%d) が bitLength (%d) と異なります", len(initBitsStr), bitLength)
		}

		initBits = make([]byte, bitLength)
		for i, v := range initBitsStr {
			if v == '0' || v == '1' {
				initBits[i] = byte(v - '0')
			} else {
				return errors.New("./initBits に 0 1 以外の文字が含まれています")
			}
		}
	}

	return nil
}

// getVal は、min から max までの範囲で val を取得する関数
func getVal(prompt string, min int, max int, defVal int) (int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return -1, err
			}
			return -1, errors.New("入力が途絶えました") // EOF
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			if defVal == -1 {
				fmt.Println("入力が必要です")
				continue
			}
			return defVal, nil
		}

		val, err := strToInt(input)
		if err != nil {
			fmt.Println("入力が数値ではありません")
			continue
		}
		if val < min || val > max {
			fmt.Printf("入力が範囲外です min: %d, max: %d\n", min, max)
			continue
		}
		return val, nil
	}
}

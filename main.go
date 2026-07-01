package main

import (
	"fmt"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	// ハイパーパラメータ設定
	// ビット長
	bitLength, err := getVal("ビット長 (デフォルト 100) : ", 3, 1000, 100)
	if err != nil {
		return err
	}
	// 初期化方法
	initMode, err := getVal("初期化方法 {自動 0 / 手動 1} (デフォルト 0): ", 0, 1, 0)
	// ルールナンバー
}

// getVal は、min から max までの範囲で val を取得する関数
func getVal(prompt string, min int, max int, defVal int) (int, error) {
	var val int = defVal

	for {
		input, err := scanInput(prompt)
		if err != nil {
			return -1, err
		}

		// デフォルト値を使用する場合
		if input == "" {
			return val, nil
		}

		val, err = strToInt(input)
		if err != nil {
			return -1, err
		}
		if val < min || val > max {
			fmt.Printf("入力が不正です min: %d, max: %d\n", min, max)
			continue
		}

		break
	}

	return val, nil
}

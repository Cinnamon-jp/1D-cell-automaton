package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func scanInput(prompt string) (string, error) {
	var input string

	fmt.Print(prompt)
	_, err := fmt.Scan(&input)
	if err != nil {
		return "", err
	}
	return input, nil
}

func strToInt(s string) (int, error) {
	return strconv.Atoi(s)
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
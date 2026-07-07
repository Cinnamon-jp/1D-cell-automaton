package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

const (
	SEED_1 = 2669
	SEED_2 = 175924051
)

var r = rand.New(rand.NewPCG(SEED_1, SEED_2))

// selector は、promptを表示してsetの要素を矢印キーで選択させる関数。
// ↑/↓キーでカーソルを移動し、Enterキーで選択を確定する。
func selector(set []string, prompt string) string {
	if len(set) == 0 {
		return ""
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// rawモードにできない場合はフォールバック
		fmt.Print(prompt)
		var input string
		fmt.Scan(&input)
		return input
	}
	defer term.Restore(fd, oldState)

	cursor := 0

	// 選択肢を描画するヘルパー
	draw := func() {
		// 描画済み行を消去してカーソルを先頭へ戻す
		// prompt行 + len(set)行 分を消去
		lines := 1 + len(set)
		for i := 0; i < lines; i++ {
			fmt.Print("\033[A\033[2K")
		}
		fmt.Printf("%s\r\n", prompt)
		for i, item := range set {
			if i == cursor {
				fmt.Printf("\033[36m> %s\033[0m\r\n", item)
			} else {
				fmt.Printf("  %s\r\n", item)
			}
		}
	}

	// 初回描画（消去なしで描画）
	fmt.Printf("%s\r\n", prompt)
	for i, item := range set {
		if i == cursor {
			fmt.Printf("\033[36m> %s\033[0m\r\n", item)
		} else {
			fmt.Printf("  %s\r\n", item)
		}
	}

	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		if n == 0 {
			continue
		}

		switch {
		case n == 1 && buf[0] == '\r': // Enter
			return set[cursor]
		case n == 1 && buf[0] == 3: // Ctrl+C
			term.Restore(fd, oldState)
			os.Exit(0)
		case n == 3 && buf[0] == 27 && buf[1] == '[' && buf[2] == 'A': // ↑
			if cursor > 0 {
				cursor--
				draw()
			}
		case n == 3 && buf[0] == 27 && buf[1] == '[' && buf[2] == 'B': // ↓
			if cursor < len(set)-1 {
				cursor++
				draw()
			}
		}
	}

	return set[cursor]
}

func hexToBits(hex string, length int) ([]byte, error) {
	if hex == "?" {
		bits, _ := randomBits(r, length)
		return bits, nil
	}

	var bits []byte
	for _, ch := range hex {
		val, err := strconv.ParseInt(string(ch), 16, 8)
		if err != nil {
			return nil, err
		}

		for i := 3; i >= 0; i-- {
			bit := (val >> i) & 1
			bits = append(bits, byte(bit))
		}
	}
	return bits, nil
}

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

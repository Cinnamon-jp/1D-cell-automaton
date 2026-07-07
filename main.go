package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	file, err := os.Open("expCond.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// ヘッダー行を読み飛ばす
	if _, err := reader.Read(); err != nil {
		return err
	}

	var eg errgroup.Group

	// 1行ずつ処理
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		cond := record
		eg.Go(func() error {
			return runExp(cond)
		})
	}

	return eg.Wait()
}

// bitsToHalfBlock は、2行分のビット列を半ブロック文字の文字列に変換する。
// upperBits が上半分、lowerBits が下半分を表す。
// lowerBits が nil の場合は下半分をすべてOFFとして処理する。
func bitsToHalfBlock(upperBits, lowerBits []byte) string {
	var sb strings.Builder
	for i, upper := range upperBits {
		var lower byte
		if lowerBits != nil {
			lower = lowerBits[i]
		}
		switch {
		case upper == 0 && lower == 0:
			sb.WriteRune(' ')
		case upper == 1 && lower == 0:
			sb.WriteRune('▀')
		case upper == 0 && lower == 1:
			sb.WriteRune('▄')
		case upper == 1 && lower == 1:
			sb.WriteRune('█')
		}
	}
	return sb.String()
}

func runExp(cond []string) error {
	// ハイパーパラメータの取得
	bitLength, err := strconv.Atoi(cond[0])
	if err != nil {
		return err
	}
	rule, err := strconv.Atoi(cond[1])
	if err != nil {
		return err
	}
	step, err := strconv.Atoi(cond[2])
	if err != nil {
		return err
	}
	initBits, err := hexToBits(cond[3], bitLength)
	if err != nil {
		return err
	}

	// ハイパーパラメータのチェック
	if bitLength < 3 || bitLength > 1000 {
		return fmt.Errorf("bitLength の値域が不正: %d", bitLength)
	}
	if rule < 0 || rule > 255 {
		return fmt.Errorf("rule の値域が不正: %d", rule)
	}
	if step < 0 {
		return fmt.Errorf("step の値域が不正: %d", step)
	}
	if len(initBits) != bitLength {
		return fmt.Errorf("hexInitBits の長さが bitLength と一致しません")
	}

	// 全ステップのビット列を蓄積（初期ビット + step回分）
	allRows := make([][]byte, 0, step+1)
	allRows = append(allRows, initBits)

	oldBits := initBits
	for range step {
		newBits, err := updateBits(oldBits, byte(rule))
		if err != nil {
			return err
		}
		allRows = append(allRows, newBits)
		oldBits = newBits
	}

	// 結果保存ファイルの作成
	os.MkdirAll("result", 0o755)
	fileName := fmt.Sprintf("result/%dbits_rule%d_%dsteps.txt", bitLength, rule, step)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 2行ずつペアにして半ブロック文字で書き込み
	for i := 0; i < len(allRows); i += 2 {
		var lowerBits []byte
		if i+1 < len(allRows) {
			lowerBits = allRows[i+1]
		}
		fmt.Fprintln(file, bitsToHalfBlock(allRows[i], lowerBits))
	}

	return nil
}

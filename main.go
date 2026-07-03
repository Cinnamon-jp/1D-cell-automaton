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

	// 結果保存ファイルの作成
	os.MkdirAll("result", 0o755)
	fileName := fmt.Sprintf("result/%dbits_rule%d_%dsteps.txt", bitLength, rule, step)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 実験の実行
	oldBits := initBits
	var line strings.Builder

	// 初期ビットのファイルへの書き込み
	for _, v := range oldBits {
		if v == 0 {
			line.WriteString("░")
		} else {
			line.WriteString("█")
		}
	}
	fmt.Fprintln(file, line.String())

	for range step {
		newBits, err := updateBits(oldBits, byte(rule))
		if err != nil {
			return err
		}

		line.Reset()

		// ファイルへの書き込み
		for _, v := range newBits {
			if v == 0 {
				line.WriteString("░")
			} else {
				line.WriteString("█")
			}
		}
		fmt.Fprintln(file, line.String())

		oldBits = newBits
	}

	return nil
}

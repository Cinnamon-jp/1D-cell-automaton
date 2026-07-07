package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
)

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
	}
}

func run() error {
	c1 := []uint8{0, 8, 32, 40, 64, 96, 128, 136, 160, 168, 192, 224, 234, 235, 238, 239, 248, 249, 250, 251, 252, 253, 254, 255}
	c2 := []uint8{4, 12, 13, 36, 44, 68, 69, 72, 76, 77, 78, 79, 92, 93, 100, 104, 132, 140, 141, 164, 172, 196, 197, 200, 202, 203, 204, 205, 206, 207, 216, 217, 218, 219, 220, 221, 222, 223, 228, 232, 233, 236, 237}
	c3 := []uint8{5, 28, 29, 70, 71, 73, 94, 95, 108, 109, 133, 156, 157, 198, 199, 201}
	c4 := []uint8{1, 2, 3, 6, 7, 9, 10, 11, 14, 15, 16, 17, 19, 20, 21, 23, 24, 25, 26, 27, 31, 33, 34, 35, 37, 38, 39, 41, 42, 43, 46, 47, 48, 49, 50, 51, 52, 53, 55, 56, 57, 58, 59, 61, 62, 63, 65, 66, 67, 74, 80, 81, 82, 83, 84, 85, 87, 88, 91, 97, 98, 99, 103, 107, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 121, 123, 124, 125, 127, 130, 131, 134, 137, 138, 139, 142, 143, 144, 145, 148, 152, 154, 155, 158, 159, 162, 163, 166, 167, 170, 171, 173, 174, 175, 176, 177, 178, 179, 180, 181, 184, 185, 186, 187, 188, 189, 190, 191, 193, 194, 208, 209, 210, 211, 212, 213, 214, 215, 226, 227, 229, 230, 231, 240, 241, 242, 243, 244, 245, 246, 247}
	c5 := []uint8{54, 147}
	c6 := []uint8{18, 22, 30, 45, 60, 75, 86, 89, 90, 101, 102, 105, 106, 120, 122, 126, 129, 135, 146, 149, 150, 151, 153, 161, 165, 169, 182, 183, 195, 225}

	slices.Sort(c1)
	slices.Sort(c2)
	slices.Sort(c3)
	slices.Sort(c4)
	slices.Sort(c5)
	slices.Sort(c6)

	totalLen := len(c1) + len(c2) + len(c3) + len(c4) + len(c6) + len(c5)

	// ルール番号網羅性チェック
	if totalLen != 256 {
		return errors.New("element count not 256")
	}

	all_c := make([]uint8, 0, totalLen)

	all_c = append(all_c, c1...)
	all_c = append(all_c, c2...)
	all_c = append(all_c, c3...)
	all_c = append(all_c, c4...)
	all_c = append(all_c, c6...)
	all_c = append(all_c, c5...)

	slices.Sort(all_c)

	// 重複チェック
	if checkDup(all_c) {
		return errors.New("element duplication")
	}

	if err := fileSlice(decToBin(all_c), "all_c.txt"); err != nil {
		return fmt.Errorf("create file: %s", err)
	}

	return nil
}

// checkDup はスライスを受け取って、要素に重複が無いかチェックする
func checkDup(s []uint8) (isDup bool) {
	var seen [256]bool
	for _, val := range s {
		if seen[val] {
			return true
		}
		seen[val] = true
	}
	return false
}

// decToBin は受け取った []uint8 の全ての要素を8ケタの2進数の文字列にして返す
func decToBin(s []uint8) []string {
	binSlice := make([]string, 0, len(s))
	for _, v := range s {
		binSlice = append(binSlice, fmt.Sprintf("%08b", v))
	}
	return binSlice
}

// fileSlice は "0", "1"で構成された []string をファイルに保存する
func fileSlice(s []string, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()


	var builder strings.Builder
	for _, binStr := range s {
		line := strings.ReplaceAll(binStr, "0", "  ")
		line = strings.ReplaceAll(line, "1", "██")
		
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	_, err = f.WriteString(builder.String())
	return err
}

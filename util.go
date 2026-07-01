package main

import (
	"fmt"
	"strconv"
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

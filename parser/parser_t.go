package parser

import (
	"fmt"
	"os"
)

func Check() {
	file, err := os.Open("parser/input.test")
	if err != nil {
		panic(err)
	}

	prc, _ := Parse(file)

	fmt.Println(prc)

	LexAndPrintTokens(file)
}

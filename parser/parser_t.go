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
	// LexAndPrintTokens(file)

	prc, err := Parse(file)

	switch {
	case err != nil:
		fmt.Println("Parsing error: ")
		fmt.Println(err)
	default:
		processes = expandUnexpandedProcesses(prc)

		for _, p := range processes {
			fmt.Println(p.Body.String())
			fmt.Println(len(p.FunctionDefinitions))
		}
	}
}

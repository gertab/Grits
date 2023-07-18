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

	// prc, err := parser.Parse(file)

	// fmt.Println(prc)
	lexer := newLexer(file)
	val := phiSymType{}
	for {

		tok := lexer.Lex(&val)
		if tok == EOF {
			break
		}

		fmt.Printf("\t%d\t%s\n", tok, val.strval)
	}
}

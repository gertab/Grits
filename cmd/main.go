package main

import (
	"fmt"
	"phi/parser"
	"phi/process"
)

// const program = ` /* SND rule */
// 		prc[pid1]: send self<pid3, self>
// 		prc[pid2]: <a, b> <- recv pid1; close self
// 		`

const program = ` /* RCV rule */
		prc[pid1]: send pid2<pid3, self>
		prc[pid2]: <a, b> <- recv self; close self
	`

func main() {
	// Execute from file
	// processes, err := parser.ParseFile("cmd/examples/ex1.txt")

	// Execute directly from string
	processes, err := parser.ParseString(program)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	process.InitializeProcesses(processes, nil)

	// Run via API
	setupAPI()
}

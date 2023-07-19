package main

import (
	"phi/parser"
	"phi/process"
)

//	func main() {
//		fmt.Println("ok")
//	}

func main() {
	// parser.Check()

	processes := parser.ParseFile("parser/input.test")

	process.InitializeProcesses(processes)
}

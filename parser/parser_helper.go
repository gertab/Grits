package parser

import (
	"fmt"
	"os"
	"phi/process"
)

type unexpandedProcesses struct {
	procs     []earlyProcess
	functions []process.FunctionDefinition
}

// Process that is currently being parses and yet to become a process.Process
type earlyProcess struct {
	Body                process.Form
	Names               []process.Name
	FunctionDefinitions *[]process.FunctionDefinition
}

func expandUnexpandedProcesses(u unexpandedProcesses) []process.Process {

	var processes []process.Process

	for _, p := range u.procs {
		for _, n := range p.Names {
			processes = append(processes, process.Process{Body: p.Body, FunctionDefinitions: p.FunctionDefinitions, Channel: n})
		}
	}

	return processes
}

func ParseFile(fileName string) []process.Process {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	// LexAndPrintTokens(file)
	prc, err := Parse(file)

	switch {
	case err != nil:
		fmt.Println("Parsing error: ")
		fmt.Println(err)
		panic("Parsing error!")
	default:
		return expandUnexpandedProcesses(prc)
	}
}

func Check() {
	processes := ParseFile("parser/input.test")

	for _, p := range processes {
		fmt.Println(p.Body.String())
		if p.FunctionDefinitions != nil {
			fmt.Println(len(*p.FunctionDefinitions))
		}
	}
}

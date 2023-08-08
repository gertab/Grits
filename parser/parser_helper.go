package parser

import (
	"fmt"
	"os"
	"phi/process"
	"strings"
)

type unexpandedProcesses struct {
	procs     []incompleteProcess
	functions []process.FunctionDefinition
}

// Process that is currently being parsed and yet to become a process.Process
type incompleteProcess struct {
	Body  Form
	Names []process.Name
	// FunctionDefinitions *[]process.FunctionDefinition
}

type Form interface {
	String() string
	FreeNames() []process.Name
	Substitute(process.Name, process.Name)
	Transition(*process.Process, *process.RuntimeEnvironment)
}

func expandUnexpandedProcesses(u unexpandedProcesses) []process.Process {

	processes := make([]process.Process, len(u.procs))

	counter := 0
	for _, p := range u.procs {
		for _, n := range p.Names {
			new_p := process.NewProcess(p.Body, n, process.LINEAR, &u.functions)
			processes[counter] = *new_p
			counter++
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
		// fmt.Println("Parsing error: ")
		fmt.Println(err)
		panic("Parsing error!")
	default:
		expandedProcesses := expandUnexpandedProcesses(prc)
		return expandedProcesses
	}
}

func ParseString(program string) []process.Process {
	r := strings.NewReader(program)

	prc, err := Parse(r)

	switch {
	case err != nil:
		// fmt.Println("Parsing error: ")
		fmt.Println(err)
		panic("Parsing error!")
	default:
		expandedProcesses := expandUnexpandedProcesses(prc)
		return expandedProcesses
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

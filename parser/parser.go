package parser

import (
	"fmt"
	"os"
	"phi/process"
	"phi/types"
	"strings"
)

type allEnvironment struct {
	procsAndFuns []unexpandedProcessOrFunction
}

type unexpandedProcessOrFunction struct {
	kind         Kind
	proc         incompleteProcess
	function     process.FunctionDefinition
	session_type types.SessionTypeDefinition
}

type Kind int

const (
	PROCESS Kind = iota
	FUNCTION_DEF
	TYPE_DEF
)

// Process that is currently being parsed and yet to become a process.Process
type incompleteProcess struct {
	Body      process.Form
	Providers []process.Name
	Type      types.SessionType
}

func ParseString(program string) ([]*process.Process, *process.GlobalEnvironment, error) {
	r := strings.NewReader(program)

	allEnvironment, err := Parse(r)

	if err != nil {
		fmt.Println(err)
		// panic("Parsing error!")
		return nil, nil, err
	}

	expandedProcesses, globalEnv := expandProcesses(allEnvironment)
	// polarizedProcesses := polarizeProcesses(expandedProcesses)
	// finalizedProcesses, globalEnv := finalizeProcesses(expandedProcesses, globalEnv)

	return expandedProcesses, globalEnv, nil
}

func ParseFile(fileName string) ([]*process.Process, *process.GlobalEnvironment, error) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	// LexAndPrintTokens(file)
	allEnvironment, err := Parse(file)

	if err != nil {
		fmt.Println(err)
		// panic("Parsing error!")
		return nil, nil, err
	}

	expandedProcesses, globalEnv := expandProcesses(allEnvironment)
	// polarizedProcesses := polarizeProcesses(expandedProcesses)
	// finalizedProcesses, globalEnv := finalizeProcesses(expandedProcesses, globalEnv)

	return expandedProcesses, globalEnv, nil
}

// func Check() {
// 	processes, _ := ParseFile("parser/input.test")

// 	for _, p := range processes {
// 		fmt.Println(p.Body.String())
// 		if p.FunctionDefinitions != nil {
// 			fmt.Println(len(*p.FunctionDefinitions))
// 		}
// 	}
// }

func expandProcesses(u allEnvironment) ([]*process.Process, *process.GlobalEnvironment) {

	var processes []*process.Process
	var functions []process.FunctionDefinition
	var types []types.SessionTypeDefinition

	// Collect all functions and types
	for _, p := range u.procsAndFuns {
		if p.kind == FUNCTION_DEF {
			functions = append(functions, p.function)
		} else if p.kind == TYPE_DEF {
			types = append(types, p.session_type)
		} else if p.kind == PROCESS {
			// Processes may have multiple names:
			// 		e.g. prc[a, b, c, d]: send self<...>
			// This remains one process with multiple providers

			// Package all processes
			new_p := process.NewProcess(p.proc.Body, p.proc.Providers, p.proc.Type, process.LINEAR)
			processes = append(processes, new_p)
		}
	}

	return processes, &process.GlobalEnvironment{FunctionDefinitions: &functions, Types: &types}
}

// func polarizeProcesses(processes []*process.Process) []*process.Process {

// 	// Set stop
// 	tries := 0
// 	limit := len(processes) * len(processes)

// 	// Queue all processes
// 	queue := make([]*process.Process, len(processes))
// 	copy(queue, processes)

// 	// List of known channel polarities
// 	// globalChannels := make(map[process.Name]process.Polarity)

// 	for len(queue) > 0 && tries < limit {
// 		tries += 1

// 		// Pick next process and discard top element of queue
// 		currentProcess := queue[0]
// 		queue = queue[1:]

// 		// Is empty ?
// 		if len(queue) == 0 {
// 			fmt.Println("Queue is empty !")
// 		}

// 		fmt.Println(currentProcess.String())

// 		// currentProcess.Body.Polarity()

// 		// var processes []process.Process
// 	}

// 	if len(queue) > 0 {
// 		fmt.Println("Error. Did not manage to find all polarities:")

// 		for i, p := range queue {
// 			fmt.Println(i, p.String())
// 		}
// 	}

// 	return processes
// }

// // Convert from a list of pointer of processes to a plain list of processes (ready to be executed)
// func finalizeProcesses(processes []*process.Process, globalEnv *process.GlobalEnvironment) ([]process.Process, *process.GlobalEnvironment) {
// 	result := make([]process.Process, 0)
// 	for _, j := range processes {
// 		result = append(result, *j)
// 	}

// 	return result, globalEnv
// }

// // Forms used as shorthand notations
// // When expanded, these are converted to the other ones
// todo: need to add polarity

// // Send: send to_c<payload_c, continuation_c>; continuation_e
// type SendMacroForm struct {
// 	to_c           process.Name
// 	payload_c      process.Name
// 	continuation_c process.Name
// 	// Extra used for shorthand notation
// 	continuation_e process.Form
// }

// func NewSendMacroForm(to_c, payload_c, continuation_c process.Name, continuation_e process.Form) *SendMacroForm {
// 	return &SendMacroForm{
// 		to_c:           to_c,
// 		payload_c:      payload_c,
// 		continuation_c: continuation_c,
// 		continuation_e: continuation_e}
// }

// func (p *SendMacroForm) String() string {
// 	var buf bytes.Buffer
// 	buf.WriteString("send ")
// 	buf.WriteString(p.to_c.String())
// 	buf.WriteString("<")
// 	buf.WriteString(p.payload_c.String())
// 	buf.WriteString(",")
// 	buf.WriteString(p.continuation_c.String())
// 	buf.WriteString(">; ")
// 	buf.WriteString(p.continuation_e.String())
// 	return buf.String()
// }

// func (p *SendMacroForm) Substitute(old, new process.Name) {
// }

// // Free names, excluding self references
// func (p *SendMacroForm) FreeNames() []process.Name {
// 	var fn []process.Name
// 	return fn
// }

// func (f *SendMacroForm) Transition(process *process.Process, re *process.RuntimeEnvironment) {
// 	// Should never be called
// 	panic("Unexpanded form found")
// }

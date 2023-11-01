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
	Body      process.Form
	Providers []process.Name
	// FunctionDefinitions *[]process.FunctionDefinition
}

func expandProcesses(u unexpandedProcesses) []*process.Process {

	var processes []*process.Process

	// First step is to duplicate process having multiple names:
	// 		e.g. prc[a, b, c, d]: send self<...>
	// becomes 4 separate processes

	// This is wrong because you can have a free name pointing to another process that becomes duplicated (which shouldn't unless shareable)
	// Todo change this to a split construct:
	//   <a, b, c, d> <- split ...

	// todo maybe throw list of names in OtherProviders
	for _, p := range u.procs {
		// for _, n := range p.Providers {
		new_p := process.NewProcess(p.Body, p.Providers, process.LINEAR, &u.functions)
		processes = append(processes, new_p)
		// }
	}

	// The next step is to get rid of all the macros
	// this is not easy since some macros may need type information
	// todo
	// for _, p := range processes {
	// 	p.Body.
	// }

	return processes
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

// Convert from a list of pointer of processes to a plain list of processes (ready to be executed)
func finalizeProcesses(processes []*process.Process) []process.Process {
	result := make([]process.Process, 0)
	for _, j := range processes {
		result = append(result, *j)
	}

	return result
}

func ParseFile(fileName string) ([]process.Process, error) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	// LexAndPrintTokens(file)
	prc, err := Parse(file)

	if err != nil {
		fmt.Println(err)
		// panic("Parsing error!")
		return nil, err
	}

	expandedProcesses := expandProcesses(prc)
	// polarizedProcesses := polarizeProcesses(expandedProcesses)
	finalizedProcesses := finalizeProcesses(expandedProcesses)

	return finalizedProcesses, nil
}

func ParseString(program string) ([]process.Process, error) {
	r := strings.NewReader(program)

	prc, err := Parse(r)

	if err != nil {
		fmt.Println(err)
		// panic("Parsing error!")
		return nil, err
	}

	expandedProcesses := expandProcesses(prc)
	// polarizedProcesses := polarizeProcesses(expandedProcesses)
	finalizedProcesses := finalizeProcesses(expandedProcesses)

	return finalizedProcesses, nil
}

func Check() {
	processes, _ := ParseFile("parser/input.test")

	for _, p := range processes {
		fmt.Println(p.Body.String())
		if p.FunctionDefinitions != nil {
			fmt.Println(len(*p.FunctionDefinitions))
		}
	}
}

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

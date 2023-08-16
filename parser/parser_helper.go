package parser

import (
	"bytes"
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

func expandProcesses(u unexpandedProcesses) []process.Process {

	var processes []process.Process

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
		processes = append(processes, *new_p)
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

func ParseFile(fileName string) []process.Process {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	// LexAndPrintTokens(file)
	prc, err := Parse(file)

	if err != nil {
		fmt.Println(err)
		panic("Parsing error!")
	}

	expandedProcesses := expandProcesses(prc)
	return expandedProcesses
}

func ParseString(program string) []process.Process {
	r := strings.NewReader(program)

	prc, err := Parse(r)

	if err != nil {
		fmt.Println(err)
		panic("Parsing error!")
	}

	expandedProcesses := expandProcesses(prc)
	return expandedProcesses
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

// Forms used as shorthand notations
// When expanded, these are converted to the other ones

// Send: send to_c<payload_c, continuation_c>; continuation_e
type SendMacroForm struct {
	to_c           process.Name
	payload_c      process.Name
	continuation_c process.Name
	// Extra used for shorthand notation
	continuation_e process.Form
}

func NewSendMacroForm(to_c, payload_c, continuation_c process.Name, continuation_e process.Form) *SendMacroForm {
	return &SendMacroForm{
		to_c:           to_c,
		payload_c:      payload_c,
		continuation_c: continuation_c,
		continuation_e: continuation_e}
}

func (p *SendMacroForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("send ")
	buf.WriteString(p.to_c.String())
	buf.WriteString("<")
	buf.WriteString(p.payload_c.String())
	buf.WriteString(",")
	buf.WriteString(p.continuation_c.String())
	buf.WriteString(">; ")
	buf.WriteString(p.continuation_e.String())
	return buf.String()
}

func (p *SendMacroForm) Substitute(old, new process.Name) {
}

// Free names, excluding self references
func (p *SendMacroForm) FreeNames() []process.Name {
	var fn []process.Name
	return fn
}

func (f *SendMacroForm) Transition(process *process.Process, re *process.RuntimeEnvironment) {
	// Should never be called
	panic("Unexpanded form found")
}

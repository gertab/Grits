package parser

import (
	"fmt"
	"os"
	"phi/position"
	"phi/process"
	"phi/types"
	"strings"
)

type allEnvironment struct {
	procsAndFuns []unexpandedProcessOrFunction
}

type unexpandedProcessOrFunction struct {
	kind                 Kind
	proc                 incompleteProcess
	function             process.FunctionDefinition
	session_type         types.SessionTypeDefinition
	assumedFreeNameTypes []process.Name
	position             position.Position
}

// Different kinds of statements
type Kind int

const (
	PROCESS_DEF Kind = iota
	FUNCTION_DEF
	TYPE_DEF
	ASSUMING_DEF
	EXEC_DEF
)

// Process that is currently being parsed and yet to become a process.Process
type incompleteProcess struct {
	Body      process.Form
	Providers []process.Name
	Type      types.SessionType
}

func ParseString(program string) ([]*process.Process, []process.Name, *process.GlobalEnvironment, error) {
	r := strings.NewReader(program)

	allEnvironment, err := Parse(r)

	if err != nil {
		return nil, nil, nil, err
	}

	expandedProcesses, assumedFreeNames, globalEnv, err := expandProcesses(allEnvironment)

	if err != nil {
		return nil, nil, nil, err
	}

	return expandedProcesses, assumedFreeNames, globalEnv, nil
}

func ParseFile(fileName string) ([]*process.Process, []process.Name, *process.GlobalEnvironment, error) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	// LexAndPrintTokens(file)
	allEnvironment, err := Parse(file)

	if err != nil {
		return nil, nil, nil, err
	}

	expandedProcesses, assumedFreeNames, globalEnv, err := expandProcesses(allEnvironment)

	if err != nil {
		return nil, nil, nil, err
	}

	return expandedProcesses, assumedFreeNames, globalEnv, nil
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

func expandProcesses(u allEnvironment) ([]*process.Process, []process.Name, *process.GlobalEnvironment, error) {

	var processes []*process.Process
	var assumedFreeNames []process.Name
	var functions []process.FunctionDefinition
	var types []types.SessionTypeDefinition

	// Collect all functions and types
	for _, p := range u.procsAndFuns {
		if p.kind == FUNCTION_DEF {

			if p.function.UsesExplicitProvider {
				// Substitute any reference to the explicit provider, with the new version which contains IsSelf = true
				p.function.Body.Substitute(p.function.ExplicitProvider, p.function.ExplicitProvider)
			}

			// Set position
			p.function.Position = p.position

			functions = append(functions, p.function)
		} else if p.kind == TYPE_DEF {
			// Set position
			p.session_type.Position = p.position

			types = append(types, p.session_type)
		} else if p.kind == PROCESS_DEF {
			// Processes may have multiple provider names:
			// 		e.g. prc[a, b, c, d]: send self<...>

			// Define process
			new_p := process.NewProcess(p.proc.Body, p.proc.Providers, p.proc.Type, process.LINEAR, p.position)

			if len(new_p.Providers) == 1 {
				// Set IsSelf to true for the explicit provider
				new_p.Body.Substitute(new_p.Providers[0], new_p.Providers[0])
			} else if len(new_p.Providers) > 1 {
				fn := new_p.Body.FreeNames()
				for _, j := range new_p.Providers {
					if j.ContainedIn(fn) {
						// Since there are multiple names for 'self', then only 'self' can be used
						return nil, nil, nil, fmt.Errorf("name '%s' cannot be referenced directly in %s. Use 'self' instead", j.String(), new_p.Body.String())
					}
				}
			}

			// Package all processes along with the types of the free names
			processes = append(processes, new_p)
		} else if p.kind == ASSUMING_DEF {
			// Collect assumed free names from 'assuming ...'
			assumedFreeNames = append(assumedFreeNames, p.assumedFreeNameTypes...)
		}
	}

	// Define process
	execCount := 0
	for _, p := range u.procsAndFuns {
		if p.kind == EXEC_DEF {
			execCount += 1

			// fetch type from the function type

			functionName := p.proc.Body.(*process.CallForm).FunctionName()

			function := process.GetFunctionByNameArity(functions, functionName, 0)
			if function == nil {
				return nil, nil, nil, fmt.Errorf("invalid calling exec on %s()", functionName)
			}
			new_p := process.NewProcess(p.proc.Body, []process.Name{{Ident: fmt.Sprintf("exec%d", execCount), IsSelf: true}}, function.Type, process.LINEAR, p.position)
			processes = append(processes, new_p)
		}
	}

	return processes, assumedFreeNames, &process.GlobalEnvironment{FunctionDefinitions: &functions, Types: &types}, nil
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

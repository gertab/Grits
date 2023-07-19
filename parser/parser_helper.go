package parser

import "phi/process"

type unexpandedProcesses struct {
	procs     []process.Process
	functions []process.FunctionDefinition
}

func expandUnexpandedProcesses(u unexpandedProcesses) []process.Process {
	// for _, p := range u.procs {
	// 	p.InsertFunctionDefinitions(&u.functions)
	// }

	return u.procs
}

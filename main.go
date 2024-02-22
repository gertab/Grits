package main

import (
	"grits/cmd"
	"grits/parser"
	"grits/process"
	"log"
	"time"
)

const development = false

func main() {
	if development {
		p := `
			prc[a] : 1 = wait b; print ok; close self
			prc[b] : 1 = close self
			`
		dev(p)
	} else {
		cmd.Cli()
	}
}

func dev(program string) {
	// For DEVELOPMENT only: we can run programs directly, bypassing the CLI version
	const (
		executionVersion = process.NORMAL_ASYNC
		typecheck        = true
		execute          = true
		delay            = 0 * time.Millisecond
	)
	var processes []*process.Process
	var assumedFreeNames []process.Name
	var globalEnv *process.GlobalEnvironment
	var err error

	processes, assumedFreeNames, globalEnv, err = parser.ParseString(program)

	if err != nil {
		log.Fatal(err)
		return
	}

	// globalEnv.LogLevels = []process.LogLevel{}
	globalEnv.LogLevels = []process.LogLevel{
		process.LOGINFO,
		process.LOGRULE,
		process.LOGPROCESSING,
		process.LOGRULEDETAILS,
		process.LOGMONITOR,
	}

	if typecheck {
		err = process.Typecheck(processes, assumedFreeNames, globalEnv)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	if execute {
		re, _, _ := process.NewRuntimeEnvironment()
		re.GlobalEnvironment = globalEnv
		re.Typechecked = typecheck
		re.Delay = delay

		process.InitializeProcesses(processes, nil, nil, re)
	}
}

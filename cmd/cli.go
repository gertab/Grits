package cmd

import (
	"flag"
	"fmt"
	"log"
	"phi/benchmarks"
	"phi/parser"
	"phi/process"
	"phi/webserver"
	"time"
)

const PHI = "phi"
const executionVersion = process.NORMAL_ASYNC

/*
Usage of ./phi:

	-benchmark
	      run benchmarks for current program
	-benchmarks
	      start all (pre-configured) benchmarks
	-execute
	      execute processes (default true)
	-noexecute
	      do not execute processes (equivalent to -execute=false)
	-typecheck
	      run typechecker (default true)
	-notypecheck
	      skip typechecker (equivalent to -typecheck=false)
	-repeat uint
	      number of repetitions do when benchmarking (default 10)
	-verbosity int
	      verbosity level (1 = least, 3 = most) (default 2)
	-webserver
	      start webserver
	-addr string
	      webserver address (default ":8081")
*/

// Entry point to run via CLI
func Cli() {
	// Execution Flags
	typecheck := flag.Bool("typecheck", true, "run typechecker")
	noTypecheck := flag.Bool("notypecheck", false, "skip typechecker (equivalent to -typecheck=false)")
	execute := flag.Bool("execute", true, "execute processes")
	noExecute := flag.Bool("noexecute", false, "do not execute processes (equivalent to -execute=false)")
	logLevel := flag.Int("verbosity", 2, "verbosity level (1 = least, 3 = most)")

	// Benchmarking flags
	doAllBenchmarks := flag.Bool("benchmarks", false, "start all (pre-configured) benchmarks")
	benchmark := flag.Bool("benchmark", false, "run benchmarks for current program")
	benchmarkRepeatCount := flag.Uint("repeat", 3, "number of repetitions do when benchmarking")

	// Webserver
	startWebserver := flag.Bool("webserver", false, "start webserver")

	// todo: add execute synchronous vs asynchronous and with polarities

	flag.Parse()
	args := flag.Args()

	if *doAllBenchmarks {
		if len(args) >= 1 {
			fmt.Println("Will run pre-configured benchmarks. Avoid including filename")
			return
		}

		// Run benchmarks and terminate
		benchmarks.Benchmarks()
		return
	}

	if *benchmark {
		if len(args) < 1 {
			fmt.Println("expected name of file to benchmark")
			return
		}

		benchmarks.BenchmarkFile(args[0], *benchmarkRepeatCount)
		return
	}

	typecheckRes := !*noTypecheck && *typecheck
	executeRes := !*noExecute && *execute

	if *logLevel < 1 {
		*logLevel = 1
	} else if *logLevel > 3 {
		*logLevel = 3
	}

	fmt.Printf("%v -- typecheck: %v, execute: %v, verbosity: %d, webserver: %v, benchmark: %v\n", PHI, typecheckRes, executeRes, *logLevel, *startWebserver, *benchmark)

	if *startWebserver {
		// Run via API
		webserver.SetupAPI()
		return
	}

	var processes []*process.Process
	var assumedFreeNames []process.Name
	var globalEnv *process.GlobalEnvironment
	var err error

	if len(args) < 1 {
		err := fmt.Errorf("expected name of file to be executed")
		log.Fatal(err)
		return
	}

	processes, assumedFreeNames, globalEnv, err = parser.ParseFile(args[0])

	if err != nil {
		log.Fatal(err)
		return
	}

	globalEnv.LogLevels = generateLogLevel(*logLevel)

	if typecheckRes {
		err = process.Typecheck(processes, assumedFreeNames, globalEnv)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	if executeRes {
		re := &process.RuntimeEnvironment{
			GlobalEnvironment: globalEnv,
			UseMonitor:        false,
			Color:             true,
			ExecutionVersion:  executionVersion,
			Typechecked:       typecheckRes,
			Delay:             0 * time.Millisecond,
			Quiet:             false,
		}

		process.InitializeProcesses(processes, nil, nil, re)
	}
}

// Generate log levels: 1 = least verbose, 3 = most verbose
// todo maybe add level 0 for quiet
func generateLogLevel(logLevel int) []process.LogLevel {
	if logLevel < 1 {
		logLevel = 1
	}

	switch logLevel {
	case 1:
		return []process.LogLevel{
			process.LOGINFO,
		}
	case 2:
		return []process.LogLevel{
			process.LOGINFO,
			process.LOGRULE,
		}
	default:
		return []process.LogLevel{
			process.LOGINFO,
			process.LOGRULE,
			process.LOGPROCESSING,
			process.LOGRULEDETAILS,
			process.LOGMONITOR,
		}
	}
}

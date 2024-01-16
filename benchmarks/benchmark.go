package benchmarks

import (
	"fmt"
	"log"
	"phi/parser"
	"phi/process"
	"time"
)

const executionVersion = process.NORMAL_ASYNC

func Benchmark() {
	fmt.Println("Benchmarking...")

	var processes []*process.Process
	var assumedFreeNames []process.Name
	var globalEnv *process.GlobalEnvironment
	var err error

	program := `
	prc[a] : 1 = close self
	prc[b] : 1 = wait a; print aa; close self
	`
	// if development {
	processes, assumedFreeNames, globalEnv, err = parser.ParseString(program)
	// } else {
	// 	processes, assumedFreeNames, globalEnv, err = parser.ParseFile(args[0])
	// }

	if err != nil {
		log.Fatal(err)
		return
	}

	globalEnv.LogLevels = []process.LogLevel{}

	err = process.Typecheck(processes, assumedFreeNames, globalEnv)
	if err != nil {
		log.Fatal(err)
		return
	}

	re, ctx, cancel := process.NewRuntimeEnvironment()
	defer cancel()
	re.GlobalEnvironment = globalEnv
	re.ExecutionVersion = executionVersion

	// Will run the benchmarks
	re.Benchmark = false

	fmt.Printf("Initializing %d processes\n", len(processes))

	channels := re.CreateChannelForEachProcess(processes)

	re.SubstituteNameInitialization(processes, channels)

	const heartbeatDelay = 50 * time.Millisecond

	go re.HeartbeatReceiver(heartbeatDelay, cancel)

	re.StartTransitions(processes)

	select {
	case <-ctx.Done():
	case err := <-re.ErrorChan():
		log.Fatal(err)
	}

	timeTaken := re.TimeTaken()
	fmt.Printf("Finished in %vms (%v) \n", timeTaken.Microseconds(), timeTaken)
}

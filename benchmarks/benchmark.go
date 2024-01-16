package benchmarks

import (
	"fmt"
	"io"
	"log"
	"os"
	"phi/parser"
	"phi/process"
	"strings"
	"time"
)

const program = `
	
let ff() : 1 =
w : 1 <- new close w;
wait w;
close self

prc[a, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12] : 1 =  f1 <- new ff();
		   f2 <- new ff();
		   f3 <- new ff();
		   f4 <- new ff();
		   f5 <- new ff();
		   f6 <- new ff();
		   f7 <- new ff();
		   f8 <- new ff();
		   f9 <- new ff();
		   f10 <- new ff();
		   f11 <- new ff();
		   f12 <- new ff();
		   f13 <- new ff();
		   f14 <- new ff();
		   f15 <- new ff();
		   f16 <- new ff();
		   f17 <- new ff();
		   f18 <- new ff();
		   f19 <- new ff();
		   f20 <- new ff();
		   wait f1; 
		   wait f19;
		   wait f3; 
		   wait f4; 
		   wait f17;
		   wait f6; 
		   wait f7; 
		   wait f15;
		   wait f9; 
		   wait f10;
		   wait f11;
		   wait f8; 
		   wait f12;
		   wait f13;
		   wait f14;
		   wait f5; 
		   wait f16;
		   wait f2; 
		   wait f18;
		   wait f20;
		   print ok;
		   close self
	`

const prog2 = `
type A = &{label : 1}
type B = 1 -* 1

let f(y : A, z : B) : A * B = send self<y, z>

assuming a : A, b : B

prc[pid1] : 1
       = x <- new f(a, b); 
				<u, v> <- recv x;  
				drop u; 
				drop v; 
				close self`

func BenchmarkFile(fileName string, repetitions uint) {
	programFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Couldn't open file: ", err)
		return
	}

	timeTaken, err := runTiming(programFile, process.NORMAL_ASYNC)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Finished in %vµs (%v) \n", timeTaken.Microseconds(), timeTaken)
}

func Benchmark(repetitions uint) {
	fmt.Println("Benchmarking...")

	timeTaken, err := runTiming(strings.NewReader(prog2), process.NORMAL_ASYNC)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Finished in %vµs (%v) \n", timeTaken.Microseconds(), timeTaken)
}

func runTiming(program io.Reader, executionVersion process.Execution_Version) (time.Duration, error) {

	var processes []*process.Process
	var assumedFreeNames []process.Name
	var globalEnv *process.GlobalEnvironment
	var err error

	processes, assumedFreeNames, globalEnv, err = parser.ParseReader(program)

	if err != nil {
		// log.Fatal(err)
		return 0, err
	}

	globalEnv.LogLevels = []process.LogLevel{}
	// globalEnv.LogLevels = []process.LogLevel{process.LOGINFO, process.LOGPROCESSING}

	err = process.Typecheck(processes, assumedFreeNames, globalEnv)
	if err != nil {
		// log.Fatal(err)
		return 0, err
	}

	re, ctx, cancel := process.NewRuntimeEnvironment()
	defer cancel()
	re.GlobalEnvironment = globalEnv
	re.ExecutionVersion = executionVersion
	re.Typechecked = true

	// Will run the benchmarks
	re.Benchmark = false

	// fmt.Printf("Initializing %d processes\n", len(processes))

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

	return re.TimeTaken(), nil
}

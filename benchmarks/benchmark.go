package benchmarks

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"phi/parser"
	"phi/process"
	"time"
)

const programExample = `
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

// Runs benchmark for one file
const detailed = false

func BenchmarkFile(fileName string, repetitions uint) {
	// runtime.GOMAXPROCS(1)
	fileNameBase := filepath.Base(fileName)

	programFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Couldn't open file: ", err)
		return
	}

	programFileBytes, _ := io.ReadAll(programFile)

	fmt.Printf("Running benchmark for %s\n", fileNameBase)

	// timeTaken, processCount, err := runTiming(bytes.NewReader(programFileBytes), process.NORMAL_ASYNC)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("Finished in %vµs (%v) -- %v processes \n", timeTaken.Microseconds(), timeTaken, processCount)

	// Run all timings repeatedly
	var allResults []TimingResult

	for i := 0; i < int(repetitions); i++ {
		result := runAllTimingsOnce(bytes.NewReader(programFileBytes))

		if result != nil {
			// fmt.Println(result)
			allResults = append(allResults, *result)
		}
	}

	if detailed {
		fmt.Println("Obtained", len(allResults), "results:")
		fmt.Println(csvHeader())
		for _, row := range allResults {
			fmt.Println(row.csvRow())
		}
	}

	// err = saveToFileCSV(fileNameBase, allResults)

	// if err != nil {
	// 	fmt.Println("Could save to file", err)
	// }

	average := getAverage(allResults)
	average.name = fileNameBase
	fmt.Println(average)
	// fmt.Println(average.csvRow())

	err = saveToFileCSV(fileNameBase, []TimingResult{*average})

	if err != nil {
		fmt.Println("Could save to file", err)
	}
}

// Runs pre-configured benchmarks
func Benchmarks(repetitions uint) {
	fmt.Println("Benchmarking...")

	BenchmarkFile("./benchmarks/compare/nat-double/nat-double-13.phi", repetitions)

	// timeTaken, processCount, err := runTiming(strings.NewReader(programExample), process.NORMAL_ASYNC)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("Finished in %vµs (%v) -- %v processes \n", timeTaken.Microseconds(), timeTaken, processCount)
}

// Run the same program using all transition variations
func runAllTimingsOnce(program io.Reader) *TimingResult {
	programFileBytes, _ := io.ReadAll(program)

	var result TimingResult

	// Version 1:
	timeTaken, count, err := runTiming(bytes.NewReader(programFileBytes), process.NON_POLARIZED_SYNC)
	if err != nil {
		return nil
	}

	result.timeNonPolarizedSync = timeTaken
	result.processCountNonPolarizedSync = count

	// Version 2 (Async):
	timeTaken2, count2, err := runTiming(bytes.NewReader(programFileBytes), process.NORMAL_ASYNC)
	if err != nil {
		return nil
	}

	result.timeNormalAsync = timeTaken2
	result.processCountNormalAsync = count2

	// Version 2 (Sync):
	timeTaken3, count3, err := runTiming(bytes.NewReader(programFileBytes), process.NORMAL_SYNC)
	if err != nil {
		return nil
	}

	result.timeNormalSync = timeTaken3
	result.processCountNormalSync = count3

	return &result
}

type TimingResult struct {
	name                         string
	timeNonPolarizedSync         time.Duration
	processCountNonPolarizedSync uint64
	timeNormalAsync              time.Duration
	processCountNormalAsync      uint64
	timeNormalSync               time.Duration
	processCountNormalSync       uint64
}

func (t *TimingResult) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("File: %v\n", t.name))
	buffer.WriteString(fmt.Sprintf("\tv1: \t\t%vµs (%v) -- %d processes\n", t.timeNonPolarizedSync.Microseconds(), t.timeNonPolarizedSync, t.processCountNonPolarizedSync))
	buffer.WriteString(fmt.Sprintf("\tv2(async):\t%vµs (%v) -- %d processes\n", t.timeNormalAsync.Microseconds(), t.timeNormalAsync, t.processCountNormalAsync))
	buffer.WriteString(fmt.Sprintf("\tv2(sync):\t%vµs (%v) -- %d processes\n", t.timeNormalSync.Microseconds(), t.timeNormalSync, t.processCountNormalSync))

	return buffer.String()
}

func (t *TimingResult) csvRow() string {
	var buffer bytes.Buffer

	buffer.WriteString(t.name)
	buffer.WriteString(separator)
	buffer.WriteString(fmt.Sprintf("%v", t.timeNonPolarizedSync.Microseconds()))
	buffer.WriteString(separator)
	buffer.WriteString(fmt.Sprintf("%d", t.processCountNonPolarizedSync))
	buffer.WriteString(separator)
	buffer.WriteString(fmt.Sprintf("%v", t.timeNormalAsync.Microseconds()))
	buffer.WriteString(separator)
	buffer.WriteString(fmt.Sprintf("%d", t.processCountNormalAsync))
	buffer.WriteString(separator)
	buffer.WriteString(fmt.Sprintf("%v", t.timeNormalSync.Microseconds()))
	buffer.WriteString(separator)
	buffer.WriteString(fmt.Sprintf("%d", t.processCountNormalSync))

	return buffer.String()
}

const separator = ","

func csvHeader() string {
	var buffer bytes.Buffer

	buffer.WriteString("name")
	buffer.WriteString(separator)
	buffer.WriteString("timeNonPolarizedSync")
	buffer.WriteString(separator)
	buffer.WriteString("processCountNonPolarizedSync")
	buffer.WriteString(separator)
	buffer.WriteString("timeNormalAsync")
	buffer.WriteString(separator)
	buffer.WriteString("processCountNormalAsync")
	buffer.WriteString(separator)
	buffer.WriteString("timeNormalSync")
	buffer.WriteString(separator)
	buffer.WriteString("processCountNormalSync")

	return buffer.String()
}

func saveToFileCSV(fileName string, results []TimingResult) error {
	const extension = ".csv"
	name := fileName + "-benchmark" + extension
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write([]byte(csvHeader()))
	f.Write([]byte("\n"))

	for _, row := range results {
		f.Write([]byte(row.csvRow()))
		f.Write([]byte("\n"))
	}

	return nil
}

func getAverage(allResults []TimingResult) *TimingResult {
	result := allResults[0]

	count := len(allResults)

	if count > 1 {
		for i := 1; i < count; i += 1 {
			result.timeNonPolarizedSync += allResults[i].timeNonPolarizedSync
			result.processCountNonPolarizedSync += allResults[i].processCountNonPolarizedSync
			result.timeNormalAsync += allResults[i].timeNormalAsync
			result.processCountNormalAsync += allResults[i].processCountNormalAsync
			result.timeNormalSync += allResults[i].timeNormalSync
			result.processCountNormalSync += allResults[i].processCountNormalSync
		}

		// Get average
		result.timeNonPolarizedSync /= time.Duration(count)
		result.processCountNonPolarizedSync /= uint64(count)
		result.timeNormalAsync /= time.Duration(count)
		result.processCountNormalAsync /= uint64(count)
		result.timeNormalSync /= time.Duration(count)
		result.processCountNormalSync /= uint64(count)
	}

	return &result
}

func runTiming(program io.Reader, executionVersion process.Execution_Version) (time.Duration, uint64, error) {

	var processes []*process.Process
	var assumedFreeNames []process.Name
	var globalEnv *process.GlobalEnvironment
	var err error

	processes, assumedFreeNames, globalEnv, err = parser.ParseReader(program)

	if err != nil {
		// log.Fatal(err)
		return 0, 0, err
	}

	err = process.Typecheck(processes, assumedFreeNames, globalEnv)
	if err != nil {
		// log.Fatal(err)
		return 0, 0, err
	}

	re, ctx, cancel := process.NewRuntimeEnvironment()
	defer cancel()
	re.GlobalEnvironment = globalEnv
	re.ExecutionVersion = executionVersion
	re.Typechecked = true

	// Suppress print and log outputs
	re.Quiet = false
	globalEnv.LogLevels = []process.LogLevel{}
	// globalEnv.LogLevels = []process.LogLevel{process.LOGINFO, process.LOGPROCESSING}

	// fmt.Printf("Initializing %d processes\n", len(processes))

	channels := re.CreateChannelForEachProcess(processes)

	re.SubstituteNameInitialization(processes, channels)

	const heartbeatDelay = 500 * time.Millisecond
	go re.HeartbeatReceiver(heartbeatDelay, cancel)

	re.StartTransitions(processes)

	select {
	case <-ctx.Done():
	case err := <-re.ErrorChan():
		log.Fatal(err)
	}

	return re.TimeTaken(), re.ProcessCount(), nil
}

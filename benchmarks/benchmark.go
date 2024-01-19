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
	"runtime"
	"sync"
	"time"
)

// All unmarked time units are in Microseconds

// The Benchmarks/Benchmark functions output the benchmark results into a CSV file containing the following columns:
//
//   - name                         : name of file being checked
//   - timeNonPolarizedSync	        : time taken to evaluate file (using v1)
//   - processCountNonPolarizedSync : number of processes spawn (when using v1)
//   - timeNormalAsync              : time taken to evaluate file (using v2-async)
//   - processCountNormalAsync      : number of processes spawn (when using v2-async)
//   - timeNormalSync               : time taken to evaluate file (using v2-sync)
//   - processCountNormalSync       : number of processes spawn (when using v2-sync)
const (
	detailedOutput = true
	// GoMaxProcs          = 4
	outputFileExtension = ".csv"
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
func BenchmarkFile(fileName string, repetitions uint, maxCores int) {
	runtime.GOMAXPROCS(maxCores)

	fileNameBase := filepath.Base(fileName)

	programFileBytes, err := readFile(fileName)

	if err != nil {
		fmt.Println("Couldn't read file: ", err)
		return
	}

	fmt.Printf("Running benchmarks for %s (using %d cores out of %v). ", fileNameBase, maxCores, runtime.NumCPU())

	if repetitions > 1 {
		fmt.Printf("Repeating runs for %d times\n", repetitions)
	} else {
		fmt.Printf("\n")
	}

	// timeTaken, processCount, err := runTiming(bytes.NewReader(programFileBytes), process.NORMAL_ASYNC)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("Finished in %vµs (%v) -- %v processes \n", timeTaken.Microseconds(), timeTaken, processCount)

	// Run all timings repeatedly
	var allResults []TimingResult

	for i := 0; i < int(repetitions); i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		var result TimingResult
		result.invalid = false
		go runAllTimingsOnce(bytes.NewReader(programFileBytes), &wg, &result)
		wg.Wait()

		// fmt.Print(".")

		if !result.invalid {
			if detailedOutput {
				fmt.Println(i+1, result.StringShort())
			}
			allResults = append(allResults, result)
		}
	}

	// fmt.Println()

	// if detailedOutput {
	// 	fmt.Println("Obtained", len(allResults), "results:")
	// 	fmt.Println(csvHeader())
	// 	for _, row := range allResults {
	// 		fmt.Println(row.csvRow())
	// 		fmt.Println(row)
	// 	}
	// }

	// err = saveToFileCSV(fileNameBase, allResults)

	// if err != nil {
	// 	fmt.Println("Could save to file", err)
	// }

	average := getAverage(allResults)
	average.name = fileNameBase
	// fmt.Println(average.csvRow())

	if len(allResults) > 1 {
		fmt.Println("Average times:")
		fmt.Println(average)
	}

	fileName, err = saveToFileCSV(fileNameBase, "benchmark", []TimingResult{*average}, maxCores)
	if err != nil {
		fmt.Println("Could save to file", err)
		return
	}

	fileNameDetailed, err := saveToFileCSV(fileNameBase, "benchmark-detailed", allResults, maxCores)
	if err != nil {
		fmt.Println("Could save to file", err)
		return
	}

	fmt.Printf("Saved results in %v (and detailed version in %s)\n", fileName, fileNameDetailed)
}

// Runs pre-configured benchmarks
// todo remove reps
func Benchmarks(maxCores int) {
	runtime.GOMAXPROCS(maxCores)
	fmt.Printf("Benchmarking... (using %d cores out of %v)\n\n", maxCores, runtime.NumCPU())

	benchmarkCases := []benchmarkCase{
		{"./benchmarks/compare/nat-double/nat-double-1.phi", 25},
		{"./benchmarks/compare/nat-double/nat-double-2.phi", 25},
		{"./benchmarks/compare/nat-double/nat-double-3.phi", 25},
		{"./benchmarks/compare/nat-double/nat-double-4.phi", 25},
		{"./benchmarks/compare/nat-double/nat-double-5.phi", 20},
		{"./benchmarks/compare/nat-double/nat-double-6.phi", 20},
		{"./benchmarks/compare/nat-double/nat-double-7.phi", 20},
		{"./benchmarks/compare/nat-double/nat-double-8.phi", 20},
		{"./benchmarks/compare/nat-double/nat-double-9.phi", 20},
		{"./benchmarks/compare/nat-double/nat-double-10.phi", 15},
		{"./benchmarks/compare/nat-double/nat-double-11.phi", 15},
		{"./benchmarks/compare/nat-double/nat-double-12.phi", 15},
		{"./benchmarks/compare/nat-double/nat-double-13.phi", 15},
		{"./benchmarks/compare/nat-double/nat-double-14.phi", 15},
		{"./benchmarks/compare/nat-double/nat-double-15.phi", 10},
		{"./benchmarks/compare/nat-double/nat-double-16.phi", 8},
	}

	runGroupedBenchmarks(benchmarkCases, "nat", maxCores)
}

func runGroupedBenchmarks(benchmarkCases []benchmarkCase, name string, maxCores int) {
	// Start writing result to file
	benchmarksFilename := name + "-benchmarks-" + fmt.Sprint(maxCores) + outputFileExtension
	f, err := os.Create(benchmarksFilename)
	if err != nil {
		fmt.Println("Couldn't open file: ", err)
		return
	}
	defer f.Close()
	f.WriteString(csvHeader() + "\n")

	var benchmarkCaseResults []benchmarkCaseResult

	for _, file := range benchmarkCases {
		repeat := ""
		if file.repetitions > 1 {
			repeat = fmt.Sprintf("(%d repetitions)", file.repetitions)
		}
		fmt.Printf("Benchmarking %s %s", file.baseName(), repeat)

		// Prepare result
		currentBenchmarkCaseResult := NewBenchmarkCaseResult(file.fileName)

		// Open file
		programFileBytes, err := readFile(file.fileName)

		if err != nil {
			fmt.Println("\nCouldn't read file: ", err)
			currentBenchmarkCaseResult.ok = false
			benchmarkCaseResults = append(benchmarkCaseResults, *currentBenchmarkCaseResult)
			continue
		}

		// Run all timings repeatedly
		var allTimingResults []TimingResult

		for i := 0; i < int(file.repetitions); i++ {
			var wg sync.WaitGroup
			wg.Add(1)
			var result TimingResult
			// Timing are obtained from heres
			go runAllTimingsOnce(bytes.NewReader(programFileBytes), &wg, &result)
			wg.Wait()
			fmt.Print(".")

			if !result.invalid {
				// fmt.Prixntln(result)
				result.name = file.baseName()
				allTimingResults = append(allTimingResults, result)
			}
		}

		currentBenchmarkCaseResult.results = allTimingResults
		currentBenchmarkCaseResult.repetitionsDone = uint(len(allTimingResults))

		fmt.Printf("\n%s\n", currentBenchmarkCaseResult)

		average := getAverage(currentBenchmarkCaseResult.results)
		f.WriteString(average.csvRow() + "\n")

		benchmarkCaseResults = append(benchmarkCaseResults, *currentBenchmarkCaseResult)
	}

	// for i, res := range benchmarkCaseResults {
	// 	fmt.Println(i, res.String())
	// }

	fileName, err := saveDetailedBenchmarks(name, benchmarkCaseResults, maxCores)
	if err != nil {
		fmt.Println("Could save to file", err)
		return
	}

	fmt.Printf("Saved results in %v\n", fileName)
}

// Saves all individual runs to a file
func saveDetailedBenchmarks(name string, benchmarkCaseResults []benchmarkCaseResult, maxCores int) (string, error) {
	fileName := name + "-detailed-benchmarks-" + fmt.Sprint(maxCores) + outputFileExtension
	f, err := os.Create(fileName)
	if err != nil {
		return fileName, err
	}
	defer f.Close()

	f.WriteString(csvHeader() + "\n")

	for _, benchmarkCaseResult := range benchmarkCaseResults {
		if benchmarkCaseResult.ok {
			for _, row := range benchmarkCaseResult.results {
				f.WriteString(row.csvRow() + "\n")
			}
		} else {
			// Empty row for erroneous ones
			f.WriteString(benchmarkCaseResult.fileName + ",,,,,,\n")
		}
	}

	return fileName, nil
}

// Input for each case
type benchmarkCase struct {
	fileName    string // Path of file to run
	repetitions uint
}

func (b *benchmarkCase) baseName() string {
	return filepath.Base(b.fileName)
}

// Output for each case
type benchmarkCaseResult struct {
	fileName        string
	ok              bool
	repetitionsDone uint
	results         []TimingResult
}

func (b *benchmarkCaseResult) baseName() string {
	return filepath.Base(b.fileName)
}

func (b *benchmarkCaseResult) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("File: %v (average of %v)\n", b.baseName(), b.repetitionsDone))

	if !b.ok {
		buffer.WriteString("Couldn't read file")
		return buffer.String()
	}

	average := getAverage(b.results)
	buffer.WriteString(average.StringShort())

	return buffer.String()
}

func NewBenchmarkCaseResult(fileName string) *benchmarkCaseResult {
	// Sets the default settings
	return &benchmarkCaseResult{
		fileName:        fileName,
		ok:              true,
		repetitionsDone: 0,
		results:         nil,
	}
}

type TimingResult struct {
	name                         string
	invalid                      bool
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

func (t *TimingResult) StringShort() string {
	var buffer bytes.Buffer

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

func saveToFileCSV(fileName string, title string, results []TimingResult, maxCores int) (string, error) {
	name := fileName + "-" + title + "-" + fmt.Sprint(maxCores) + outputFileExtension
	f, err := os.Create(name)
	if err != nil {
		return name, err
	}
	defer f.Close()

	f.Write([]byte(csvHeader()))
	f.Write([]byte("\n"))

	for _, row := range results {
		f.Write([]byte(row.csvRow()))
		f.Write([]byte("\n"))
	}

	return name, nil
}

func getAverage(allResults []TimingResult) *TimingResult {
	// allResults are not modified
	count := len(allResults)

	if count == 0 {
		return &TimingResult{}
	}

	result := TimingResult{allResults[0].name, false, allResults[0].timeNonPolarizedSync, allResults[0].processCountNonPolarizedSync, allResults[0].timeNormalAsync, allResults[0].processCountNormalAsync, allResults[0].timeNormalSync, allResults[0].processCountNormalSync}
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

	return &result
}

// Takes a file, reads it and returns the content as an array of bytes
func readFile(fileName string) ([]byte, error) {
	programFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(programFile)
}

// Run the same program using all transition variations
// This should run in a separate goroutine, storing the final results in 'result'
func runAllTimingsOnce(program io.Reader, wg *sync.WaitGroup, result *TimingResult) {
	defer wg.Done()

	programFileBytes, _ := io.ReadAll(program)

	// Version 1:
	timeTaken, count, err := runTiming(bytes.NewReader(programFileBytes), process.NON_POLARIZED_SYNC)
	if err != nil {
		result.invalid = true
		fmt.Println(err)
		return
	}

	result.timeNonPolarizedSync = timeTaken
	result.processCountNonPolarizedSync = count

	// Version 2 (Async):
	timeTaken2, count2, err := runTiming(bytes.NewReader(programFileBytes), process.NORMAL_ASYNC)
	if err != nil {
		result.invalid = true
		return
	}

	result.timeNormalAsync = timeTaken2
	result.processCountNormalAsync = count2

	// Version 2 (Sync):
	timeTaken3, count3, err := runTiming(bytes.NewReader(programFileBytes), process.NORMAL_SYNC)
	if err != nil {
		result.invalid = true
		return
	}

	result.timeNormalSync = timeTaken3
	result.processCountNormalSync = count3
	result.invalid = false
}

// Performs the actual execution
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
	re.Quiet = true
	globalEnv.LogLevels = []process.LogLevel{}
	// globalEnv.LogLevels = []process.LogLevel{process.LOGINFO, process.LOGPROCESSING}

	// fmt.Printf("Initializing %d processes\n", len(processes))

	channels := re.CreateChannelForEachProcess(processes)

	re.SubstituteNameInitialization(processes, channels)

	const heartbeatDelay = 600 * time.Millisecond
	go re.HeartbeatReceiver(heartbeatDelay, cancel)

	re.StartTransitions(processes)

	select {
	case <-ctx.Done():
	case err := <-re.ErrorChan():
		log.Fatal(err)
	}

	timeTaken := re.TimeTaken()
	processCount := re.ProcessCount()

	re = nil
	processes = nil
	assumedFreeNames = nil
	return timeTaken, processCount, nil
}

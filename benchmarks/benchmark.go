package benchmarks

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"grits/parser"
	"grits/process"
	"image/color"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// All unmarked time units are in Microseconds

// The Benchmarks/Benchmark functions output the benchmark results into a CSV file containing the following columns:
//
//   - name                         : name of file being checked
//   - timeSyncV1NP	        : time taken to evaluate file (using v1)
//   - processCountSyncV1NP : number of processes spawn (when using v1)
//   - timeAsyncV2              : time taken to evaluate file (using v2-async)
//   - processCountAsyncV2      : number of processes spawn (when using v2-async)

const (
	detailedOutput      = true
	outputFileExtension = ".csv"
	outputFolder        = "benchmark-results"
)

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

	// Run all timings repeatedly
	var allResults []TimingResult

	for i := 0; i < int(repetitions); i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		var result TimingResult
		result.name = fileNameBase
		result.invalid = false
		result.caseNumber = -1 // ignore case number in individual files
		go runAllTimingsOnce(bytes.NewReader(programFileBytes), &wg, &result)
		wg.Wait()

		// fmt.Print(".")

		if !result.invalid {
			fmt.Println(i+1, result.StringShort())
			allResults = append(allResults, result)
		}
	}

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

	fmt.Printf("Saved results in %v (and detailed version in %s)\n", filepath.Join(outputFolder, fileName), filepath.Join(outputFolder, fileNameDetailed))
}

// Runs pre-configured benchmarks (stored in the folder benchmarks/compare)
func SampleBenchmarks(maxCores int) {
	runtime.GOMAXPROCS(maxCores)

	if err := checkBenchmarksAvailable(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Benchmarking... (using %d cores out of %v)\n\n", maxCores, runtime.NumCPU())

	folder := "nat-double"
	benchmarkCases := []benchmarkCase{
		{"nat-double-1.grits", 1, 4},
		{"nat-double-2.grits", 2, 4},
		{"nat-double-3.grits", 3, 4},
		{"nat-double-4.grits", 4, 4},
		{"nat-double-5.grits", 5, 4},
		{"nat-double-6.grits", 6, 4},
		{"nat-double-7.grits", 7, 4},
		{"nat-double-8.grits", 8, 4},
		{"nat-double-9.grits", 9, 4},
		{"nat-double-10.grits", 10, 4},
		{"nat-double-11.grits", 11, 4},
		{"nat-double-12.grits", 12, 3},
		{"nat-double-13.grits", 13, 3},
		{"nat-double-14.grits", 14, 1},
	}

	resultsFile, err := runGroupedBenchmarks(folder, benchmarkCases, maxCores)
	if err == nil {
		visualisePlots(resultsFile, true)
	}

	folder = "nat-double-parallel"
	benchmarkCases = []benchmarkCase{
		{"nat-double-parallel-2.grits", 2, 4},
		{"nat-double-parallel-8.grits", 8, 4},
		{"nat-double-parallel-14.grits", 14, 4},
		{"nat-double-parallel-20.grits", 20, 4},
		{"nat-double-parallel-26.grits", 26, 4},
		{"nat-double-parallel-32.grits", 32, 4},
		{"nat-double-parallel-38.grits", 38, 4},
		{"nat-double-parallel-44.grits", 44, 4},
		{"nat-double-parallel-50.grits", 50, 4},
	}

	resultsFile, err = runGroupedBenchmarks(folder, benchmarkCases, maxCores)
	if err == nil {
		visualisePlots(resultsFile, false)
	}
}

// Looks for the benchmark/compare folder to make sure that the benchmarking directory is available
func checkBenchmarksAvailable() error {
	msg :=
		`to run benchmarks, you need to be located at the root folder of the Grits source code. The source code and benchmark files can be obtained using the following commands:
> git clone https://github.com/gertab/Grits.git
> cd Grits 
> go build .
> ./grits --sample-benchmarks
> cd benchmark-results
`
	if _, err := os.Stat(filepath.Join("benchmarks", "compare")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(msg)
	}

	return nil
}

// Fetches files from /benchmarks/compare/<folder>/
func runGroupedBenchmarks(folder string, benchmarkCases []benchmarkCase, maxCores int) (string, error) {
	// Start writing result to file
	benchmarksFilename := folder + "-benchmarks-" + fmt.Sprint(maxCores) + outputFileExtension
	fileNameWithFolder := filepath.Join(outputFolder, benchmarksFilename)

	if err := prepareOutputFolder(fileNameWithFolder); err != nil {
		return benchmarksFilename, err
	}

	// Create file
	f, err := os.Create(fileNameWithFolder)
	if err != nil {
		return benchmarksFilename, fmt.Errorf("couldn't open file: %v", err)
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

		// Look for files in benchmarks/compare/
		fullPathName := filepath.Join("benchmarks", "compare", folder, file.fileName)

		// Prepare result
		currentBenchmarkCaseResult := NewBenchmarkCaseResult(fullPathName)
		currentBenchmarkCaseResult.caseNumber = file.caseNumber

		// Open file
		programFileBytes, err := readFile(fullPathName)

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
				result.name = file.baseName()
				result.caseNumber = file.caseNumber
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

	benchmarksDetailed, err := saveDetailedBenchmarks(folder, benchmarkCaseResults, maxCores)
	if err != nil {
		return benchmarksFilename, fmt.Errorf("could not save to file: %v", err)
	}

	fmt.Printf("Saved detailed results in %v\n", benchmarksDetailed)
	return benchmarksFilename, nil
}

// Saves all individual runs to a file
func saveDetailedBenchmarks(name string, benchmarkCaseResults []benchmarkCaseResult, maxCores int) (string, error) {
	fileName := name + "-detailed-benchmarks-" + fmt.Sprint(maxCores) + outputFileExtension
	fileNameWithFolder := filepath.Join(outputFolder, fileName)

	// Prepare folder
	if err := prepareOutputFolder(fileNameWithFolder); err != nil {
		return fileNameWithFolder, err
	}

	// Create file
	f, err := os.Create(fileNameWithFolder)
	if err != nil {
		return fileNameWithFolder, err
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

	return fileNameWithFolder, nil
}

// Input for each case
type benchmarkCase struct {
	fileName    string // Path of file to run
	caseNumber  int
	repetitions uint
}

func (b *benchmarkCase) baseName() string {
	return filepath.Base(b.fileName)
}

// Output for each case
type benchmarkCaseResult struct {
	fileName        string
	ok              bool
	caseNumber      int
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
	name                 string
	invalid              bool
	caseNumber           int
	timeSyncV1NP         time.Duration
	processCountSyncV1NP uint64
	timeAsyncV2          time.Duration
	processCountAsyncV2  uint64
	// timeSyncV2           time.Duration
	// processCountSyncV2   uint64
}

func (t *TimingResult) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("File: %v\n", t.name))
	if t.caseNumber >= 0 {
		buffer.WriteString(fmt.Sprintf(" (%v)\n", t.caseNumber))
	}
	buffer.WriteString(fmt.Sprintf("\tv1 (sync): \t%vµs (%v) -- %d processes\n", t.timeSyncV1NP.Microseconds(), t.timeSyncV1NP, t.processCountSyncV1NP))
	buffer.WriteString(fmt.Sprintf("\tv2 (async):\t%vµs (%v) -- %d processes\n", t.timeAsyncV2.Microseconds(), t.timeAsyncV2, t.processCountAsyncV2))
	// buffer.WriteString(fmt.Sprintf("\tv2 (sync):\t%vµs (%v) -- %d processes\n", t.timeSyncV2.Microseconds(), t.timeSyncV2, t.processCountSyncV2))

	return buffer.String()
}

func (t *TimingResult) StringShort() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("\tv1 (sync): \t%vµs (%v) -- %d processes\n", t.timeSyncV1NP.Microseconds(), t.timeSyncV1NP, t.processCountSyncV1NP))
	buffer.WriteString(fmt.Sprintf("\tv2 (async):\t%vµs (%v) -- %d processes\n", t.timeAsyncV2.Microseconds(), t.timeAsyncV2, t.processCountAsyncV2))
	// buffer.WriteString(fmt.Sprintf("\tv2 (sync):\t%vµs (%v) -- %d processes\n", t.timeSyncV2.Microseconds(), t.timeSyncV2, t.processCountSyncV2))

	return buffer.String()
}

func (t *TimingResult) csvRow() string {
	var buffer bytes.Buffer

	buffer.WriteString(t.name)
	buffer.WriteString(separator)
	if t.caseNumber >= 0 {
		buffer.WriteString(fmt.Sprintf("%v", t.caseNumber))
	}
	buffer.WriteString(separator)
	buffer.WriteString(fmt.Sprintf("%v", t.timeSyncV1NP.Microseconds()))
	buffer.WriteString(separator)
	buffer.WriteString(fmt.Sprintf("%d", t.processCountSyncV1NP))
	buffer.WriteString(separator)
	buffer.WriteString(fmt.Sprintf("%v", t.timeAsyncV2.Microseconds()))
	buffer.WriteString(separator)
	buffer.WriteString(fmt.Sprintf("%d", t.processCountAsyncV2))
	// buffer.WriteString(separator)
	// buffer.WriteString(fmt.Sprintf("%v", t.timeSyncV2.Microseconds()))
	// buffer.WriteString(separator)
	// buffer.WriteString(fmt.Sprintf("%d", t.processCountSyncV2))

	return buffer.String()
}

const separator = ","

func csvHeader() string {
	var buffer bytes.Buffer

	buffer.WriteString("name")
	buffer.WriteString(separator)
	buffer.WriteString("caseNumber")
	buffer.WriteString(separator)
	buffer.WriteString("timeSyncV1NP")
	buffer.WriteString(separator)
	buffer.WriteString("processCountSyncV1NP")
	buffer.WriteString(separator)
	buffer.WriteString("timeAsyncV2")
	buffer.WriteString(separator)
	buffer.WriteString("processCountAsyncV2")
	// buffer.WriteString(separator)
	// buffer.WriteString("timeSyncV2")
	// buffer.WriteString(separator)
	// buffer.WriteString("processCountSyncV2")

	return buffer.String()
}

func saveToFileCSV(fileName string, title string, results []TimingResult, maxCores int) (string, error) {
	name := fileName + "-" + title + "-" + fmt.Sprint(maxCores) + outputFileExtension
	nameWithFolder := filepath.Join(outputFolder, name)

	// Prepare folder
	if err := prepareOutputFolder(nameWithFolder); err != nil {
		return name, err
	}

	// Create file
	f, err := os.Create(nameWithFolder)
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

// Averages a list of results.
// 0 seconds timings are skipped
func getAverage(allResults []TimingResult) *TimingResult {
	if len(allResults) == 0 {
		return &TimingResult{invalid: true}
	}

	result := TimingResult{allResults[0].name, false, allResults[0].caseNumber, 0, 0, 0, 0}

	count := 0
	for _, curResult := range allResults {

		if curResult.timeSyncV1NP == 0 && curResult.timeAsyncV2 == 0 {
			// in case of zero timings, skip record
			continue
		}

		count += 1

		result.timeSyncV1NP += curResult.timeSyncV1NP
		result.processCountSyncV1NP += curResult.processCountSyncV1NP
		result.timeAsyncV2 += curResult.timeAsyncV2
		result.processCountAsyncV2 += curResult.processCountAsyncV2
		// result.timeSyncV2 += curResult.timeSyncV2
		// result.processCountSyncV2 += curResult.processCountSyncV2
	}

	if count == 0 {
		result.invalid = true
		return &result
	}

	// Get average
	result.timeSyncV1NP /= time.Duration(count)
	result.processCountSyncV1NP /= uint64(count)
	result.timeAsyncV2 /= time.Duration(count)
	result.processCountAsyncV2 /= uint64(count)
	// result.timeSyncV2 /= time.Duration(count)
	// result.processCountSyncV2 /= uint64(count)

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

	result.timeSyncV1NP = timeTaken
	result.processCountSyncV1NP = count

	// Version 2 (Async):
	timeTaken2, count2, err := runTiming(bytes.NewReader(programFileBytes), process.NORMAL_ASYNC)
	if err != nil {
		result.invalid = true
		return
	}

	result.timeAsyncV2 = timeTaken2
	result.processCountAsyncV2 = count2

	// Version 2 can also be executes synchronously, however we do not include it in the benchmarks
	// // Version 2 (Sync):
	// timeTaken3, count3, err := runTiming(bytes.NewReader(programFileBytes), process.NORMAL_SYNC)
	// if err != nil {
	// 	result.invalid = true
	// 	return
	// }

	// result.timeSyncV2 = timeTaken3
	// result.processCountSyncV2 = count3
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

	const heartbeatDelay = 200 * time.Millisecond
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

// Create output folder (if nonexistent)
func prepareOutputFolder(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0770)
}

////////////////////////////////////////////
// Visualising csv file into a plot
////////////////////////////////////////////

// Reads a csv file containing the following headers and outputs a graph
// containing two lines: one showing the synchronous (non-polarized) version
// and the other showing the asynchronous (polarized) version.
// CSV Headers: name, caseNumber, timeSyncV1NP, processCountSyncV1NP, timeAsyncV2, processCountAsyncV2

func visualisePlots(resultsFile string, logScale bool) {
	fullResultsPath := filepath.Join(outputFolder, resultsFile)

	// Open the CSV file
	file, err := os.Open(fullResultsPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Skip the header row
	if _, err := reader.Read(); err != nil {
		panic(err)
	}

	// Read in all the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// Create a new plot
	p := plot.New()

	//Get data for the first line (Synchronous, non-polarized version)
	var data plotter.XYs
	for _, record := range records {
		x, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			// panic(err)
			fmt.Println("Couldn't draw plot")
			return
		}

		y, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			// panic(err)
			fmt.Println("Couldn't draw plot")
			return
		}

		if y > 0 {
			var d plotter.XY
			d.X = float64(x)
			d.Y = float64(y)
			data = append(data, d)
		}
	}

	// Add first line to the plot
	// Synchronous semantics (using the non-polarized version)
	line, err := plotter.NewLine(data)
	if err != nil {
		panic(err)
	}
	line.LineStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	line.LineStyle.Width = vg.Points(1)

	p.Add(line)
	p.Legend.Add("Grits (Synchronous sem.)", line)

	// Add second line to the plot

	// Convert the CSV data to a plotter.Values
	var data2 plotter.XYs
	for _, record := range records {
		x, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			// panic(err)
			fmt.Println("Couldn't draw plot")
			return
		}
		y, err := strconv.ParseInt(record[4], 10, 64)
		if err != nil {
			// panic(err)
			fmt.Println("Couldn't draw plot")
			return
		}

		if y > 0 {
			var d plotter.XY
			d.X = float64(x)
			d.Y = float64(y)
			data2 = append(data2, d)
		}
	}

	// Asynchronous semantics (using the polarized, buffered version)
	line2, err := plotter.NewLine(data2)
	if err != nil {
		panic(err)
	}
	line2.LineStyle.Color = color.RGBA{R: 255, G: 127, B: 80, A: 255}
	line2.LineStyle.Width = vg.Points(1)

	p.Add(line2)
	p.Legend.Add("Grits (Asynchronous sem.)", line2)

	// Set the axes labels
	p.Title.Text = filepath.Base(resultsFile)
	p.X.Label.Text = "Count"
	p.Y.Label.Text = "Time taken (µs)"

	if logScale {
		p.Y.Label.Text = "Time taken (µs, log)"
		p.Y.Tick.Marker = plot.LogTicks{Prec: -1}
		p.Y.Scale = plot.LogScale{}
	}

	// Save the plot to a PNG file
	plotFileName := strings.TrimSuffix(filepath.Base(resultsFile), filepath.Ext(resultsFile)) + ".png"
	plotFilePath := filepath.Join(outputFolder, plotFileName)

	err = p.Save(14*vg.Centimeter, 14*vg.Centimeter, plotFilePath)

	if err != nil {
		fmt.Println("Couldn't save graph", err)
	}
}

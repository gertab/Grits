package benchmarks

import (
	"bytes"
	"encoding/csv"
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
//   - timeNonPolarizedSync	        : time taken to evaluate file (using v1)
//   - processCountNonPolarizedSync : number of processes spawn (when using v1)
//   - timeNormalAsync              : time taken to evaluate file (using v2-async)
//   - processCountNormalAsync      : number of processes spawn (when using v2-async)
//   - timeNormalSync               : time taken to evaluate file (using v2-sync)
//   - processCountNormalSync       : number of processes spawn (when using v2-sync)
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
			if detailedOutput {
				fmt.Println(i+1, result.StringShort())
			}
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
func Benchmarks(maxCores int) {
	runtime.GOMAXPROCS(maxCores)
	fmt.Printf("Benchmarking... (using %d cores out of %v)\n\n", maxCores, runtime.NumCPU())

	folder := "nat-double"
	visualisePlots("nat-double-benchmarks-10.csv", true)
	visualisePlots("nat-double-parallel-benchmarks-10.csv", false)
	return
	benchmarkCases := []benchmarkCase{
		{"nat-double-1.grits", 1, 2},
		{"nat-double-2.grits", 2, 2},
		{"nat-double-3.grits", 3, 2},
		{"nat-double-4.grits", 4, 2},
		{"nat-double-5.grits", 5, 2},
		{"nat-double-6.grits", 6, 2},
		{"nat-double-7.grits", 7, 2},
		{"nat-double-8.grits", 8, 2},
		{"nat-double-9.grits", 9, 2},
		{"nat-double-10.grits", 10, 2},
		{"nat-double-11.grits", 11, 2},
		{"nat-double-12.grits", 12, 2},
		{"nat-double-13.grits", 13, 2},
		{"nat-double-14.grits", 14, 2},
		{"nat-double-15.grits", 15, 2},
		{"nat-double-16.grits", 16, 2},
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
				// fmt.Println(result)
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
	name                         string
	invalid                      bool
	caseNumber                   int
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
	if t.caseNumber >= 0 {
		buffer.WriteString(fmt.Sprintf(" (%v)\n", t.caseNumber))
	}
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
	if t.caseNumber >= 0 {
		buffer.WriteString(fmt.Sprintf("%v", t.caseNumber))
	}
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
	buffer.WriteString("caseNumber")
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

func getAverage(allResults []TimingResult) *TimingResult {
	// allResults are not modified
	count := len(allResults)

	if count == 0 {
		return &TimingResult{}
	}

	result := TimingResult{
		allResults[0].name,
		false,
		allResults[0].caseNumber,
		allResults[0].timeNonPolarizedSync,
		allResults[0].processCountNonPolarizedSync,
		allResults[0].timeNormalAsync,
		allResults[0].processCountNormalAsync,
		allResults[0].timeNormalSync,
		allResults[0].processCountNormalSync,
	}

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
// CSV Headers: name, caseNumber, timeNonPolarizedSync, processCountNonPolarizedSync, timeNormalAsync, processCountNormalAsync, timeNormalSync, processCountNormalSync

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
	data := make(plotter.XYs, len(records))
	for i, record := range records {
		x, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			// panic(err)
			fmt.Println("Couldn't draw plot")
			return
		}
		y, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			// panic(err)
			fmt.Println("Couldn't draw plot")
			return
		}
		data[i].X = x
		data[i].Y = y
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
	data2 := make(plotter.XYs, len(records))
	for i, record := range records {
		x, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			// panic(err)
			fmt.Println("Couldn't draw plot")
			return
		}
		y, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			// panic(err)
			fmt.Println("Couldn't draw plot")
			return
		}
		data2[i].X = x
		data2[i].Y = y
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
	plotFileName := filepath.Base(resultsFile) + ".png"
	plotFilePath := filepath.Join(outputFolder, plotFileName)

	if err := p.Save(14*vg.Centimeter, 14*vg.Centimeter, plotFilePath); err != nil {
		panic(err)
	}
}

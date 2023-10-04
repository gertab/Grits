package main

import (
	"bytes"
	"phi/parser"
	"phi/process"
	"sort"
	"testing"
)

const numberOfIterations = 200

// We compare only the process names and rules (i.e. the steps) without comparing the order
func TestSimpleFWDRCV(t *testing.T) {
	// go test -timeout 30s -run ^TestSimpleToken$ phi/cmd
	// Case 1: FWD + RCV
	input := ` 	/* FWD + RCV rule */
		let
		in
		prc[pid1]: send pid2<pid5, self>
		prc[pid2]: -fwd self pid3
		prc[pid3]: <a, b> <- recv self; close a
		end`
	expected := []traceOption{
		{steps{{"pid2", process.FWD}, {"pid2", process.RCV}, {"pid1", process.RCV}}, 2},
		// {steps{{"pid2", process.FWD}, {"pid1", process.RCV}, {"pid2", process.RCV}}, 2},
	}
	checkInputRepeatedly(t, input, expected)

	// sort.Sort(steps(ss3))
}

func TestSimpleSND(t *testing.T) {
	// Case 2: SND
	input := ` 	/* SND rule */
	let
	in
		prc[pid1]: send self<pid3, self>
		prc[pid2]: <a, b> <- recv pid1; close self
	end`
	expected := []traceOption{
		{steps{{"pid1", process.SND}, {"pid2", process.SND}}, 1},
		// {steps{{"pid2", process.SND}, {"pid1", process.SND}}, 1},
	}
	checkInputRepeatedly(t, input, expected)
}

func TestSimpleCUTSND(t *testing.T) {
	// Case 3: CUT + SND
	input := ` 	/* CUT + SND rule */
		let
		in
		prc[pid1]: x <- +new (<a, b> <- recv pid2; close b); close self
		prc[pid2]: send self<pid5, self>
		end`
	expected := []traceOption{
		{steps{{"pid1", process.CUT}, {"x", process.SND}, {"pid2", process.SND}}, 1},
		{steps{{"pid1", process.CUT}, {"pid1", process.SND}, {"x", process.SND}}, 1},
		// {steps{{"pid1", process.CUT}, {"pid2", process.SND}, {"x", process.SND}}, 1},
		// {steps{{"x", process.SND}, {"pid2", process.SND}, {"pid1", process.CUT}}, 1},
		// {steps{{"x", process.SND}, {"pid1", process.CUT}, {"pid2", process.SND}}, 1},
		// {steps{{"pid2", process.SND}, {"pid1", process.CUT}, {"x", process.SND}}, 1},
		// {steps{{"pid2", process.SND}, {"x", process.SND}, {"pid1", process.CUT}}, 1},
	}
	checkInputRepeatedly(t, input, expected)
}

func TestSimpleCUTSNDFWDRCV(t *testing.T) {
	// Case 3: CUT + SND
	input := `   /* CUT + inner blocking SND + FWD + RCV rule */
			let
			in
			prc[pid1]: send pid2<pid5, self>
			prc[pid2]: -fwd self pid3
			prc[pid3]: ff <- -new (send ff<pid5, ff>); <a, b> <- recv self; close self
			end`
	expected := []traceOption{
		// {steps{{"ff", process.SND}, {"pid2", process.FWD}, {"pid2", process.RCV}, {"pid1", process.RCV}}, 3},
		// {steps{{"ff", process.SND}, {"pid1", process.RCV}, {"pid2", process.RCV}, {"pid2", process.FWD}, {"pid3", process.CUT}}, 3},
		// {steps{{"ff", process.SND}, {"pid1", process.RCV}, {"pid2", process.RCV}, {"pid2", process.CUT}, {"pid2", process.FWD}}, 3},
		{steps{{"pid2", process.FWD}, {"pid2", process.RCV}, {"pid1", process.RCV}}, 2},
		{steps{{"pid1", process.RCV}, {"pid2", process.RCV}, {"pid2", process.FWD}, {"pid3", process.CUT}}, 2},
		{steps{{"pid1", process.RCV}, {"pid2", process.RCV}, {"pid2", process.CUT}, {"pid2", process.FWD}}, 2},
	}
	checkInputRepeatedly(t, input, expected)
}

func TestSimpleMultipleFWD(t *testing.T) {
	// Case 4: FWD + FWD + RCV

	input := ` 	/* FWD + RCV rule */
		let
		in
			prc[pid1]: send pid2<pid5, self>
			prc[pid2]: -fwd self pid3
			prc[pid3]: -fwd self pid4
			prc[pid4]: <a, b> <- recv self; close a
		end`
	expected := []traceOption{
		{steps{{"pid2", process.FWD}, {"pid2", process.FWD}, {"pid1", process.RCV}, {"pid2", process.RCV}}, 3},
		{steps{{"pid2", process.FWD}, {"pid3", process.FWD}, {"pid1", process.RCV}, {"pid2", process.RCV}}, 3},
	}
	checkInputRepeatedly(t, input, expected)
}

func TestSimpleMultipleProvidersInitially(t *testing.T) {
	// Case 5: Implicit SPLIT

	input := ` /* SND rule with process having multiple names */
		prc[pa, pb, pc, pd]: send self<pid0, self>
		prc[pid2]: <a, b> <- recv pa; close self
		prc[pid3]: <a, b> <- recv pb; close self
		prc[pid4]: <a, b> <- recv pc; close self
		prc[pid5]: <a, b> <- recv pd; close self
		`
	expected := []traceOption{
		// not sure about these
		{steps{{"pa", process.SND}, {"pb", process.SND}, {"pc", process.SND}, {"pd", process.DUP}, {"pd", process.DUP}, {"pid2", process.SND}, {"pid3", process.SND}, {"pid4", process.SND}, {"pid5", process.SND}}, 9},
		{steps{{"pa", process.SND}, {"pb", process.SND}, {"pc", process.SND}, {"pd", process.SND}, {"pd", process.DUP}, {"pid2", process.SND}, {"pid3", process.SND}, {"pid4", process.SND}, {"pid5", process.SND}}, 4},
	}
	checkInputRepeatedly(t, input, expected)
}

// func TestSimpleSPLITCUT(t *testing.T) {
// 	// Case 6: SPLIT + CUT + SND

// 	input := ` 	/* SPLIT + CUT + SND rule */
// 		let
// 		in
// 			prc[pid0]: <x1, x2> <- +split pid1; close self
// 			prc[pid1]: x <- +new (<a, b> <- recv pid2; close b); close self
// 			prc[pid2]: send self<pid5, self>
// 		end`
// 	expected := []traceOption{
// 		// Either the split finishes before the CUT/SND rules, so the entire tree gets DUPlicated first, thus SND happens twice
// 		{steps{{"pid0", process.SPLIT}, {"x1", process.FWD}, {"pid2", process.DUP}, {"x1", process.CUT}, {"pid2", process.FWD}, {"x2", process.CUT}, {"x1", process.DUP}, {"pid2", process.SND}, {"pid2", process.SND}, {"x", process.SND}, {"x", process.SND}}, 4},
// 		// Or the SND takes place before the SPLIT/DUP, so only one SND is needed
// 		{steps{{"pid0", process.SPLIT}, {"pid1", process.CUT}, {"pid2", process.SND}, {"x", process.SND}}, 1},

// 		// // Added for half done rules (x1...)
// 		// {steps{{"pid0", process.SPLIT}, {"pid1", process.CUT}, {"pid2", process.SND}, {"x", process.SND}, {"x1", process.FWD}}, 1},
// 		// {steps{{"pid0", process.SPLIT}, {"pid2", process.SND}, {"pid2", process.FWD}, {"x1", process.CUT}, {"x1", process.FWD}, {"x1", process.DUP}, {"x2", process.CUT}}, 2},
// 	}
// 	checkInputRepeatedly(t, input, expected)
// }

func TestSimpleSPLITSNDSND(t *testing.T) {
	// Case 7: SPLIT + SND rule (x 2)

	input := ` 	/* Simple SPLIT + SND rule (x 2) */
		let
		in
			prc[pid1]: <a, b> <- +split pid2; <a2, b2> <- recv a; <a3, b3> <- recv b; close self
			prc[pid2]: send self<pid3, self>
		end`
	expected := []traceOption{
		{steps{{"pid1", process.SPLIT}, {"a", process.FWD}, {"a", process.DUP}, {"a", process.SND}, {"pid1", process.SND}, {"b", process.SND}, {"pid1", process.SND}}, 3},
	}
	checkInputRepeatedly(t, input, expected)
}

func TestSimpleSPLITCALL(t *testing.T) {
	// Case 7: SPLIT + CALL

	input := ` /* SPLIT + CALL rule */
			let
				D1(c) =  <a, b> <- recv c; close a
			in
				prc[pid0]: <x1, x2> <- +split pid1; close self
				prc[pid1]: +D1(pid2)
				prc[pid2]: send self<pid3, self>
			end`
	expected := []traceOption{
		{steps{{"pid0", process.SPLIT}, {"pid1", process.CALL}, {"pid2", process.SND}, {"pid2", process.SND}, {"pid2", process.FWD}, {"pid2", process.DUP}, {"x1", process.SND}, {"x1", process.FWD}, {"x1", process.DUP}, {"x2", process.SND}}, 4},
		{steps{{"pid0", process.SPLIT}, {"pid1", process.SND}, {"pid1", process.CALL}, {"pid2", process.SND}}, 1},
		{steps{{"pid0", process.SPLIT}, {"pid2", process.SND}, {"pid2", process.SND}, {"pid2", process.FWD}, {"pid2", process.DUP}, {"x1", process.SND}, {"x1", process.CALL}, {"x1", process.FWD}, {"x1", process.DUP}, {"x2", process.SND}}, 4},
	}
	checkInputRepeatedly(t, input, expected)
}

type step struct {
	processName string
	rule        process.Rule
}

type steps []step

func (a steps) Len() int      { return len(a) }
func (a steps) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a steps) Less(i, j int) bool {
	// Sort by processName and then by rule
	if a[i].processName == a[j].processName {
		return a[i].rule < a[j].rule
	} else {
		return a[i].processName < a[j].processName
	}
}

type traceOption struct {
	trace         steps
	deadProcesses int
}

func checkInputRepeatedly(t *testing.T, input string, expectedOptions []traceOption) {
	// If you increase the number of repetitions to a very high number, make sure to increase
	// the monitor inactiveTimer (to avoid the monitor timing out before terminating).

	done := make(chan bool)
	for i := 0; i < numberOfIterations; i++ {
		go checkInput(t, input, expectedOptions, done)
	}

	for i := 0; i < numberOfIterations; i++ {
		<-done
	}
}

func checkInput(t *testing.T, input string, expectedOptions []traceOption, done chan bool) {
	processes, err := parser.ParseString(input)

	if err != nil {
		t.Errorf("Error during parsing")
		done <- true

		return
	}
	deadProcesses, rulesLog, _ := initProcesses(processes)
	stepsGot := convertRulesLog(rulesLog)

	if len(stepsGot) == 0 {
		t.Errorf("Zero transitions: %s (Increase monitor timeout value) \n", stingifySteps(stepsGot))
	}

	// Make sure that at least the rulesLog match to one of the trance options
	if !compareAllTraces(t, stepsGot, expectedOptions, len(deadProcesses)) {
		// All failed so compare to each expected trace
		printAllTraces(t, stepsGot, expectedOptions, input)
	}
	done <- true
}

func compareAllTraces(t *testing.T, got steps, cases []traceOption, deadProcesses int) bool {
	// Sort trace
	sort.Sort(steps(got))

	for _, c := range cases {
		if len(c.trace) == len(got) {
			// Sort trace
			sort.Sort(steps(c.trace))

			if compareSteps(t, c.trace, got) {
				if c.deadProcesses != deadProcesses {
					t.Errorf("Expected %d dead processes but got %d\n", c.deadProcesses, deadProcesses)
					return false
				}
				return true
			}
		}
	}

	return false
}

func printAllTraces(t *testing.T, got steps, cases []traceOption, input string) {
	for i := range cases {

		if len(cases[i].trace) == len(got) {
			t.Errorf("Got %s, expected %s\n%s\n", stingifySteps(got), stingifySteps(cases[i].trace), input)
		} else {
			t.Errorf("Error: trace length not equal. Expected %d (%s), got %d (%s).\n%s\n", len(cases[i].trace), stingifySteps(cases[i].trace), len(got), stingifySteps(got), input)
		}
	}
}

func compareSteps(t *testing.T, got steps, expected steps) bool {
	if len(got) != len(expected) {
		// t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return false
	}

	for index := range got {
		if got[index] != expected[index] {
			// t.Errorf("got %s, expected %s\n", "sa", "de")
			// t.Errorf("got %s: %s, expected %s: %s\n", got[index].processName, process.RuleString[got[index].rule],
			// expected[index].processName, process.RuleString[expected[index].rule])
			return false
			// tr
		}
	}

	return true
}

func stingifySteps(steps steps) string {
	var buf bytes.Buffer

	for _, s := range steps {
		buf.WriteString(s.processName)
		buf.WriteString(":")
		buf.WriteString(process.RuleString[s.rule])
		buf.WriteString(" ")
	}

	return buf.String()
}

func convertRulesLog(monRulesLog []process.MonitorRulesLog) (log steps) {
	for _, c := range monRulesLog {
		log = append(log, step{processName: c.Process.Providers[0].Ident, rule: c.Rule})
	}

	return log
}

func initProcesses(processes []process.Process) ([]process.Process, []process.MonitorRulesLog, uint64) {

	l := []process.LogLevel{
		process.LOGERROR,
		process.LOGINFO,
		process.LOGPROCESSING,
		process.LOGRULE,
		process.LOGRULEDETAILS,
	}

	debug := true

	re := process.NewRuntimeEnvironment(l, debug, true)

	// fmt.Printf("Initializing %d processes\n", len(processes))

	channels := re.CreateChannelForEachProcess(processes)

	re.SubstituteNameInitialization(processes, channels)

	if debug {
		started := make(chan bool)

		re.InitializeMonitor(started, nil)
		re.InitializeController(started)

		// Ensure that both servers are running
		<-started
		<-started
	}

	re.StartTransitions(processes)

	deadProcesses, rulesLog := re.WaitForMonitorToFinish()

	return deadProcesses, rulesLog, re.ProcessCount
}

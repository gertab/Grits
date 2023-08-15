package main

import (
	"bytes"
	"phi/parser"
	"phi/process"
	"testing"
)

type step struct {
	processName string
	rule        process.Rule
}

type traceOption struct {
	trace         []step
	deadProcesses int
}

var input string
var expected []traceOption

func TestSimpleFWDRCV(t *testing.T) {
	// go test -timeout 30s -run ^TestSimpleToken$ phi/cmd

	// Case 1: FWD + RCV
	input = ` 	/* FWD + RCV rule */
		let
		in
		prc[pid1]: send pid2<pid5, self>
		prc[pid2]: fwd self pid3
		prc[pid3]: <a, b> <- recv self; close a
		end`
	expected = []traceOption{
		{[]step{{"pid2", process.FWD}, {"pid2", process.RCV}, {"pid1", process.RCV}}, 2},
		{[]step{{"pid2", process.FWD}, {"pid1", process.RCV}, {"pid2", process.RCV}}, 2},
	}
	checkInputRepeatedly(t, input, expected)
}

func TestSimpleSND(t *testing.T) {
	// Case 2: SND
	input = ` 	/* SND rule */
	let
	in
		prc[pid1]: send self<pid3, self>
		prc[pid2]: <a, b> <- recv pid1; close self
	end`
	expected = []traceOption{
		{[]step{{"pid1", process.SND}, {"pid2", process.SND}}, 1},
		{[]step{{"pid2", process.SND}, {"pid1", process.SND}}, 1},
	}
	checkInputRepeatedly(t, input, expected)
}

func TestSimpleCUTSND(t *testing.T) {
	// Case 3: CUT + SND
	input = ` 	/* CUT + SND rule */
		let
		in
		prc[pid1]: x <- new (<a, b> <- recv pid2; close b); close self
		prc[pid2]: send self<pid5, self>
		end`
	expected = []traceOption{
		{[]step{{"pid1", process.CUT}, {"x", process.SND}, {"pid2", process.SND}}, 1},
		{[]step{{"pid1", process.CUT}, {"pid1", process.SND}, {"x", process.SND}}, 1},
		{[]step{{"pid1", process.CUT}, {"pid2", process.SND}, {"x", process.SND}}, 1},
		{[]step{{"x", process.SND}, {"pid2", process.SND}, {"pid1", process.CUT}}, 1},
		{[]step{{"x", process.SND}, {"pid1", process.CUT}, {"pid2", process.SND}}, 1},
	}
	checkInputRepeatedly(t, input, expected)

}

func TestSimpleCUTSNDFWDRCV(t *testing.T) {
	// Case 3: CUT + SND
	input = `   /* CUT + inner blocking SND + FWD + RCV rule */
			let
			in
			prc[pid1]: send pid2<pid5, self>
			prc[pid2]: fwd self pid3
			prc[pid3]: ff <- new (send ff<pid5, ff>); <a, b> <- recv self; close self
			end`
	expected = []traceOption{
		{[]step{{"pid3", process.CUT}, {"pid2", process.FWD}, {"pid2", process.RCV}, {"pid1", process.RCV}}, 2},
		{[]step{{"pid3", process.CUT}, {"pid2", process.FWD}, {"pid1", process.RCV}, {"pid2", process.RCV}}, 2},
		{[]step{{"pid2", process.FWD}, {"pid2", process.CUT}, {"pid1", process.RCV}, {"pid2", process.RCV}}, 2},
		{[]step{{"pid2", process.FWD}, {"pid2", process.CUT}, {"pid2", process.RCV}, {"pid1", process.RCV}}, 2},
	}
	checkInputRepeatedly(t, input, expected)
}

func checkInputRepeatedly(t *testing.T, input string, expectedOptions []traceOption) {
	repetitions := 100
	// If you increase the number of repetitions to a very high number, make sure to increase
	// the monitor inactiveTimer (to avoid the monitor timing out before terminating).

	done := make(chan bool)
	for i := 0; i < repetitions; i++ {
		go checkInput(t, input, expectedOptions, done)
	}

	for i := 0; i < repetitions; i++ {
		<-done
	}
}

func checkInput(t *testing.T, input string, expectedOptions []traceOption, done chan bool) {
	processes := parser.ParseString(input)
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

func compareAllTraces(t *testing.T, got []step, cases []traceOption, deadProcesses int) bool {
	for _, c := range cases {
		if len(c.trace) == len(got) {
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

func printAllTraces(t *testing.T, got []step, cases []traceOption, input string) {
	for i := range cases {

		if len(cases[i].trace) == len(got) {
			t.Errorf("Got %s, expected %s\n%s\n", stingifySteps(got), stingifySteps(cases[i].trace), input)
		} else {
			t.Errorf("Error: trace length not equal. Expected %d (%s), got %d (%s).\n%s\n", len(cases[i].trace), stingifySteps(cases[i].trace), len(got), stingifySteps(got), input)
		}
	}
}

func compareSteps(t *testing.T, got []step, expected []step) bool {
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

func stingifySteps(steps []step) string {
	var buf bytes.Buffer

	for _, s := range steps {
		buf.WriteString(s.processName)
		buf.WriteString(":")
		buf.WriteString(process.RuleString[s.rule])
		buf.WriteString(" ")
	}

	return buf.String()
}

func convertRulesLog(monRulesLog []process.MonitorRulesLog) (log []step) {
	for _, c := range monRulesLog {
		log = append(log, step{processName: c.Process.Provider.Ident, rule: c.Rule})
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

		re.InitializeMonitor(started)
		re.InitializeController(started)

		// Ensure that both servers are running
		<-started
		<-started
	}

	re.StartTransitions(processes)

	deadProcesses, rulesLog := re.WaitForMonitorToFinish()

	return deadProcesses, rulesLog, re.ProcessCount
}

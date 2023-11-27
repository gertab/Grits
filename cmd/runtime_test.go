package main

import (
	"bytes"
	"phi/parser"
	"phi/process"
	"sort"
	"sync"
	"testing"
	"time"
)

const numberOfIterations = 30
const monitorTimeout = 40 * time.Millisecond

// Invalidate all cache
// go clean -testcache

// We compare only the process names and rules (i.e. the steps) without comparing the order
func TestSimpleFWDRCV(t *testing.T) {

	// go test -timeout 30s -run ^TestSimpleToken$ phi/cmd
	// Case 1: FWD + RCV
	input := ` 	/* FWD + RCV rule */
		prc[pid1] = send pid2<pid5, self>
		prc[pid2] = fwd self -pid3
		prc[pid3] = <a, b> <- recv self; close a`
	expected := []traceOption{
		{steps{{"pid2", process.FWD}, {"pid2", process.RCV}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleFWDSND(t *testing.T) {
	// go test -timeout 30s -run ^TestSimpleFWDSND$ phi/cmd
	// Case 1: FWD + SND
	input := ` 	/* FWD + RCV rule */
		prc[pid1] = <a, b> <- recv pid2; close a
		prc[pid2] = fwd self +pid3
		prc[pid3] = send self<pid5, self>
		`
	expected := []traceOption{
		{steps{{"pid2", process.FWD}, {"pid1", process.SND}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSND(t *testing.T) {
	// Case 2: SND
	input := ` 	/* SND rule */
	prc[pid1] = send self<pid3, self>
	prc[pid2] = <a, b> <- recv pid1; close self`
	expected := []traceOption{
		{steps{{"pid2", process.SND}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleCUTSND(t *testing.T) {
	// Case 3: CUT + SND
	input := ` 	/* CUT + SND rule */
		prc[pid1] = x <- new (<a, b> <- recv pid2; close b); close self
		prc[pid2] = send self<pid5, self>`
	expected := []traceOption{
		{steps{{"pid1", process.CUT}, {"x", process.SND}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleCUTSNDFWDRCV(t *testing.T) {
	// Case 3: CUT + SND
	input := `   /* CUT + inner blocking SND + FWD + RCV rule */
			prc[pid1] = send pid2<pid5, self>
			prc[pid2] = fwd self -pid3
			prc[pid3] = ff <- new (send ff<pid5, ff>); <a, b> <- recv self; close self`
	expected := []traceOption{
		{steps{{"pid2", process.RCV}, {"pid2", process.FWD}, {"pid3", process.CUT}}},

		// non polarized
		{steps{{"pid2", process.RCV}, {"pid2", process.CUT}, {"pid2", process.FWD}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleMultipleFWD(t *testing.T) {
	// Case 4: FWD + FWD + RCV

	input := ` 	/* FWD + RCV rule */
		prc[pid1] = send pid2<pid5, self>
		prc[pid2] = fwd self -pid3
		prc[pid3] = fwd self -pid4
		prc[pid4] = <a, b> <- recv self; close a`
	expected := []traceOption{
		{steps{{"pid2", process.FWD}, {"pid2", process.FWD}, {"pid2", process.RCV}}},
		{steps{{"pid2", process.FWD}, {"pid3", process.FWD}, {"pid2", process.RCV}}},
		{steps{{"pid2", process.RCV}, {"pid3", process.FWD}, {"pid2", process.RCV}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSPLITCUT(t *testing.T) {
	// Case 5: SPLIT + CUT + SND

	input := ` 	/* SPLIT + CUT + SND rule */
			prc[pid0] = <x1, x2> <- split +pid1; close self
			prc[pid1] = x <- new (<a, b> <- recv pid2; close b); close self
			prc[pid2] = send self<pid5, self>`
	expected := []traceOption{
		// Either the split finishes before the CUT/SND rules, so the entire tree gets DUPlicated first, thus SND happens twice
		{steps{{"pid0", process.SPLIT}, {"pid2", process.DUP}, {"pid2", process.FWD}, {"x", process.SND}, {"x", process.SND}, {"x1", process.CUT}, {"x1", process.DUP}, {"x1", process.FWD}, {"x2", process.CUT}}},

		// Or the SND takes place before the SPLIT/DUP, so only one SND is needed
		{steps{{"pid0", process.SPLIT}, {"pid1", process.CUT}, {"x", process.SND}, {"x1", process.DUP}, {"x1", process.FWD}}},
		// {steps{{"pid0", process.SPLIT}, {"pid2", process.SND}, {"pid2", process.FWD}, {"x1", process.CUT}, {"x1", process.FWD}, {"x1", process.DUP}, {"x2", process.CUT}}, 2},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSPLITSNDSND(t *testing.T) {
	// Case 6: SPLIT + SND rule (x 2)

	input := ` 	/* Simple SPLIT + SND rule (x 2) */
	prc[pid1] : 1 = <a, b> <- split +pid2; 
				<a2, b2> <- recv a; drop a2; drop b2;
				<a3, b3> <- recv b; drop a3; drop b3;
				close self
	prc[pid2] : 1 * 1 = send self<+pid3, +pid4>
	prc[pid3] : 1 = close self
	prc[pid4] : 1 = close self`

	expected := []traceOption{
		{steps{{"pid1", process.SPLIT}, {"a", process.FWD}, {"a", process.DUP}, {"pid1", process.SND}, {"pid1", process.SND}, {"pid1", process.DROP}, {"pid1", process.DROP}, {"pid1", process.DROP}, {"pid1", process.DROP}, {"pid3", process.DUP}, {"pid3", process.FWD}, {"pid4", process.DUP}, {"pid4", process.FWD}}},
		// {steps{{"pid1", process.SPLIT}, {"a", process.FWD}, {"a", process.DUP}, {"pid1", process.SND}, {"pid1", process.SPLIT}}},
		// not sure about this ^
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSPLITRCVRCV(t *testing.T) {
	// Case 6: SPLIT + RCV rule (x 2)

	input := ` 	/* Simple SPLIT + RCV rule (x 2) */
	prc[pid1]= <pid2_first, pid2_second> <- split -pid2; 
					k<- new send pid2_first<pid3, self>;
					wait k;
					send pid2_second<pid4, self>
	prc[pid2] = <a, b> <- recv self; 
						 drop a; 
						 close self`
	expected := []traceOption{
		{steps{{"pid1", process.SPLIT}, {"pid2_first", process.FWD}, {"pid2_first", process.DUP}, {"pid1", process.CUT}, {"pid1", process.CLS}, {"k", process.DROP}, {"pid1", process.DROP}, {"pid2_first", process.RCV}, {"pid2_second", process.RCV}}},
		// {steps{{"pid1", process.SPLIT}, {"pid2_first", process.RCV}, {"pid2_first", process.DUP}, {"pid1", process.CUT}, {"pid1", process.CLS}, {"k", process.DROP}, {"pid1", process.DROP}, {"pid2_first", process.RCV}, {"pid2_second", process.RCV}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSPLITRCVRCVWithTyping(t *testing.T) {
	// Case 6: SPLIT + RCV rule (x 2) (needs types)

	input := ` 	/* Simple SPLIT + RCV rule (x 2) */
	assuming pid3 : 1, pid4 : 1

	prc[pid1] : 1 = <pid2_first, pid2_second> <- split pid2; /* split gets its polarity from the types */
					k : 1 <- new send pid2_first<pid3, self>;
					wait k;
					send pid2_second<pid4, self>
	prc[pid2] : 1 -* 1 = <a, b> <- recv self; 
						 drop a; 
						 close self`
	expected := []traceOption{
		{steps{{"pid1", process.SPLIT}, {"pid2_first", process.FWD}, {"pid2_first", process.DUP}, {"pid1", process.CUT}, {"pid1", process.CLS}, {"k", process.DROP}, {"pid1", process.DROP}, {"pid2_first", process.RCV}, {"pid2_second", process.RCV}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSPLITCALL(t *testing.T) {
	// Case 7: SPLIT + CALL

	input := ` /* SPLIT + CALL rule */
				let D1(c) =  <a, b> <- recv c; close a

				prc[pid0] = <x1, x2> <- split +pid1; close self
				prc[pid1] = D1(pid2)
				prc[pid2] = send self<pid3, pid4>`
	expected := []traceOption{
		{steps{{"pid0", process.SPLIT}, {"pid1", process.CALL}, {"pid2", process.FWD}, {"pid2", process.DUP}, {"x1", process.SND}, {"x1", process.FWD}, {"x1", process.DUP}, {"x2", process.SND}}},
		{steps{{"pid0", process.SPLIT}, {"pid1", process.SND}, {"pid1", process.CALL}}},
		{steps{{"pid0", process.SPLIT}, {"pid2", process.FWD}, {"pid2", process.DUP}, {"x1", process.SND}, {"x1", process.CALL}, {"x1", process.FWD}, {"x1", process.DUP}, {"x2", process.SND}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleMultipleProvidersInitially(t *testing.T) {
	// Case 8: Implicit SPLIT

	input := ` /* SND rule with process having multiple names */
		prc[pa, pb, pc, pd] = send self<+pid0, +pid00>
		prc[pid2] = <a, b> <- recv pa; close self
		prc[pid3] = <a, b> <- recv pb; close self
		prc[pid4] = <a, b> <- recv pc; close self
		prc[pid5] = <a, b> <- recv pd; close self`
	expected := []traceOption{
		{steps{{"pa", process.DUP}, {"pid2", process.SND}, {"pid3", process.SND}, {"pid4", process.SND}, {"pid5", process.SND}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleDUP(t *testing.T) {
	// Case 9: DUP at the top level

	input := ` 
		prc[a, b, c, d] = send self<+x, +y>

		prc[m] = <f, g> <- recv a; wait f; wait g; close self
		prc[n] = <f, g> <- recv b; wait f; wait g; close self
		prc[o] = <f, g> <- recv c; wait f; wait g; close self
		prc[p] = <f, g> <- recv d; wait f; wait g; close self

		prc[x] = close self
		prc[y] = close self`
	expected := []traceOption{
		{steps{{"a", process.DUP}, {"m", process.SND}, {"m", process.CLS}, {"m", process.CLS}, {"n", process.SND}, {"n", process.CLS}, {"n", process.CLS}, {"o", process.SND}, {"o", process.CLS}, {"o", process.CLS}, {"p", process.SND}, {"p", process.CLS}, {"p", process.CLS}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleFunctionCalls(t *testing.T) {
	// Case 9: Function calls, with and without explicit self passed

	input := ` 
		let f(x,y) = send x<y, self>
		let g() = <a, b> <- recv self; wait a; close w
		
		prc[pid1] = f(pid2, pid3)
		prc[pid2] = g()
		prc[pid3] = close self
		
		prc[pid4] = f(self, pid5, pid6)
		prc[pid5] = g(self)
		prc[pid6] = close self
		
		let ff[w, x, y] = send x<y, w>
		let gg[w] = <a, b> <- recv w; wait a; close w
		
		prc[pid7] = ff(pid8, pid9)
		prc[pid8] = gg()
		prc[pid9] = close self
		
		prc[pid10] = ff(self, pid11, pid12)
		prc[pid11] = gg(self)
		prc[pid12] = close self`
	expected := []traceOption{
		{steps{{"pid1", process.CLS}, {"pid1", process.CALL}, {"pid10", process.CLS}, {"pid10", process.CALL}, {"pid11", process.RCV}, {"pid11", process.CALL}, {"pid2", process.RCV}, {"pid2", process.CALL}, {"pid4", process.CLS}, {"pid4", process.CALL}, {"pid5", process.RCV}, {"pid5", process.CALL}, {"pid7", process.CLS}, {"pid7", process.CALL}, {"pid8", process.RCV}, {"pid8", process.CALL}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestExec(t *testing.T) {
	// Case 9: Function calls, with and without explicit self passed

	input := ` 
	type A = 1

	let f() : A = x : A <- new close x; 
					wait x; 
					close self
	
	exec f()`
	expected := []traceOption{
		{steps{{"exec1", process.CALL}, {"exec1", process.CUT}, {"exec1", process.CLS}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestCall(t *testing.T) {
	// Case 10: Function calls

	input := ` 
	type A = &{label : 1}
	type B = 1

	let f1(x : A) : B = x.label<self>
	let f2() : A = case self (label<zz> => close self )

	prc[y] : B = f1(z)
	prc[z] : A = f2()`
	expected := []traceOption{
		{steps{{"y", process.CALL}, {"z", process.CALL}, {"z", process.CSE}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestNew(t *testing.T) {
	// Case 10: New

	input := ` 
	type A = +{label : 1}
	
	prc[y] : 1 = m : A <- new self.label<z>; 
			     case m (label<zz> => wait zz; 
				 close self)
	prc[z] : 1 = close self`

	expected := []traceOption{
		{steps{{"y", process.CUT}, {"y", process.SEL}, {"y", process.CLS}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestFwdPolarity(t *testing.T) {
	// Case 10: Forwards (explicit polarities)

	input := ` 
	type A = &{labelok : 1}

	prc[x] : 1 = y.labelok<self> 
	prc[y] : A = fwd self -z
	prc[z] : A = case self ( labelok<b> => close b )`

	expected := []traceOption{
		{steps{{"y", process.FWD}, {"y", process.CSE}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestFwdCutPolarityWithTyping(t *testing.T) {
	// Case 10: Forwards and cut (implicit polarities from types)

	input := `
	type A = &{labelok : 1}

	prc[x] : 1 = y.labelok<self>
	prc[y] : A = zz : A <- new fwd self z;
				 fwd self zz
				// y is a -ve fwd
	prc[z] : A = case self ( labelok<b> => close b )`

	expected := []traceOption{
		{steps{{"y", process.CUT}, {"y", process.FWD}, {"zz", process.FWD}, {"y", process.CSE}}},
		{steps{{"y", process.CUT}, {"y", process.FWD}, {"y", process.FWD}, {"y", process.CSE}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestFwdPolarityWithTyping(t *testing.T) {
	// Case 10: Forwards (implicit polarities from types)

	input := ` 
	type A = &{labelok : 1}

	prc[x] : 1 = y.labelok<self> 
	prc[y] : A = fwd self z /* this is a negative fwd*/
	prc[z] : A = case self ( labelok<b> => close b )`

	expected := []traceOption{
		{steps{{"y", process.FWD}, {"y", process.CSE}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestNestedSelectWithTyping(t *testing.T) {
	// Case 10: Nested branch/select

	input := ` 
	type A = &{label : +{next : 1}}
	let f1(x : A) : +{next : 1} = x.label<self>
	let f2(y : 1) : A = case self (label<zz> => zz.next<y> )

	prc[x] : +{next : 1} = f1(z)
	prc[z] : A = f2(y)
	prc[y] : 1 = close self
	prc[final] : 1 = case x (next<z> => drop z; close self)`

	expected := []traceOption{
		{steps{{"x", process.CALL}, {"z", process.CALL}, {"z", process.CSE}, {"final", process.SEL}, {"final", process.DROP}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestNestedReceiveWithTyping(t *testing.T) {
	// Case 10: Nested receive/send

	input := ` 
	type A = 1 -* (1 -* (1 * 1))
	prc[x1] : 1 -* (1 * 1) = send z<yy, self>
	prc[x2] : 1 * 1 = send x1<xx, self>
	prc[z] : A = <x, y> <- recv self; 
				 <xx, y> <- recv y; 
				 send y<x, xx>
	prc[xx] : 1 = close self
	prc[yy] : 1 = close self
	
	prc[final] : 1 = <g1, g2> <- recv x2;
					 drop g1;
					 drop g2;
					 close self`

	expected := []traceOption{
		{steps{{"x1", process.RCV}, {"z", process.RCV}, {"final", process.SND}, {"final", process.DROP}, {"final", process.DROP}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestNestedPolarizedFwd(t *testing.T) {
	// Case 10: Nested receive/send

	input := ` 
	// Positive fwd
	type A = +{label1 : B}
	type B = 1
	prc[y] : 1 = case ff (label1<cont> => wait cont; close self)
	prc[ff] : A = fwd self +z
	prc[z] : A = self.label1<x>
	prc[x] : B = close self`

	expected := []traceOption{
		{steps{{"ff", process.FWD}, {"y", process.CLS}, {"y", process.SEL}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)

	input = ` 
	// Positive fwd
	type A = &{label1 : B}
	type B = 1
	prc[y] : 1 = ff.label1<self>
	prc[ff] : A = fwd self -z
	prc[z] : A = case self (label1<cont> => close self)
	prc[x] : B = close self`

	expected = []traceOption{
		{steps{{"ff", process.FWD}, {"ff", process.CSE}}},
	}
	checkInputRepeatedly(t, input, expected, typecheck)

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
	trace steps
	// deadProcesses int
}

func checkInputRepeatedly(t *testing.T, input string, expectedOptions []traceOption, typecheck bool) {
	// If you increase the number of repetitions to a very high number, make sure to increase
	// the monitor inactiveTimer (to avoid the monitor timing out before terminating).

	wg := new(sync.WaitGroup)
	wg.Add(numberOfIterations)
	for i := 0; i < numberOfIterations; i++ {
		go checkInput(t, input, expectedOptions, wg, typecheck)
	}

	wg.Wait()
}

func checkInput(t *testing.T, input string, expectedOptions []traceOption, wg *sync.WaitGroup, typecheck bool) {
	defer wg.Done()

	// Test all operational semantic versions
	execVersions := []process.Execution_Version{
		process.NORMAL_ASYNC,
		process.NORMAL_SYNC,
		process.NON_POLARIZED_SYNC,
	}

	for _, execVersion := range execVersions {
		processes, assumedFreeNames, globalEnv, err := parser.ParseString(input)
		if err != nil {
			t.Errorf("Error during parsing")
			return
		}

		if typecheck {
			err = process.Typecheck(processes, assumedFreeNames, globalEnv)
			if err != nil {
				t.Errorf("typing error: %s", err)
				return
			}
		}

		deadProcesses, rulesLog, _ := initProcesses(processes, globalEnv, execVersion, typecheck)
		stepsGot := convertRulesLog(rulesLog)

		if len(stepsGot) == 0 {
			t.Errorf("Zero transitions: %s (Increase monitor timeout value) \n", stingifySteps(stepsGot))
		}

		// Make sure that at least the rulesLog match to one of the trance options
		if !compareAllTraces(t, stepsGot, expectedOptions, len(deadProcesses)) {
			// All failed so compare to each expected trace
			printAllTraces(t, stepsGot, expectedOptions, input)
		}
	}
}

func compareAllTraces(t *testing.T, got steps, cases []traceOption, deadProcesses int) bool {
	// Sort trace
	sort.Sort(steps(got))

	for _, c := range cases {
		if len(c.trace) == len(got) {
			// Sort trace
			sort.Sort(steps(c.trace))

			if compareSteps(t, c.trace, got) {
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
			return false
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

func initProcesses(processes []*process.Process, globalEnv *process.GlobalEnvironment, execVersion process.Execution_Version, typechecked bool) ([]process.Process, []process.MonitorRulesLog, uint64) {

	l := []process.LogLevel{}

	debug := true

	re := &process.RuntimeEnvironment{
		GlobalEnvironment: globalEnv,
		Debug:             true,
		Color:             true,
		LogLevels:         l,
		Delay:             0,
		ExecutionVersion:  execVersion,
		Typechecked:       typechecked,
	}

	channels := re.CreateChannelForEachProcess(processes)

	re.SubstituteNameInitialization(processes, channels)

	if debug {
		startedWg := new(sync.WaitGroup)
		startedWg.Add(1)

		newMonitor := process.NewMonitor(re, nil)
		newMonitor.SetInactivityTimer(monitorTimeout)

		re.InitializeGivenMonitor(startedWg, newMonitor, nil)

		// Ensure that both servers are running
		startedWg.Wait()
	}

	re.StartTransitions(processes)

	deadProcesses, rulesLog := re.WaitForMonitorToFinish()

	return deadProcesses, rulesLog, re.ProcessCount()
}

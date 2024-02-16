package cmd

import (
	"bytes"
	"phi/parser"
	"phi/process"
	"sort"
	"sync"
	"testing"
	"time"
)

const numberOfIterations = 10
const timeout = 50 * time.Millisecond

// Invalidate all cache
// go clean -testcache
// go test -timeout 30s -run ^TestSimpleToken$ phi/cmd

// We compare only the process names and rules (i.e. the steps) without comparing the order
func TestSimpleFWDRCV(t *testing.T) {
	// FWD + RCV

	input := ` 	/* FWD + RCV rule */
		prc[pid1] = send pid2<pid5, self>
		prc[pid2] = fwd self -pid3
		prc[pid3] = <a, b> <- recv self; close b`

	expected := []traceOption{
		{steps{{"pid2", process.FWD}, {"pid2", process.RCV}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleFWDSND(t *testing.T) {
	// FWD + SND
	input := ` 	/* FWD + RCV rule */
		prc[pid1] = <a, b> <- recv pid2; close self
		prc[pid2] = fwd self +pid3
		prc[pid3] = send self<pid5, self>`

	expected := []traceOption{
		{steps{{"pid2", process.FWD}, {"pid1", process.SND}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSND(t *testing.T) {
	// SND
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
	// CUT + SND
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
	// CUT + SND
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
	// FWD + FWD + RCV

	input := ` 	/* FWD + RCV rule */
		prc[pid1] = send pid2<pid5, self>
		prc[pid2] = fwd self -pid3
		prc[pid3] = fwd self -pid4
		prc[pid4] = <a, b> <- recv self; close b`

	expected := []traceOption{
		{steps{{"pid2", process.FWD}, {"pid2", process.FWD}, {"pid2", process.RCV}}},
		{steps{{"pid2", process.FWD}, {"pid3", process.FWD}, {"pid2", process.RCV}}},
		{steps{{"pid2", process.RCV}, {"pid3", process.FWD}, {"pid2", process.RCV}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSPLITCLS(t *testing.T) {
	// SPLIT + CLS

	input := ` 	/* SPLIT + CLS rule */
	prc[pid0] : 1 = <x1, x2> <- split +pid1; wait x1; wait x2; close self
	prc[pid1] : 1 = close self`

	expected := []traceOption{
		{steps{{"pid0", process.CLS}, {"pid0", process.CLS}, {"pid0", process.SPLIT}, {"x1", process.DUP}, {"x1", process.FWD}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSPLITSND(t *testing.T) {
	// SPLIT + SND

	input := ` 	/* SPLIT + SND rule */
	assuming pid5 : 1, pid6 : 1
	prc[pid0] : 1 = <x1, x2> <- split +pid1; drop +x1; drop +x2; close self
	prc[pid1] : 1 = <y1, y2> <- recv pid2; drop +y1; drop +y2; close self
	prc[pid2] : 1 * 1 = send self<pid5, pid6>`

	expected := []traceOption{
		// pid0:DROP pid0:DROP pid0:SPLIT pid2:DUP pid2:FWD x1:SND x1:DROP x1:DROP x1:DUP x1:FWD x2:SND x2:DROP x2:DROP
		{steps{{"pid0", process.DROP}, {"pid0", process.DROP}, {"pid0", process.SPLIT}, {"pid2", process.DUP}, {"pid2", process.FWD}, {"x1", process.SND}, {"x1", process.DROP}, {"x1", process.DROP}, {"x1", process.DUP}, {"x1", process.FWD}, {"x2", process.SND}, {"x2", process.DROP}, {"x2", process.DROP}}},

		// pid0:DROP pid0:DROP pid0:SPLIT pid1:SND pid1:DROP pid1:DROP x1:DUP x1:FWD
		{steps{{"pid0", process.DROP}, {"pid0", process.DROP}, {"pid0", process.SPLIT}, {"pid1", process.SND}, {"pid1", process.DROP}, {"pid1", process.DROP}, {"x1", process.DUP}, {"x1", process.FWD}}},

		// // pid0:DROP pid0:DROP pid0:SPLIT pid2:DUP pid2:FWD x1:DUP x1:FWD x2:SND x2:DROP x2:DROP
		// {steps{{"pid0", process.DROP}, {"pid0", process.DROP}, {"pid0", process.SPLIT}, {"pid2", process.DUP}, {"pid2", process.FWD}, {"x1", process.DUP}, {"x1", process.FWD}, {"x2", process.SND}, {"x2", process.DROP}, {"x2", process.DROP}}},

		// pid0:DROP pid0:DROP pid0:SPLIT pid1:SND x1:DROP x1:DROP x1:DUP x1:FWD x2:DROP x2:DROP
		{steps{{"pid0", process.DROP}, {"pid0", process.DROP}, {"pid0", process.SPLIT}, {"pid1", process.SND}, {"x1", process.DROP}, {"x1", process.DROP}, {"x1", process.DUP}, {"x1", process.FWD}, {"x2", process.DROP}, {"x2", process.DROP}}},

		// pid0:DROP pid0:DROP pid0:SPLIT pid1:SND pid1:DROP x1:DROP x1:DUP x1:FWD x2:DROP
		{steps{{"pid0", process.DROP}, {"pid0", process.DROP}, {"pid0", process.SPLIT}, {"pid1", process.SND}, {"pid1", process.DROP}, {"x1", process.DROP}, {"x1", process.DUP}, {"x1", process.FWD}, {"x2", process.DROP}}},

		// Either the split finishes before the CUT/SND rules, so the entire tree gets DUPlicated first, thus SND happens twice
		// {steps{{"pid0", process.SPLIT}, {"pid2", process.DUP}, {"pid2", process.FWD}, {"x", process.SND}, {"x", process.SND}, {"x1", process.CUT}, {"x1", process.DUP}, {"x1", process.FWD}, {"x2", process.CUT}}},

		// pid0:DROP pid0:DROP pid0:SPLIT pid1:CUT x:SND x:DROP x:DROP x:CALL x:DUP x:FWD x1:DUP x1:FWD x2:CLS
		// {steps{{"pid0", process.DROP}, {"pid0", process.DROP}, {"pid0", process.SPLIT}, {"pid1", process.CUT}, {"x", process.SND}, {"x", process.DROP}, {"x", process.DROP}, {"x", process.CALL}, {"x", process.DUP}, {"x", process.FWD}, {"x1", process.DUP}, {"x1", process.FWD}, {"x2", process.CLS}}},
		// pid0:DROP pid0:DROP pid0:SPLIT pid2:DUP pid2:FWD x:SND x:DROP x:DROP x:CALL x:CALL x1:CLS x1:CUT x1:DUP x1:FWD x2:CUT
		// {steps{{"pid0", process.DROP}, {"pid0", process.DROP}, {"pid0", process.SPLIT}, {"pid2", process.DUP}, {"pid2", process.FWD}, {"x", process.SND}, {"x", process.DROP}, {"x", process.DROP}, {"x", process.CALL}, {"x", process.CALL}, {"x1", process.CLS}, {"x1", process.CUT}, {"x1", process.DUP}, {"x1", process.FWD}, {"x2", process.CUT}}},

		// Or the SND takes place before the SPLIT/DUP, so only one SND is needed
		// {steps{{"pid0", process.SPLIT}, {"pid1", process.CUT}, {"x", process.SND}, {"x1", process.DUP}, {"x1", process.FWD}}},
		// {steps{{"pid0", process.SPLIT}, {"pid2", process.SND}, {"pid2", process.FWD}, {"x1", process.CUT}, {"x1", process.FWD}, {"x1", process.DUP}, {"x2", process.CUT}}, 2},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSPLITSNDSND(t *testing.T) {
	// SPLIT + SND rule (x 2)

	input := ` 	/* Simple SPLIT + SND rule (x 2) */
	prc[pid1] : 1 = <a, b> <- split +pid2; 
				<a2, b2> <- recv a; drop +a2; drop +b2;
				<a3, b3> <- recv b; drop +a3; drop +b3;
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

	typecheck = true
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSPLITRCVRCV(t *testing.T) {
	// SPLIT + RCV rule (x 2)

	input := ` 	/* Simple SPLIT + RCV rule (x 2) */
	prc[pid1]= <pid2_first, pid2_second> <- split -pid2; 
					k <- new send pid2_first<pid3, self>;
					wait k;
					send pid2_second<pid4, self>
	prc[pid2] = <a, b> <- recv self; 
						 drop a; 
						 close self`
	expected := []traceOption{
		{steps{{"pid1", process.SPLIT}, {"pid2_first", process.FWD}, {"pid2_first", process.DUP}, {"pid1", process.CUT}, {"pid1", process.CLS}, {"k", process.DROP}, {"pid1", process.DROP}, {"pid2_first", process.RCV}, {"pid2_second", process.RCV}}},
		{steps{{"pid1", process.SPLIT}, {"pid2_first", process.RCV}, {"pid2_first", process.DUP}, {"pid1", process.CUT}, {"pid1", process.CLS}, {"k", process.DROP}, {"pid1", process.DROP}, {"pid2_first", process.RCV}, {"pid2_second", process.RCV}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleSPLITRCVRCVWithTyping(t *testing.T) {
	// SPLIT + RCV rule (x 2) (needs types)

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
	// SPLIT + CALL

	input := ` /* SPLIT + CALL rule */
				let D1(c) =  <a, b> <- recv c; close self

				prc[pid0] = <x1, x2> <- split +pid1; close self
				prc[pid1] = D1(pid2)
				prc[pid2] = send self<pid3, pid4>`

	expected := []traceOption{
		// pid0:SPLIT pid1:CALL pid2:DUP pid2:FWD x1:SND x1:DUP x1:FWD x2:SND
		{steps{{"pid0", process.SPLIT}, {"pid1", process.CALL}, {"pid2", process.DUP}, {"pid2", process.FWD}, {"x1", process.SND}, {"x1", process.DUP}, {"x1", process.FWD}, {"x2", process.SND}}},
		// // pid0:SPLIT pid1:CALL pid2:DUP pid2:FWD x1:DUP x1:FWD x2:SND
		// {steps{{"pid0", process.SPLIT}, {"pid1", process.CALL}, {"pid2", process.DUP}, {"pid2", process.FWD}, {"x1", process.DUP}, {"x1", process.FWD}, {"x2", process.SND}}},
		// pid0:SPLIT pid1:SND pid1:CALL x1:DUP x1:FWD
		{steps{{"pid0", process.SPLIT}, {"pid1", process.SND}, {"pid1", process.CALL}, {"x1", process.DUP}, {"x1", process.FWD}}},
		// // pid0:SPLIT pid1:CALL pid2:DUP pid2:FWD x1:SND x1:DUP x1:FWD
		// {steps{{"pid0", process.SPLIT}, {"pid1", process.CALL}, {"pid2", process.DUP}, {"pid2", process.FWD}, {"x1", process.SND}, {"x1", process.DUP}, {"x1", process.FWD}}},

		// {steps{{"pid0", process.SPLIT}, {"pid1", process.CALL}, {"pid2", process.FWD}, {"pid2", process.DUP}, {"x1", process.SND}, {"x1", process.FWD}, {"x1", process.DUP}, {"x2", process.SND}}},
		// {steps{{"pid0", process.SPLIT}, {"pid1", process.SND}, {"pid1", process.CALL}, {"x1", process.DUP}, {"x1", process.FWD}}},
		// {steps{{"pid0", process.SPLIT}, {"pid2", process.FWD}, {"pid2", process.DUP}, {"x1", process.SND}, {"x1", process.CALL}, {"x1", process.FWD}, {"x1", process.DUP}, {"x2", process.SND}}},
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
		//  a:DUP m:SND m:CLS m:CLS n:SND n:CLS n:CLS o:SND o:CLS o:CLS p:SND p:CLS p:CLS x:DUP x:FWD y:DUP y:FWD
		{steps{{"a", process.DUP}, {"m", process.SND}, {"m", process.CLS}, {"m", process.CLS}, {"n", process.SND}, {"n", process.CLS}, {"n", process.CLS}, {"o", process.SND}, {"o", process.CLS}, {"o", process.CLS}, {"p", process.SND}, {"p", process.CLS}, {"p", process.CLS}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},

		// // a:DUP m:SND n:SND n:CLS n:CLS o:SND p:SND x:DUP x:FWD y:DUP y:FWD
		// {steps{{"a", process.DUP}, {"m", process.SND}, {"n", process.SND}, {"n", process.CLS}, {"n", process.CLS}, {"o", process.SND}, {"p", process.SND}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},
		// // a:DUP m:SND m:CLS m:CLS n:SND o:SND p:SND x:DUP x:FWD y:DUP y:FWD
		// {steps{{"a", process.DUP}, {"m", process.SND}, {"m", process.CLS}, {"m", process.CLS}, {"n", process.SND}, {"o", process.SND}, {"p", process.SND}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},
		// // a:DUP m:SND n:SND o:SND o:CLS o:CLS p:SND x:DUP x:FWD y:DUP y:FWD
		// {steps{{"a", process.DUP}, {"m", process.SND}, {"n", process.SND}, {"o", process.SND}, {"o", process.CLS}, {"o", process.CLS}, {"p", process.SND}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},
		// // a:DUP m:SND n:SND o:SND p:SND p:CLS p:CLS x:DUP x:FWD y:DUP y:FWD
		// {steps{{"a", process.DUP}, {"m", process.SND}, {"n", process.SND}, {"o", process.SND}, {"p", process.SND}, {"p", process.CLS}, {"p", process.CLS}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},

		// {steps{{"a", process.DUP}, {"m", process.SND}, {"m", process.CLS}, {"m", process.CLS}, {"n", process.SND}, {"n", process.CLS}, {"n", process.CLS}, {"o", process.SND}, {"o", process.CLS}, {"o", process.CLS}, {"p", process.SND}, {"p", process.CLS}, {"p", process.CLS}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},
	}

	typecheck := false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleDUPTyping(t *testing.T) {
	// DUP at the top level

	input := ` 
		type A = 1
		type B = 1
		
		prc[a, b, c, d] : A * B = send self<x, y>
		
		prc[m] : 1 = <f, g> <- recv a; wait f; wait g; close self
		prc[n] : 1 = <f, g> <- recv b; wait f; wait g; close self
		prc[o] : 1 = <f, g> <- recv c; wait f; wait g; close self
		prc[p] : 1 = <f, g> <- recv d; wait f; wait g; close self
		
		prc[x] : A = close self
		prc[y] : B = close self`

	expected := []traceOption{

		// // a:DUP m:SND n:SND n:CLS n:CLS o:SND p:SND x:DUP x:FWD y:DUP y:FWD
		// {steps{{"a", process.DUP}, {"m", process.SND}, {"n", process.SND}, {"n", process.CLS}, {"n", process.CLS}, {"o", process.SND}, {"p", process.SND}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},

		// // a:DUP m:SND m:CLS m:CLS n:SND o:SND p:SND x:DUP x:FWD y:DUP y:FWD
		// {steps{{"a", process.DUP}, {"m", process.SND}, {"m", process.CLS}, {"m", process.CLS}, {"n", process.SND}, {"o", process.SND}, {"p", process.SND}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},

		// // a:DUP m:SND n:SND o:SND o:CLS o:CLS p:SND x:DUP x:FWD y:DUP y:FWD
		// {steps{{"a", process.DUP}, {"m", process.SND}, {"n", process.SND}, {"o", process.SND}, {"o", process.CLS}, {"o", process.CLS}, {"p", process.SND}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},

		// // a:DUP m:SND n:SND o:SND p:SND p:CLS p:CLS x:DUP x:FWD y:DUP y:FWD
		// {steps{{"a", process.DUP}, {"m", process.SND}, {"n", process.SND}, {"o", process.SND}, {"p", process.SND}, {"p", process.CLS}, {"p", process.CLS}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},

		{steps{{"a", process.DUP}, {"m", process.SND}, {"m", process.CLS}, {"m", process.CLS}, {"n", process.SND}, {"n", process.CLS}, {"n", process.CLS}, {"o", process.SND}, {"o", process.CLS}, {"o", process.CLS}, {"p", process.SND}, {"p", process.CLS}, {"p", process.CLS}, {"x", process.DUP}, {"x", process.FWD}, {"y", process.DUP}, {"y", process.FWD}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestSimpleFunctionCalls(t *testing.T) {
	// Function calls, with and without explicit self passed

	input := ` 
		let f(x,y) = send x<y, self>
		let g() = <a, b> <- recv self; wait a; close b
		
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
	// Function calls, with and without explicit self passed

	input := ` 
	type A = 1

	let f() : A = x : A <- new close x; 
				  wait x; 
				  close self
	
	main f()`
	expected := []traceOption{
		{steps{{"exec1", process.CALL}, {"exec1", process.CUT}, {"exec1", process.CLS}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)

	typecheck = false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestCall(t *testing.T) {
	// Function calls

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

	typecheck = true
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestNew(t *testing.T) {
	// New

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

	typecheck = false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestFwdPolarity(t *testing.T) {
	// Forwards (explicit polarities)

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

	typecheck = true
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
	// Nested branch/select

	input := ` 
	type A = &{label : +{next : 1}}
	let f1(x : A) : +{next : 1} = x.label<self>
	let f2(y : 1) : A = case self (label<zz> => zz.next<y> )

	prc[x] : +{next : 1} = f1(z)
	prc[z] : A = f2(y)
	prc[y] : 1 = close self
	prc[final] : 1 = case x (next<z> => drop +z; close self)`

	expected := []traceOption{
		{steps{{"x", process.CALL}, {"z", process.CALL}, {"z", process.CSE}, {"final", process.SEL}, {"final", process.DROP}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)
	typecheck = false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestNestedReceiveWithTyping(t *testing.T) {
	// Nested receive/send

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
	// typecheck = false
	// checkInputRepeatedly(t, input, expected, typecheck)
}

func TestNestedPolarizedFwd(t *testing.T) {
	// Nested receive/send

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

func TestUpShift(t *testing.T) {
	// Upshift

	input := ` 
		type A = linear 1
		
		prc[a] : linear /\ affine A  = y <- shift self; close y
		prc[b] : linear A = cast a<self>`

	expected := []traceOption{
		{steps{{"a", process.SHF}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)

	typecheck = false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestDownShift(t *testing.T) {
	// DownShift

	input := ` 
	type A = affine 1

	assuming x : A
	prc[a] : affine A = y <- shift b; drop +y; close self
	prc[b] : affine \/ linear A = cast self<x>`

	expected := []traceOption{
		{steps{{"a", process.CST}, {"a", process.DROP}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)

	typecheck = false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestTwoReceives(t *testing.T) {
	// TwoReceives

	input := ` 
		prc[a] : 1             = send b<u, self>
		prc[b] : (1 -* 1) -* 1 = <x, y> <- recv self; send x<z, y>
		prc[u] : 1 -* 1        = <x, y> <- recv self; wait x; close y
		prc[z] : 1             = close self`

	expected := []traceOption{
		{steps{{"b", process.RCV}, {"u", process.RCV}, {"a", process.CLS}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)

	typecheck = false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestFwdDrop(t *testing.T) {
	// FwdDrop

	input := ` 
		prc[a] : 1 -* 1 = <x, y> <- recv self; wait +u; wait +v; drop +x; close y
		prc[b] : 1 = drop -a; close self
		prc[u] : 1 = close self
		prc[v] : 1 = close self`

	expected := []traceOption{
		{steps{{"a", process.FWD}, {"b", process.DROP}}},
		{steps{{"b", process.DROP}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)

	typecheck = false
	checkInputRepeatedly(t, input, expected, typecheck)
}

func TestShifting(t *testing.T) {
	// Shifting

	input := ` 
	prc[a] : lin 1 = x <- shift b; wait x; close self
	prc[b] : rep \/ lin 1 = cast self<c>
	prc[c] : rep 1 = close self`

	expected := []traceOption{
		{steps{{"a", process.CST}, {"a", process.CLS}}},
	}

	typecheck := true
	checkInputRepeatedly(t, input, expected, typecheck)

	typecheck = false
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

	for _, c := range expectedOptions {
		// Sort trace
		sort.Sort(steps(c.trace))
	}

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

		// No logs during testing
		globalEnv.LogLevels = []process.LogLevel{}

		// processCount, deadProcessCount
		deadProcesses, rulesLog, _, _ := initProcesses(processes, globalEnv, execVersion, typecheck)
		stepsGot := convertRulesLog(rulesLog)

		if len(stepsGot) == 0 {
			t.Errorf("Zero transitions: %s (Increase timeout value) \n", stingifySteps(stepsGot))
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
			if compareSteps(t, c.trace, got) {
				return true
			}
		}
	}

	return false
}

func printAllTraces(t *testing.T, got steps, cases []traceOption, input string) {

	var buffer bytes.Buffer

	buffer.WriteString("Got ")
	buffer.WriteString(stingifySteps(got))
	buffer.WriteString("[WRONG] , expected ")

	for i := range cases {
		buffer.WriteString(stingifySteps(cases[i].trace))

		if i < len(cases)-1 {
			buffer.WriteString("\n or ")
		}
	}
	t.Errorf(buffer.String())
}

func compareSteps(t *testing.T, got steps, expected steps) bool {
	if len(got) != len(expected) {
		// t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return false
	}

	gotString := stingifySteps(got)
	expectedString := stingifySteps(expected)

	return gotString == expectedString
}

func stingifySteps(steps steps) string {
	var buffer bytes.Buffer

	for _, s := range steps {
		buffer.WriteString(s.processName)
		buffer.WriteString(":")
		buffer.WriteString(process.RuleString[s.rule])
		buffer.WriteString(" ")
	}

	return buffer.String()
}

func convertRulesLog(monRulesLog []process.MonitorRulesLog) (log steps) {
	for _, c := range monRulesLog {
		log = append(log, step{processName: c.Process.Providers[0].Ident, rule: c.Rule})
	}

	return log
}

func initProcesses(processes []*process.Process, globalEnv *process.GlobalEnvironment, execVersion process.Execution_Version, typechecked bool) ([]process.Process, []process.MonitorRulesLog, uint64, uint64) {

	re, _, cancel := process.NewRuntimeEnvironment()
	defer cancel()

	re.GlobalEnvironment = globalEnv
	re.UseMonitor = true
	re.Color = true
	re.Delay = 0
	re.ExecutionVersion = execVersion
	re.Typechecked = typechecked

	channels := re.CreateChannelForEachProcess(processes)

	re.SubstituteNameInitialization(processes, channels)

	// Start monitor
	startedWg := new(sync.WaitGroup)
	startedWg.Add(1)
	newMonitor := process.NewMonitor(re, nil)
	re.InitializeGivenMonitor(startedWg, newMonitor, nil)
	startedWg.Wait()

	// Start heartbeat receiver which cleans up any remaining processes
	go re.HeartbeatReceiver(timeout, cancel)

	// Initiate transtitions
	re.StartTransitions(processes)

	select {
	case <-re.Ctx().Done():
		deadProcesses, rulesLog := re.StopMonitor()
		return deadProcesses, rulesLog, re.ProcessCount(), re.DeadProcessCount()
	case <-re.ErrorChan():
		return []process.Process{}, []process.MonitorRulesLog{}, re.ProcessCount(), re.DeadProcessCount()
	}
}

package main

import (
	"phi/parser"
	"phi/process"
	"testing"
)

// Checks related to the typechecker

// Parses the cases and expects each to pass the typechecking phase successfully (or not)
func runThroughTypechecker(t *testing.T, cases []string, pass bool) {
	for i, c := range cases {
		processes, processesFreeNames, globalEnv, err := parser.ParseString(c)

		if err != nil {
			t.Errorf("compilation error in case #%d: %s\n", i, err.Error())
		}

		err = process.Typecheck(processes, processesFreeNames, globalEnv)

		if pass {
			if err != nil {
				t.Errorf("expected no type errors in case #%d, but found %s\n", i, err.Error())
			}
		} else {
			if err == nil {
				t.Errorf("expected type error in case #%d, but didn't find any\n", i)
			}
		}
	}
}

// Typechecker -> these programs should pass the typechecker
func TestTypecheckCorrectSendReceive(t *testing.T) {

	cases := []string{
		/* 0 */ "type A = 1",
		"let f() : 1 = close self",
		// send
		// MulR
		"let f(a : 1, b : 1) : 1 * 1 = send self<a, b>",
		`type A = +{l1 : 1}
			type B = 1 -* 1
			let f(a : A, b : B) : A * B = send self<a, b>`,
		`type A = 1 * 1
		let f(a : 1, b : 1) : A = send self<a, b>`,
		// ImpL
		"let f2(a : 1 -* 1, b : 1) : 1 = send a<b, self>",
		/* 5 */ `type A = +{l1 : 1}
			type B = 1 * 1
			let f(a : A -* B, b : A) : B = send a<b, self>`,
		`type A = 1 -* 1
		let f2(a : A, b : 1) : 1 = send a<b, self>`,
		// receive
		// ImpR
		"let f1() : 1 -* 1 = <x, y> <- recv self; wait x; close y",
		"let f2(b : 1) : 1 -* (1 * 1) = <x, y> <- recv self; send y<x, b>",
		"let f1() : (1 * 1) -* 1 = <x, y> <- recv self; <x2, y2> <- recv x; wait x2; wait y2; close y",
		`type A = 1 -* 1
		let f1() : A = <x, y> <- recv self; wait x; close y`,
		// MulL
		"let f1(u : 1 * 1) : 1 = <x, y> <- recv u; wait x; wait y; close self",
		`type A = 1 * 1
		let f1(u : A) : 1 = <x, y> <- recv u; wait x; wait y; close self`,
	}

	runThroughTypechecker(t, cases, true)
}

// Typechecker -> these programs should fail
func TestTypecheckIncorrectSendReceive(t *testing.T) {
	cases := []string{
		/* 0 */ "type A = B",
		"prc[a] : A = close self",
		"let f() : 1 -* A = close self",
		// MulL (extra non used names)
		"let f(c : 1, a : 1, b : 1) : 1 * 1 = send self<a, b>",
		// MulL (missing names)
		"let f(b : 1) : 1 * 1 = send self<a, b>",
		// MulL (incorrect self type)
		/* 5 */ "let f(a : 1, b : 1) : &{a : 1} = send self<a, b>",
		"let f(a : 1, b : 1) : 1 * &{a : 1} = send self<a, b>",
		// MulL (incorrect self type)
		"let f(a : &{a : 1}, b : &{b : 1}) : &{a : 1} * 1 = send self<a, b>",
		// MulL (wrong types)
		"let f(a : 1 -* 1, b : 1) : 1 * 1 = send self<a, b>",
		"let f(a : 1, b : &{a: 1}) : 1 * 1 = send self<a, b>",
		/* 10 */ "let f(a : 1, b : 1) : 1 * (1 * 1) = send self<a, b>",
		// ImpL
		"let f2(a : 1 -* 1, b : 1) : +{x : 1} = send a<b, self>",
		"let f2(a : 1 -* 1, b : +{x : 1}) : 1 = send a<b, self>",
		"let f2(a : 1 * 1, b : 1) : 1 = send a<b, self>",
		// MulR/ImpL
		"let f2(a : 1 -* 1, b : 1) : 1 = send a<b, c>",
		/* 15 */ "let f2(a : 1 -* 1, c : 1) : 1 = send a<self, c>",
		// ImpR
		"let f2() : 1 * 1 = <x, y> <- recv self; close y",
		"let f2(b : 1) : 1 -* (1 * 1) = <x, y> <- recv self; send x<y, b>",
		"let f1() : 1 -* 1 = <x, self> <- recv self; close y",
		// MulL
		"let f1(u : 1 -* 1) : 1 = <x, y> <- recv u; close y",
		/* 20 */
		"let f1(u : 1 * 1) : 1 = <self, y> <- recv u; close y",
		"let f1() : (1 -* 1) -* 1 = <x, y> <- recv self; <x2, y2> <- recv x; close y",
		`type B = &{label33 : 1}
		let f2(x : +{label1 : 1, label2 : 1, label3 : 1}) : B -* (1 * B) = 
					<x, y> <- recv self; send y<a, x>`,
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectUnit(t *testing.T) {

	cases := []string{
		// Close
		// EndR
		"let f1() : 1 = close self",
		` type A = 1
		let f1() : A = close self`,
		// EndL
		"let f1(x : 1) : 1 = wait x; close self",
		"let f1() : 1 -* 1 = <x, y> <- recv self; wait x; close y",
		`type A = 1
		let f1(x : A) : A = wait x; close self`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckIncorrectUnit(t *testing.T) {
	cases := []string{
		// EndR
		"let f1(u : 1) : 1 = close u",
		"let f1(u : 1 * 1) : 1 = close self",
		"let f1() : 1 * 1 = close self",
		// EndL
		"let f1() : 1 = wait self; close self",
		"let f1(g : 1 * 1, x : 1) : 1 = wait x; close self",
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectForward(t *testing.T) {

	cases := []string{
		// ID
		"let f1(x : 1 * 1) : 1 * 1 = fwd self x",
		"let f1(x : 1 * 1) : 1 * 1 = fwd self x",
		"let f1() : 1 -* 1 = <x, y> <- recv self; fwd y x",
		"let f1(g : (&{a : 1})) : 1 -* (&{a : 1}) = <x, y> <- recv self; wait x; fwd y g",
		`type A = 1 * 1
		type B = 1 * 1
		let f1(x : A) : B = fwd self x`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckIncorrectForward(t *testing.T) {
	cases := []string{
		// ID
		"let f1(x : 1 * 1) : 1 = fwd self x",
		"let f1(x : 1 * 1) : 1 -* 1 = fwd self x",
		"let f1(x : &{hello : 1}) : 1 = fwd self x",
		"let f1(x : 1 * 1) : 1 * 1 = fwd x self",
		"let f1(x : 1 * 1, y : 1) : 1 * 1 = fwd self x",
		"let f1(g : (+{a : 1})) : 1 -* (&{a : 1}) = <x, y> <- recv self; wait x; fwd y g",
		"let f1(x : 1 * 1, y : 1 * 1) : 1 * 1 = fwd x y",
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectDrop(t *testing.T) {

	cases := []string{
		// Drop
		"let f1(x : 1 * 1, g : &{a : 1}) : 1 * 1 = drop g; fwd self x",
		`type A = 1 * 1
		type B = 1 * 1
		type G = &{a : 1}
		let f1(x : A, g : G) : B = drop g; fwd self x`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckIncorrectDrop(t *testing.T) {
	cases := []string{
		// Drop
		"let f1(x : 1 * 1) : 1 * 1 = drop g; fwd self x",
		"let f1(x : 1 * 1, g : &{a : 1}) : 1 * 1 = drop self; fwd self x",
		// Drop and use later
		"let f1() : 1 -* 1 = <x, y> <- recv self; drop x; wait x; close y",
		"let f1() : 1 -* 1 = drop x; <x, y> <- recv self;  wait x; close y",
		// Missed drop
		"let f1(x : 1 * 1, g : &{a : 1}) : 1 * 1 = fwd self x",
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectSelect(t *testing.T) {

	cases := []string{
		// Select
		// IChoiceR
		"let f1(cont : 1) : +{label1 : 1} = self.label1<cont>",
		"let f1(cont : 1 -* 1) : +{label0 : 1, label1 : 1 -* 1} = self.label1<cont>",
		// EChoiceL
		"let f1(to_c : &{label1 : 1}) : 1 = to_c.label1<self>",
		"let f1(to_c : &{label0 : 1, label1 : 1 -* 1}) : 1 -* 1 = to_c.label1<self>",
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckIncorrectSelect(t *testing.T) {
	cases := []string{
		// Select
		// IChoiceR
		"let f1(cont : 1) : &{label1 : 1} = self.label1<cont>",
		"let f1(cont : 1 -* 1) : +{label0 : 1, label1 : 1 -* 1} = self.otherLabel<cont>",
		"let f1(cont : 1) : +{label1 : 1} = a.label1<cont>",
		// EChoiceL
		"let f1(to_c : +{label1 : 1}) : 1 = to_c.label1<self>",
		"let f1(to_c : &{label0 : 1, label1 : 1 -* 1}) : 1 -* 1 = to_c.label2<self>",
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectBranch(t *testing.T) {

	cases := []string{
		// Branch
		// EChoiceR
		"let f1() : &{label1 : 1} = case self (label1<a> => close self)",
		`let f2() : &{label1 : 1, label2 : 1, label3 : 1} = 
					case self (label1<a> => close self
							  |label2<a> => close self
							  |label3<a> => close self)`,
		`let f2() : &{label1 : 1, label2 : 1, label3 : 1} = 
					case self (label2<a> => close self
							  |label3<a> => close self
							  |label1<a> => close self)`,
		`type A = 1 -* (1 * 1)
		let f1(b : 1) : &{label1 : A} = 
				case self (
					label1<a> => <x, y> <- recv a; send y<x, b>
				)`,
		`let f1() : &{label1 : (1 * 1) -* 1 } = 
		case self (
			label1<a> => <x, y> <- recv a; 
						 <xx, yy> <- recv x; 
						 wait xx; 
						 wait yy; 
						 close y
						//  close self
		) `,
		// IChoiceL
		"let f1(x : +{label1 : 1}) : 1 = case x (label1<a> => wait a; close self)",
		`let f2(x : +{label1 : 1, label2 : 1, label3 : 1}) : 1 = 
		case x (label1<a> => wait a; close self
			   |label2<a> => drop a; close self
			   |label3<a> => wait a; close self)`,

		`let f1(x : +{label1 : (1 * 1) * 1 }) : 1 = 
		case x (label1<a> => <x, y> <- recv a; 
						     <xx, yy> <- recv x; 
						     wait xx; 
						     wait y; 
						     wait yy; 
						     close self)`,
		`type B = &{label33 : 1}
		let f2(x : +{label1 : 1, label2 : 1, label3 : 1}) : B -* (1 * B) = 
							case x ( label2<a> => <x, y> <- recv self; send y<a, x>
								| label3<a> => <x, y> <- recv self; send y<a, x>
								| label1<a> => <x, y> <- recv self; send y<a, x>)`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckIncorrectBranch(t *testing.T) {
	cases := []string{
		// Branch
		// EChoiceR
		// Wrong type
		"let f1() : +{label1 : 1} = case self (label1<a> => close self)",
		// Missing labels in type
		`let f2() : &{label1 : 1, label3 : 1} = 
			case self (label1<a> => close self
					  |label2<a> => close self
					  |label3<a> => close self)`,
		// Missing labels in branch
		`let f2() : &{label1 : 1, label2 : 1, label3 : 1} = 
					case self (label1<a> => close self
							|label3<a> => close self)`,
		// Incorrect inner type
		`type A = 1 -* (1 -* 1)
		let f1(b : 1) : &{label1 : A} = 
				case self (
					label1<a> => <x, y> <- recv a; send y<x, b>
				)`,
		`let f1(x : +{label1 : (1 -* 1) * 1 }) : 1 = 
			case x (
				label1<a> => <x, y> <- recv a; 
							<xx, yy> <- recv x; 
							wait xx; 
							wait y; 
							wait yy; 
							close self
			)`,
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectFunctionCall(t *testing.T) {

	cases := []string{
		// FunctionCall
		`let f1() : 1 = close self
		let f2() : 1 = +f1()`,
		`let f3(x : 1 -* 1, y : 1) : 1 = send x<y, self>
		let f4(x2 : 1 -* 1, y2 : 1) : 1 = +f3(x2, y2)`,
		`let f5(x : 1 -* &{label : 1}, y : 1) : &{label : 1} = send x<y, self>
		let f6(x2 : 1 -* &{label : 1}, y2 : 1) : &{label : 1} = +f5(x2, y2)`,
		// Explicit self
		`let f1() : 1 = close self
		let f2() : 1 -* 1 = <x, y> <- recv self; drop x; +f1(y)`,
		`let f3(x : 1 -* 1, y : 1) : 1 = send x<y, self>
		let f4(x2 : 1 -* 1, y2 : 1) : 1 = +f3(self, x2, y2)`,
		`let f5(x : 1 -* &{label : 1}, y : 1) : &{label : 1} = send x<y, self>
		let f6(x2 : 1 -* &{label : 1}, y2 : 1) : &{label : 1} = +f5(self, x2, y2)`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckIncorrectFunctionCall(t *testing.T) {
	cases := []string{
		// FunctionCall
		`let f5(x : 1 -* &{label : 1}, y : 1) : &{label : 1} = send x<y, self>
		let f6(x2 : 1 -* 1, y2 : 1) : &{label : 1} = +f5(x2, y2)`,
		// Explicit self
		`let f1() : 1 = close self
		let f2() : 1 -* (1 * 1) = <x, y> <- recv self; +f1(y)`,
		`let f3(x : 1 -* 1, y : 1) : 1 * 1 = send x<y, self>
		let f4(x2 : 1 -* 1, y2 : 1) : 1 = +f3(y2)`,
		`let f5(x : 1 -* &{label : 1}, y : 1) : &{label : 1} = send x<y, self>
		let f6(x2 : 1 -* &{label2 : 1}, y2 : 1) : &{label : 1} = +f5(x2, y2)`,
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectProcessDefinitions(t *testing.T) {

	cases := []string{
		// send, MulR
		"prc[x] : 1 * 1 = send self<a, b> % a : 1, b : 1",
		`type A = +{l1 : 1}
		type B = 1 -* 1
		prc[y] : A * B = send self<a, b> % a : A, b : B`,
		`type A = 1 * 1
		prc[x] : A = send self<a, b> 		% a : 1, b : 1`,
		// ImpL
		`let f2(a : 1 -* 1, b : 1) : 1 = send a<b, self>
		prc[x] : 1 = f2(aa, bb) 		% aa : 1 -* 1, bb : 1
		prc[y] : 1 = f2(self, aa, bb) 		% aa : 1 -* 1, bb : 1`,
		// With explicit self
		`let f3[w: 1, a : 1 -* 1, b : 1] = send a<b, self>
		prc[x] : 1 = f3(aa, bb) 			% aa : 1 -* 1, bb : 1
		prc[y] : 1 = f3(self, aa, bb) 	% aa : 1 -* 1, bb : 1`,
		`type A = +{l1 : 1}
		type B = 1 * 1
		prc[x]: B = send a<b, self> 		% b : A, a : A -* B`,
		`prc[x] : (1 * 1) -* 1 = 
				<x, y> <- recv self; 
				<x2, y2> <- recv x; 
				wait x2; 
				wait y2; 
				close y`,
		`prc[x] : 1 -* (1 * 1) = 
		<x, y> <- recv self; 
		send y<x, b>  			%  b : 1`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckCorrectProcessDefinitionsIncorrect(t *testing.T) {

	cases := []string{
		// send, MulR
		"prc[x] : 1 = send self<a, b> % a : 1, b : 1",
		"prc[x] : 1 * 1 = send self<a, b> % a : 1, b : 1 -* 1",
		"prc[x] : 1 * 1 = send self<a, b> % a : 1",
		"prc[x] = send self<a, b> % a : 1, b : 1",
		"prc[x] = send self<a, b> % b : 1",
		"prc[x] : 1 * 1 = send self<a, b> % a : 1, b : 1, c : 1",
		"prc[x] : 1 * 1 = send self<a, b> % c : 1",
		"prc[x] : 1 * 1 = send self<a, b>",
		"prc[x] = send self<a, b>",

		// With explicit self
		`let f3[w: 1, a : 1 -* 1, b : 1] = send a<b, self>
		prc[x] : 1 * 1= f3(aa, bb) 			% aa : 1 -* 1, bb : 1
		prc[y] : 1 = f3(self, aa, bb) 		% aa : 1 -* 1, bb : 1`,
		`let f3[w: 1, a : 1 -* 1, b : 1] = send a<b, self>
		prc[x] : 1 = f3(aa, bb) 			% aa : 1 -* 1, bb : 1
		prc[y] : 1 = f3(self, aa, bb) 		% aa : 1 , bb : 1`,
		`let f3[w: 1, a : 1 -* 1, b : 1] = send a<b, self>
		prc[x] : 1 = f3(aa, bb) 			% aa : 1 -* 1, bb : 1
		prc[y] : 1 = f3(self, aa, bb) 		% bb : 1`,
		`let f3[w: 1, a : 1 -* 1, b : 1] = send a<b, self>
		prc[x] : 1 = f3(aa, bb) 			% aa : 1 -* 1, bb : 1
		prc[y] : 1 * 1 = f3(self, aa, bb) 	% aa : 1 -* 1, bb : 1`,
	}

	runThroughTypechecker(t, cases, false)
}

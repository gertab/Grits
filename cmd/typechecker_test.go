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
		processes, assumedFreeNames, globalEnv, err := parser.ParseString(c)

		if err != nil {
			t.Errorf("compilation error in case #%d: %s\n", i, err.Error())
		}

		err = process.Typecheck(processes, assumedFreeNames, globalEnv)

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
		"type A = 1",
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
		`type A = +{l1 : 1}
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
		`assuming b : 1
		 prc[y] : 1 -* (1 -* (1 -* (1 * 1))) = <x, y> <- recv self;
								 drop x;
								 <x, y> <- recv y;
								 drop x;
								 <x, y> <- recv y;
								 send y<x, b>`,
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
		"type A = B",
		"prc[a] : A = close self",
		"let f() : 1 -* A = close self",
		// MulL (extra non used names)
		"let f(c : 1, a : 1, b : 1) : 1 * 1 = send self<a, b>",
		// MulL (missing names)
		"let f(b : 1) : 1 * 1 = send self<a, b>",
		// MulL (incorrect self type)
		"let f(a : 1, b : 1) : &{a : 1} = send self<a, b>",
		"let f(a : 1, b : 1) : 1 * &{a : 1} = send self<a, b>",
		// MulL (incorrect self type)
		"let f(a : &{a : 1}, b : &{b : 1}) : &{a : 1} * 1 = send self<a, b>",
		// MulL (wrong types)
		"let f(a : 1 -* 1, b : 1) : 1 * 1 = send self<a, b>",
		"let f(a : 1, b : &{a: 1}) : 1 * 1 = send self<a, b>",
		"let f(a : 1, b : 1) : 1 * (1 * 1) = send self<a, b>",
		// ImpL
		"let f2(a : 1 -* 1, b : 1) : +{x : 1} = send a<b, self>",
		"let f2(a : 1 -* 1, b : +{x : 1}) : 1 = send a<b, self>",
		"let f2(a : 1 * 1, b : 1) : 1 = send a<b, self>",
		// MulR/ImpL
		"let f2(a : 1 -* 1, b : 1) : 1 = send a<b, c>",
		"let f2(a : 1 -* 1, c : 1) : 1 = send a<self, c>",
		// ImpR
		"let f2() : 1 * 1 = <x, y> <- recv self; close y",
		"let f2(b : 1) : 1 -* (1 * 1) = <x, y> <- recv self; send x<y, b>",
		"let f1() : 1 -* 1 = <x, self> <- recv self; close y",
		// MulL
		"let f1(u : 1 -* 1) : 1 = <x, y> <- recv u; close y",
		"let f1(u : 1 * 1) : 1 = <self, y> <- recv u; close y",
		"let f1() : (1 -* 1) -* 1 = <x, y> <- recv self; <x2, y2> <- recv x; close y",
		`type B = &{label33 : 1}
		 let f2(x : +{label1 : 1, label2 : 1, label3 : 1}) : B -* (1 * B) =
					<x, y> <- recv self; send y<a, x>`,
		"prc[c] : 1 -* 1 = <x, x> <- recv self; close x",
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectUnit(t *testing.T) {

	cases := []string{
		// Close
		// EndR
		"let f1() : 1 = close self",
		"let f1[w : 1] = close w",
		`type A = 1
		 let f1() : A = close self`,
		// "prc[x] : 1 = close x",
		"prc[x] : 1 = close self",
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
		// Assuming a
		`type A = 1
		 assuming a : A
		 prc[a] : A = close self
		 prc[b] : 1 = wait a; close self`,
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
		 let f2() : 1 = f1()`,
		`let f3(x : 1 -* 1, y : 1) : 1 = send x<y, self>
		 let f4(x2 : 1 -* 1, y2 : 1) : 1 = f3(x2, y2)`,
		`let f5(x : 1 -* &{label : 1}, y : 1) : &{label : 1} = send x<y, self>
		 let f6(x2 : 1 -* &{label : 1}, y2 : 1) : &{label : 1} = f5(x2, y2)`,
		// Explicit self
		`let f1() : 1 = close self
		 let f2() : 1 -* 1 = <x, y> <- recv self; drop x; f1(y)`,
		`let f3(x : 1 -* 1, y : 1) : 1 = send x<y, self>
		 let f4(x2 : 1 -* 1, y2 : 1) : 1 = f3(self, x2, y2)`,
		`let f5(x : 1 -* &{label : 1}, y : 1) : &{label : 1} = send x<y, self>
		 let f6(x2 : 1 -* &{label : 1}, y2 : 1) : &{label : 1} = f5(self, x2, y2)`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckIncorrectFunctionCall(t *testing.T) {
	cases := []string{
		// FunctionCall
		`let f5(x : 1 -* &{label : 1}, y : 1) : &{label : 1} = send x<y, self>
		 let f6(x2 : 1 -* 1, y2 : 1) : &{label : 1} = f5(x2, y2)`,
		// Explicit self
		`let f1() : 1 = close self
		 let f2() : 1 -* (1 * 1) = <x, y> <- recv self; f1(y)`,
		`let f3(x : 1 -* 1, y : 1) : 1 * 1 = send x<y, self>
		 let f4(x2 : 1 -* 1, y2 : 1) : 1 = f3(y2)`,
		`let f5(x : 1 -* &{label : 1}, y : 1) : &{label : 1} = send x<y, self>
		 let f6(x2 : 1 -* &{label2 : 1}, y2 : 1) : &{label : 1} = f5(x2, y2)`,
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectProcessDefinitions(t *testing.T) {

	cases := []string{
		// send, MulR
		`assuming a : 1, b : 1
		 prc[x] : 1 * 1 = send self<a, b>`,
		`type A = +{l1 : 1}
		 type B = 1 -* 1
		 assuming a : A, b : B
		 prc[y] : A * B = send self<a, b>`,
		`type A = 1 * 1
		 assuming a : 1, b : 1
		 prc[x] : A = send self<a, b>`,
		// ImpL
		`let f2(a : 1 -* 1, b : 1) : 1 = send a<b, self>
		 assuming a1 : 1 -* 1, a2 : 1 -* 1, b1 : 1, b2 : 1
		 prc[x] : 1 = f2(a1, b1)
		 prc[y] : 1 = f2(self, a2, b2)`,
		// With explicit self
		`let f3[w: 1, a : 1 -* 1, b : 1] = send a<b, self>
		 prc[x] : 1 = f3(aa, bb)
		 assuming aa : 1 -* 1, aaa : 1 -* 1, bb : 1
		 assuming b2 : 1
		 prc[y] : 1 = f3(self, aaa, b2)`,
		`type A = +{l1 : 1}
		 type B = 1 * 1
		 assuming b : A, a : A -* B
		 prc[x]: B = send a<b, self>`,
		`prc[x] : (1 * 1) -* 1 =
				<x, y> <- recv self;
				<x2, y2> <- recv x;
				wait x2;
				wait y2;
				close y`,
		`assuming b : 1
		 prc[x] : 1 -* (1 * 1) =
		 <x, y> <- recv self;
		 send y<x, b> `,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckCorrectProcessDefinitionsIncorrect(t *testing.T) {

	cases := []string{
		// send, MulR
		`assuming a : 1, b : 1
		 prc[x] : 1 = send self<a, b>`,
		`assuming a : 1, b : 1 -* 1
		 prc[x] : 1 * 1 = send self<a, b>`,
		`assuming a : 1
		 prc[x] : 1 * 1 = send self<a, b>`,
		`assuming a : 1, b : 1
		 prc[x] = send self<a, b>`,
		`assuming b : 1
		 prc[x] = send self<a, b>`,
		`assuming a : 1, b : 1, c : 1
		 prc[x] : 1 * 1 = send self<a, b>`,
		`assuming c : 1
		 prc[x] : 1 * 1 = send self<a, b>`,
		"prc[x] : 1 * 1 = send self<a, b>",
		"prc[x] = send self<a, b>",

		// With explicit self
		`assuming aa : 1 -* 1, bb : 1, x : 1 -* 1, y : 1
		 let f3[w: 1, a : 1 -* 1, b : 1] = send a<b, self>
		 prc[x] : 1 * 1= f3(aa, bb)
		 prc[y] : 1 = f3(self, x, y)`,
		`assuming aa : 1 -* 1, bb : 1, x : 1 , y : 1
		 let f3[w: 1, a : 1 -* 1, b : 1] = send a<b, self>
		 prc[x] : 1 = f3(aa, bb)
		 prc[y] : 1 = f3(self, x, y)`,
		`assuming aa : 1 -* 1, bb : 1, x : 1
		 let f3[w: 1, a : 1 -* 1, b : 1] = send a<b, self>
		 prc[x] : 1 = f3(aa, bb)
		 prc[y] : 1 = f3(self, aa, x)`,
		`assuming aa : 1 -* 1, bb : 1, x : 1 -* 1, y : 1
		 let f3[w: 1, a : 1 -* 1, b : 1] = send a<b, self>
		 prc[x] : 1 = f3(aa, bb)
		 prc[y] : 1 * 1 = f3(self, aa, bb)`,
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectCut(t *testing.T) {

	cases := []string{
		// Cut
		`prc[pid1] : 1 = x : 1 <- new ( close x ); wait x; close self
		 prc[pid2] : 1 = x : 1 <- new ( close self ); wait x; close self`,
		// Cut [call]
		`let f() : 1 = close self
		 prc[pid1] : 1 = x : 1 <- new f(); wait x; close self`,
		`let f2[w : 1] = close w
		 prc[pid2] : 1 = x : 1 <- new f2(); wait x; close self`,
		`let f(p : &{labelok : 1}) : 1 = p.labelok<self>
		 prc[pid1] : 1 = x <- new (f(pid2)); drop x; close self
		 prc[pid2] : &{labelok : 1} = case self (labelok<b> => close b)`,
		`prc[pid1] : 1 = xy : +{labelok : 1} <- new ( self.labelok<ff> );
					case xy (labelok<b> => print b; wait b; close self)
		 prc[ff] : 1 = close self`,
		`let f() : 1 = close self
		 prc[pid1] : 1 = x : 1 <- new f(x); wait x; close self`,
		`type A = &{label : 1}
		 type B = 1 -* 1
		 let f(a : A, b : B) : A * B = send self<a, b>
		 assuming a : A, b : B
		 prc[pid1] : 1 = x <- new f(a, b); <u, v> <- recv x;  drop u; drop v; close self`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckIncorrectCut(t *testing.T) {

	cases := []string{
		// Cut [Call]
		`let f(p : &{labelok : 1}) : 1 = p.labelok<self>
		 assuming pid2 : +{labelok : 1}
		 prc[pid1] : 1 = x <- new (f(pid2)); drop x; close self
		 prc[pid2] : &{labelok : 1} = case self (labelok<b> => close b)`,

		// Cut [axiomatic version]
		`assuming pid2 : &{labelok : 1}
		 prc[pid1] : 1 = x <- new ( pid2.labelok<x> ); drop x; close self
		 prc[pid2] : &{labelok : 1} = case self (labelok<b> => close b) `,
		`assuming ff : 1
		 prc[pid1] : 1 = x : 1 <- new ( self.labelok<ff> ); wait x; close self`,
		`assuming ff : 1
		 prc[pid2] : 1 = x : 1 <- new ( x.labelok<ff> ); wait x; close self`,
		`assuming y : 1
		 let f() : 1 = close self
		 prc[pid1] : 1 = x : 1 <- new f(y); wait x; close self`,
		`let f() : 1 = close self
		 prc[pid1] : 1 = x : 1 <- new f_other(); wait x; close self`,
	}

	runThroughTypechecker(t, cases, false)
}

func TestTypecheckCorrectSplit(t *testing.T) {

	cases := []string{
		// Split
		`prc[pid0] : 1 = <u, v> <- split x; wait u; wait v; close self
		 prc[x] : 1 = close self`,
		`type A = 1 * 1
		 prc[pid1] : 1 = <a, b> <- split +pid2;
		 				 <a2, b2> <- recv a;
						 <a3, b3> <- recv b;
						 wait a2;
						 wait b2;
						 wait a3;
						 wait b3;
						 close self
	     assuming pid3 : 1, pid4 : 1
		 prc[pid2] : A = send self<pid3, pid4>`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckIncorrectSplit(t *testing.T) {

	cases := []string{
		// Split
		`assuming x : 1
		 prc[pid0] : 1 = <u, u> <- split x; wait u; close self
		 prc[x] : 1 = close self`,
		`let f(x: 1, y : 1) : 1 = <u, x> <- split y; wait x; close self`,
		`let f(x: 1, y : 1) : 1 = <u1, u2> <- split y; wait x; close self`,
	}

	runThroughTypechecker(t, cases, false)
}

// Preliminary Checks
func TestPreliminaryFunctionDefChecksCorrect(t *testing.T) {

	cases := []string{
		`type A = 1`,
		`type A = 1
		 type B = A -* 1`,
		`type C = D
		 type D = E
		 type E = 1`,
		`let f() : 1 = close self
		 let f2() : 1 = close self`,
		`let f(x : 1) : 1 = wait x; close self`,
		`let f(x : 1, y : 1) : 1 = wait x; drop y; close self`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestPreliminaryFunctionDefChecksIncorrect(t *testing.T) {

	cases := []string{
		// B is undefined
		`type A = B`,
		// Duplicate label A
		`type A = 1
		 type A = 1`,
		// D undefined
		`type A = 1
		 type B = A -* 1
		 type C = A -* D`,
		// Non-contractive type
		`type C = D
		 type D = E
		 type E = C`,
		// Duplicate function name
		`let f() : 1 = close self
		 let f() : 1 = close self`,
		`let f() : 1 = close self
		 let f[w : 1] = close w`,
		// No provider type
		`let f() = close self`,
		`let f(w) = close w`,
		// No parameter type
		`let f(x) : 1 = wait x; close self`,
		// Duplicate parameter names
		`let f(x : 1, x : 1) : 1 = wait x; drop x; close self`,
		// Invalid parameter type
		`type C = D
		 let f(x : C) : 1 = wait x; close self`,
	}

	runThroughTypechecker(t, cases, false)
}

func TestPreliminaryProcessesChecksCorrect(t *testing.T) {

	cases := []string{
		`type A = 1`,
		`assuming a : 1, b : 1
		 prc[x] : 1 = drop a; wait b; close self`,
		`prc[a, b] : 1 = close self`,
		`assuming x : 1, y : 1
		 prc[a, b] : 1 = close self
		 prc[c] : 1 = wait a; close self
		 prc[d] : 1 = wait b; close self
		 prc[e] : 1 = wait x; wait y; close self`,
		`assuming x : 1
		 prc[a, b] : 1 = close self
		 prc[c] : 1 = wait a; wait x; close self
		 prc[d] : 1 = wait b; close self`,
		`assuming x : 1, y : 1
		 prc[a, b] : 1 = close self
		 prc[c] : 1 = wait a; wait x; close self
		 prc[d] : 1 = wait b; drop y; close self`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestPreliminaryProcessesChecksIncorrect(t *testing.T) {

	cases := []string{
		// Duplicate name assumptions
		`assuming a : 1, b : 1, a : 1`,
		`assuming a : 1, b : 1
		 assuming a : 1`,
		// Assumptions without types
		`assuming a : 1, b`,
		// Undefined type B
		`assuming a : 1, b : 1 * B`,
		// Missing type of provider
		`prc[a] = close self`,
		// Duplicate provider names
		`prc[a, b, a] : 1 = close self`,
		`prc[a, b] : 1 = close self
		 prc[a] : 1 = close self`,
		`assuming x : 1, y : 1
		 prc[b] : 1 = close aaa
		 prc[c] : 1 = wait aaa; close self`,
		`assuming x : 1
		 prc[a, b] : 1 = close aaa
		 prc[c] : 1 = wait a; wait x; close self`,
		// assumed name used elsewhere
		`assuming x : 1
		 prc[a, b] : 1 = close self
		 prc[c] : 1 = wait a; wait x; close self
		 prc[d] : 1 = wait b; close self
		 prc[e] : 1 = wait x; close self`,
		// process name used elsewhere
		`assuming x : 1
		 prc[a, b] : 1 = close self
		 prc[c] : 1 = wait a; wait x; close self
		 prc[d] : 1 = wait b; close self
		 prc[e] : 1 = wait b; close self`,
		// remaining unused assumed names
		`assuming x : 1, y : 1
		 prc[a, b] : 1 = close self
		 prc[c] : 1 = wait a; wait x; close self
		 prc[d] : 1 = wait b; close self`,
	}

	runThroughTypechecker(t, cases, false)
}

func TestExecCorrect(t *testing.T) {

	cases := []string{
		`type A = 1

		let f() : A = x : A <- new close x; 
						  wait x; 
						  close self
		
		exec f()`,
		`type A = 1

		let f[w : A] = x : A <- new close x; 
						  wait x; 
						  close w
		
		exec f()
		exec f()`,
	}

	runThroughTypechecker(t, cases, true)
}

// Polarities

func TestTypecheckCorrectPolarity(t *testing.T) {

	cases := []string{
		// ID
		"let f1(x : 1 * 1) : 1 * 1 = fwd self x",
		"let f1(x : 1 * 1) : 1 * 1 = fwd self +x",
		"let f1(x : 1 -* 1) : 1 -* 1 = fwd self x",
		"let f1(x : 1 -* 1) : 1 -* 1 = fwd self -x",
		`assuming x : 1 -* 1
		 prc[f] : 1 -* 1 = fwd self x`,
		`assuming x : 1 -* 1
		 prc[f] : 1 -* 1 = fwd self -x`,
		`prc[a] : 1 = close self
		 prc[b] : 1 = fwd self a
		 prc[c] : 1 = wait b; close self`,
		`prc[a] : 1 = close self
		 prc[b] : 1 = fwd self +a
		 prc[c] : 1 = wait b; close self`,
		`type A = &{labelok : 1}
		prc[x] : 1 = y.labelok<self>
		prc[y] : A = fwd self -z
		prc[z] : A = case self ( labelok<b> => close b )`,
		`type A = &{labelok : 1}
		prc[x] : 1 = y.labelok<self>
		prc[y] : A = zz : A <- new fwd self -z;
					 fwd self -zz
		prc[z] : A = case self ( labelok<b> => close b )`,
		// New
		`let f1(x : 1) : +{label : 1} = self.label<x>
		 prc[y] : 1 = m <- new f1(z); case m (label<zz> => wait zz; close self)
		 prc[z] : 1 = close self`,
		`let f1(x : 1) : +{label : 1} = self.label<x>
		 prc[y] : 1 = m <- new f1(z); case m (label<zz> => wait zz; close self)
		 prc[z] : 1 = close self`,
		`type A = +{label : 1}
		 assuming z : 1
		 prc[y] : 1 = m : A <- new self.label<z>; case m (label<zz> => wait zz; close self)`,
		`type A = +{label : 1}
		 assuming z : 1
		 prc[y] : 1 = m : A <- new self.label<z>; case m (label<zz> => wait zz; close self)`,
		// Call
		`type A = +{label : 1}
		 type B = 1
		 let f1(x : B) : A = self.label<x>
		 prc[y] : A = f1(z)
		 prc[y2] : 1 = case y (label<zz> => wait zz; close self )
		 prc[z] : 1 = close self`,
		`type A = +{label : 1}
		 type B = 1
		 let f1(x : B) : A = self.label<x>
		 prc[y] : A = f1(z)
		 prc[y2] : 1 = case y (label<zz> => wait zz; close self )
		 prc[z] : 1 = close self`,
		`type A = &{label : 1}
		 type B = 1
		 let f1(x : A) : B = x.label<self>
		 let f2() : A = case self (label<zz> => close self )
		 prc[y] : B = f1(z)
		 prc[z] : A = f2()`,
	}

	runThroughTypechecker(t, cases, true)
}

func TestTypecheckIncorrectPolarity(t *testing.T) {
	cases := []string{
		// ID
		"let f1(x : 1 * 1) : 1 * 1 = fwd self -x",
		"let f1(x : 1 -* 1) : 1 -* 1 = fwd self +x",
		`assuming x : 1 -* 1
		 prc[f] : 1 -* 1 = fwd self +x`,
		`prc[a] : 1 = close self
		prc[b] : 1 = fwd self -a
		prc[c] : 1 = wait b; close self`,
		`type A = &{labelok : 1}
		prc[x] : 1 = y.labelok<self>
		prc[y] : A = fwd self +z
		prc[z] : A = case self ( labelok<b> => close b )`,
		// New
		`let f1(x : 1) : +{label : 1} = self.label<x>
		 prc[y] : 1 = m <- new f1(z); case -m (label<zz> => wait zz; close self )
		 prc[z] : 1 = close self`,
		`type A = +{label : 1}
		 assuming z : 1
		 prc[y] : 1 = m : A <- new self.label< -z >; case m (label<zz> => wait zz; close self )`,
		`type A = +{label : 1}
		 type B = 1
		 let f1(x : B) : A = self.label<x>
		 prc[y] : A = f1(z)
		 prc[y2] : 1 = case y (label<zz> => wait -zz; close self )
		 prc[z] : 1 = close self`,
		// Call
		`type A = &{label : 1}
		 type B = 1
		 let f1(x : A) : B = +x.label<self>
		 let f2() : A = case self (label<zz> => close self )
		 prc[y] : B = f1(z)
		 prc[z] : A = f2()`,
	}

	runThroughTypechecker(t, cases, false)
}

func TestFunctionDefinitionModes(t *testing.T) {

	inputProgram :=
		`
		// type C = linear 1 * 1
		// type Q = affine 1 -* (1 * +{ab : 1, cd : 1})
		// type R =  1 * (1 * linear /\ affine +{ab : 1, cd : 1})
		// type S =  1 -* (1 * affine /\ linear &{ab : 1, cd : 1})
		
		
		type A = +{l1 : 1}
		type B = 1 -* 1
		let f1(a : A, b : B) : A * B = send self<a, b>
		
		type C = 1 -* 1
		let f2(a : linear +{l1 : 1}, b : C) : linear +{l1 : 1} * C = 
				send self<a, b>
		 `

	cases := []struct {
		input1 string
		input2 []string
	}{
		{"[unr]A [unr]* [unr]B", []string{"[unr]A", "[unr]B"}},                       // f1
		{"lin+{l1 : [lin]1} [lin]* [lin]C", []string{"lin+{l1 : [lin]1}", "[unr]C"}}, // f2
	}

	processes, assumedFreeNames, globalEnv, err := parser.ParseString(inputProgram)

	if err != nil {
		t.Errorf("compilation error in program: %s\n", err.Error())
	}

	err = process.Typecheck(processes, assumedFreeNames, globalEnv)

	if err != nil {
		t.Errorf("expected no type errors in program, but found %s\n", err.Error())
	}

	functionDefs := *globalEnv.FunctionDefinitions

	if len(functionDefs) != len(cases) {
		t.Errorf("number of cases do not match with the type definitions\n")
	}

	for i, c := range functionDefs {
		if c.Type.StringWithModality() != cases[i].input1 {
			t.Errorf("error in case #%d: Got %s, but expected %s\n", i, c.Type.StringWithModality(), cases[i].input1)
		}

		for j := range c.Parameters {
			if c.Parameters[j].Type.StringWithModality() != cases[i].input2[j] {
				t.Errorf("error in parameter case #%d: Got %s, but expected %s\n", i, c.Parameters[j].Type.StringWithModality(), cases[i].input2[j])
			}
		}
	}
}

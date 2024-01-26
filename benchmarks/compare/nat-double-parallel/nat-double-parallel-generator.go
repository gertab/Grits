package main

import (
	"bytes"
	"fmt"
	"os"
)

// To generate files:
// cd benchmarks/compare/nat-double-parallel
// go run .

const (
	fileName  = "nat-double-parallel"
	extension = ".phi"
	upto      = 50
	every     = 2
)

// Script to generate the nat-double files
func main() {
	for i := every; i <= upto; i += every {
		name := fileName + "-" + fmt.Sprint(i) + extension
		f, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		buffer := synthesize(i)
		f.Write(buffer.Bytes())
	}

	// Synthesize equivalent Sax program
	name := fileName + ".sax"
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buffer := synthesizeSax(every, upto)
	f.Write(buffer.Bytes())
}

func synthesize(parallelThreads int) bytes.Buffer {
	var buffer bytes.Buffer

	const firstPart = `
///////// Initiate execution /////////

prc[a] : listNat = runTests()
//prc[b] : 1 = printList(a)

///////// Natural number type and function definitions /////////

type nat = +{zero : 1, succ : nat}
type listNat = +{cons : nat * listNat, nil : 1}
`

	buffer.WriteString(firstPart)

	const commonFunctions = `
let double(x : nat) : nat =
case x (
		zero<x'> => self.zero<x'>
	| succ<x'> => h <- new double(x');
					d : nat <- new d.succ<h>;
					self.succ<d>
)

// Creates an empty list
let emptyList() : listNat =
  term : 1 <- new close self;
  self.nil<term>

// Appends an element to an existing list
let appendElement(elem : nat, K : listNat) : listNat =
  next : nat * listNat <- new (send self<elem, K>);
  self.cons<next>
`
	buffer.WriteString(commonFunctions)

	const useDoublingFunc = `
// Doubles a number 5 times (i.e. produces 2^5). It needs to receive a 'start' message to initiate execution
let performSomeDoubling() : &{start : nat} =
	case self (
		start<result> =>
		a <- new nat1();
		d1 <- new double(a);
		d2 <- new double(d1);
		d3 <- new double(d2);
		d4 <- new double(d3);
		d5 <- new double(d4);
		fwd result d5
	)
	`
	buffer.WriteString(useDoublingFunc)

	const testPart1 = `
// Creates the testing environment
let runTests() : listNat =
    // Spawn all parallel instances
`
	buffer.WriteString(testPart1)

	for i := 1; i <= parallelThreads; i += 1 {
		testPart2 := fmt.Sprintf("    instance%d <- new performSomeDoubling();\n", i)
		buffer.WriteString(testPart2)
	}

	const testPart3 = `
    // Ask all instances to start
`
	buffer.WriteString(testPart3)

	for i := 1; i <= parallelThreads; i += 1 {
		testPart4 := fmt.Sprintf("    instance%dresult : nat <- new instance%d.start<instance%dresult>;\n", i, i, i)
		buffer.WriteString(testPart4)
	}

	const testPart5 = `
    // Collect all results in one list
    list0  <- new emptyList();
`
	buffer.WriteString(testPart5)

	for i := 1; i <= parallelThreads; i += 1 {
		testPart6 := fmt.Sprintf("    list%d <- new appendElement(instance%dresult, list%d);\n", i, i, i-1)
		buffer.WriteString(testPart6)
	}

	const testPart7 = `
    // Forward the list result
`
	buffer.WriteString(testPart7)

	testPart8 := fmt.Sprintf("    fwd self list%d\n", parallelThreads)
	buffer.WriteString(testPart8)

	const remainingFunctions = `
///////// Natural numbers constants /////////

// 1 : S(0)
let nat1() : nat =
  t : 1 <- new close t;
  z  : nat <- new z.zero<t>;
  s0 : nat <- new s0.succ<z>;
  fwd self s0

///////// Printing Helpers /////////

let consumeNat(n : nat) : 1 =
        case n ( zero<c> => print zero; wait c; close self
               | succ<c> => print succ; consumeNat(c))

let printNat(n : nat) : 1 =
          y <- new consumeNat(n);
          wait y;
          close self

let consumeList(l : listNat) : 1 =
        case l ( cons<c> => print _cons_;
                            <b, L2> <- recv c;
                            bConsume <- new consumeNat(b);
                            wait bConsume;
                            consumeList(L2)
               | nil<c>  => print _nil_;
                            wait c;
                            close self)

let printList(l : listNat) : 1 =
          y <- new consumeList(l);
          wait y;
          close self
`
	buffer.WriteString(remainingFunctions)

	return buffer
}

// Creates equivalent Sax program
func synthesizeSax(every, upto int) bytes.Buffer {
	var buffer bytes.Buffer

	const firstPart = `
% ./sax -q nat-double-parallel.sax

type nat = +{'zero : 1, 'succ : nat}
type listNat = +{'cons : nat * listNat, 'nil : 1}

proc double (r : nat) (x : nat) =
  recv x ( 'zero() => send r 'zero()
         | 'succ(x') => 
            x'' <- call double x'' x';
            send r 'succ('succ(x''))
  )

% Doubles a number 5 times. It needs to receive a 'start' message to initiate execution
proc performSomeDoubling (r : &{'start : nat}) = 
    recv r (
      'start(result) => 
        x : nat <- send x 'succ('zero()) ;
        d1 <- call double d1 x;
        d2 <- call double d2 d1;
        d3 <- call double d3 d2;
        d4 <- call double d4 d3;
        d5 <- call double d5 d4;
        fwd result d5
    )

% Creates an empty list
proc emptyList (l : listNat) = 
  send l 'nil()

% Appends an element to an existing list
proc appendElement (l : listNat) (elem : nat) (K : listNat) =
  send l 'cons(elem, K)


% Creates the testing environment
`

	buffer.WriteString(firstPart)

	for i := every; i <= upto; i += every {
		b := saxFunc(i)
		buffer.Write(b.Bytes())
	}

	for i := every; i <= upto; i += every {
		f1 := fmt.Sprintf("exec runTests%d\n", i)
		buffer.WriteString(f1)
	}

	return buffer
}

/*
Produces sax programs, similar to:

proc runTests2 (result : listNat) =
    % Spawn all parallel instances
    instance1 <- call performSomeDoubling instance1;
    instance2 <- call performSomeDoubling instance2;

    % Ask all instances to start
    instance1result : nat <- send instance1 'start(instance1result);
    instance2result : nat <- send instance2 'start(instance2result);

    % Collect all results in one list
    list  <- call emptyList list;

    list1 <- call appendElement list1 instance1result list;
    list2 <- call appendElement list2 instance2result list1;

    % Forward the list result
    fwd result list2
*/

func saxFunc(parallelThreads int) bytes.Buffer {
	var buffer bytes.Buffer

	f1 := fmt.Sprintf("proc runTests%d (result : listNat) =\n    %% Spawn all parallel instances\n", parallelThreads)
	buffer.WriteString(f1)

	for i := 1; i <= parallelThreads; i += 1 {
		f2 := fmt.Sprintf("    instance%d <- call performSomeDoubling instance%d;\n", i, i)
		buffer.WriteString(f2)
	}

	f3 := `
    %% Ask all instances to start
`
	buffer.WriteString(f3)

	for i := 1; i <= parallelThreads; i += 1 {
		f4 := fmt.Sprintf("    instance%dresult : nat <- send instance%d 'start(instance%dresult);\n", i, i, i)
		buffer.WriteString(f4)
	}

	f5 := `
    % Collect all results in one list
    list0  <- call emptyList list0;
`
	buffer.WriteString(f5)

	for i := 1; i <= parallelThreads; i += 1 {
		f6 := fmt.Sprintf("    list%d <- call appendElement list%d instance%dresult list%d;\n", i, i, i, i-1)
		buffer.WriteString(f6)
	}

	f7 := `
    %% Forward the list result
`
	buffer.WriteString(f7)

	f8 := fmt.Sprintf("    fwd result list%d\n\n", parallelThreads)
	buffer.WriteString(f8)

	return buffer
}

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
	count     = 16
)

// Script to generate the nat-double files
func main() {
	for i := 1; i <= count; i++ {
		name := fileName + "-" + fmt.Sprint(i) + extension
		f, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		buffer := repeat(i)
		f.Write(buffer.Bytes())
	}

}

func repeat(n int) bytes.Buffer {
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
// Doubles a number X times. It needs to receive a 'start' message to initiate execution
let performSomeDoubling() : &{start : nat} =
	case self (
		start<result> =>
		// print starting;
		a <- new nat1();
		d1 <- new double(a);
		d2 <- new double(d1);
		d3 <- new double(d2);
		d4 <- new double(d3);
		d5 <- new double(d4);
		d6 <- new double(d5);
		d7 <- new double(d6);
		d8 <- new double(d7);
		fwd result d8
	)
	`
	// for i := 1; i <= n; i += 1 {
	// 	processPart2 := fmt.Sprintf("    d%d <- new double(d%d);\n", i, i-1)
	// 	buffer.WriteString(processPart2)
	// }

	// processPart3 := fmt.Sprintf("    fwd result d%d\n", n)
	// buffer.WriteString(processPart3)
	buffer.WriteString(useDoublingFunc)

	const testPart1 = `
// Creates the testing environment
let runTests() : listNat =
    // Spawn all parallel instances
`
	buffer.WriteString(testPart1)

	for i := 1; i <= n; i += 1 {
		testPart2 := fmt.Sprintf("    instance%d <- new performSomeDoubling();\n", i)
		buffer.WriteString(testPart2)
	}

	const testPart3 = `
    // Ask all instances to start
`
	buffer.WriteString(testPart3)

	for i := 1; i <= n; i += 1 {
		testPart4 := fmt.Sprintf("    instance%dresult : nat <- new instance%d.start<instance%dresult>;\n", i, i, i)
		buffer.WriteString(testPart4)
	}

	const testPart5 = `
    // Collect all results in one list
    list0  <- new emptyList();
`
	buffer.WriteString(testPart5)

	for i := 1; i <= n; i += 1 {
		testPart6 := fmt.Sprintf("    list%d <- new appendElement(instance%dresult, list%d);\n", i, i, i-1)
		buffer.WriteString(testPart6)
	}

	const testPart7 = `
    // Forward the list result
`
	buffer.WriteString(testPart7)

	testPart8 := fmt.Sprintf("    fwd self list%d\n", n)
	buffer.WriteString(testPart8)

	const remainingFunctions = `
///////// Natural numbers constants /////////

// 1 : S(0)
let nat1() : nat =
  t : 1 <- new close t;
  z  : nat <- new z.zero<t>;
  s0 : nat <- new s0.succ<z>;
  fwd self s0

// 5 : S(S(S(S(S(0)))))
let nat5() : nat =
  t : 1 <- new close t;
  z  : nat <- new z.zero<t>;
  s1 : nat <- new s1.succ<z>;
  s2 : nat <- new s2.succ<s1>;
  s3 : nat <- new s3.succ<s2>;
  s4 : nat <- new s4.succ<s3>;
  s5 : nat <- new s5.succ<s4>;
  fwd self s5

// 10 : S(S(S(S(S(S(S(S(S(S(0))))))))))
let nat10() : nat =
  t   : 1   <- new close t;
  z   : nat <- new z.zero<t>;
  s1  : nat <- new s1.succ<z>;
  s2  : nat <- new s2.succ<s1>;
  s3  : nat <- new s3.succ<s2>;
  s4  : nat <- new s4.succ<s3>;
  s5  : nat <- new s5.succ<s4>;
  s6  : nat <- new s6.succ<s5>;
  s7  : nat <- new s7.succ<s6>;
  s8  : nat <- new s8.succ<s7>;
  s9  : nat <- new s9.succ<s8>;
  s10 : nat <- new s10.succ<s9>;
  fwd self s10

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

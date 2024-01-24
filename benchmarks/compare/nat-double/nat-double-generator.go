package main

import (
	"bytes"
	"fmt"
	"os"
)

// To generate files:
// cd benchmarks/compare/nat-double
// go run .

const (
	fileName  = "nat-double"
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
type nat = +{zero : 1, succ : nat}

let double(x : nat) : nat =
    case x (
          zero<x'> => self.zero<x'>
        | succ<x'> => h <- new double(x');
                      d : nat <- new d.succ<h>;
                      self.succ<d>
    )

let plus1(y : nat) : nat = 
    case y (
          zero<x'> => x'' : nat <- new x''.zero<x'>;
                      self.succ<x''>
        | succ<x'> => x'' <- new plus1(x');
                      self.succ<x''>
    )

// 1 : S(0)
let nat1() : nat =
    t : 1 <- new close t;
    z  : nat <- new z.zero<t>;
    s0 : nat <- new s0.succ<z>; 
    fwd self s0
	
// Print result
let printNat(n : nat) : 1 = 
          y <- new consumeNat(n); 
          wait y;
          close self

let consumeNat(n : nat) : 1 = 
        case n ( zero<c> => print zero; wait c; close self
               | succ<c> => print succ; consumeNat(c))


`

	buffer.WriteString(firstPart)

	// The last part should be similar to: (where the number of calls to double is varied)
	/*
		// Initiate execution
			prc[d0] : nat = nat1()
			prc[b] : nat =
					d1 <- new double(d0);
					d2 <- new double(d1);
					fwd self d2
		    prc[c] : 1 = printNat(b)

	*/

	const processPart1 = `// Initiate execution
prc[d0] : nat = nat1()
prc[b] : nat = 
`
	buffer.WriteString(processPart1)
	// const processPart2 = `       d1 <- new double(d0);`

	for i := 1; i <= n; i += 1 {
		processPart2 := fmt.Sprintf("    d%d <- new double(d%d);\n", i, i-1)
		buffer.WriteString(processPart2)
	}

	processPart3 := fmt.Sprintf("    fwd self d%d\n", n)
	buffer.WriteString(processPart3)

	const processPart4 = `prc[c] : 1 = printNat(b)`
	buffer.WriteString(processPart4)

	return buffer
}

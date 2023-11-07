package main

import (
	"log"
	"phi/parser"
	"phi/process"
)

// const program = `
// prc[pid1]: case a
//
//	( label1<b> => wait a; close self
//	| label2<b> => close self )
//
// prc[a]: self.label1<self>
//
//	`

const program22 = `
type Receive = 1bc
type Label = label
type Unit = 1
type Select = +{a : b}
type Select2 = +{a : b, c : d}
type Branch = &{a : b}
type Branch2 = &{a : b, c : d}
type Send = a * b
type Receive = c -o b
type Brack = (a)
type Complex = +{a : (x -o &{a : f * g}), c : d}

`
const program = `
type Unit = 1
type Select = +{a : b}

let func3(next_pid : B) : A = send self< next_pid, self>
let func1(next_pid : a * b) : a * b = send self< next_pid, self>

let func2(next_pid : s) = send self< next_pid, self>

prc[pid1] = <a, b> <- recv pid2; wait a; close self
prc[pid2] = +fwd self pid3
prc[pid3] = +fwd self pid4
prc[pid4] = +func1(pid5)
prc[pid5] = close self
`

const program1 = `
prc[a]: f.label2<self>
prc[f]: case self
		( label1<b> => close self
		| label2<b> => close self )
`

const program2 = `
prc[a]: <a, b> <- recv c; close self
prc[c]: send self<d, self>
`

const program3 = `
let
	false() = send self.false<self>
in
  prc[pid1]: x <- -new (send pid2 <a, x>); close self
  prc[pid2]: -fwd self pid3
  prc[pid3]: <x, y> <- recv self; close y
end		
`

const program4 = `
let
	false() = self.false<self>
	true() = self.true<self>
	neg(a) = case a ( true<b> => self.false<self>
					| false<b> => self.true<self> )
in
	prc[pid0]: +true()
    prc[pid1]: +neg(pid0)

	prc[result]: case pid1 ( true<b> => wait res_true; close self
						  | false<b> => wait res_false; close self )

   	prc[res_true]: close self
   	prc[res_false]: close self
end
  `
const program5 = `
let false(): A = self.false<self>
let true(): B = self.true<self>
let neg(a): C = case a ( true<b> => self.false<self>
					| false<b> => self.true<self> )
prc[pid0]: D = +true()
prc[pid1]: E = +neg(pid0)

prc[result]: case pid1 ( true<b> => wait res_true; close self
						| false<b> => wait res_false; close self )

prc[res_true]: close self
prc[res_false]: close self
  `

// const program = `
// prc[pid1]: <a, b> <- recv pid2; close self
// prc[pid2]: send self<pid3, self>
// 	`

// const program = `
// prc[pid0]: <x, y> <- +split pid2; <a, b> <- recv x; <c, d> <- recv y; close d
// prc[pid2]: send self <xx, self>
// prc[xx]: close self
//     `

// const program = ` 	/* FWD + RCV rule  -- ok with the original scenario */
// 	let
// 	in
// 	prc[pid1]: send pid2<pidother, self>
//  	prc[pid2]: -fwd self pid3
//  	prc[pid3]: -fwd self pid4
// 	prc[pid4]: <a, b> <- recv self; close a
// 	end`

// const program = ` 	/* FWD + RCV rule  -- ok with the original scenario */
// 	let
// 	in
// prc[pid1]: send pid2<pid5, self>
// prc[pid2]: -fwd self pid3
// prc[pid3]: -fwd self pid4
// prc[pid4]: <a, b> <- recv self; close a
// 	end`

// const program = ` 	/* FWD + SND rule -- the problematic scenario*/
// 	let
// 	in
// 	prc[pid1]: <a, b> <- recv pid2; close a
// 	prc[pid2]: +fwd self pid3
// 	// prc[pid3]: +fwd self pid4
// 	prc[pid3]: send self<pid5, self>
// 	end`

// const program = ` 	/* CLS rule*/
// 	let
// 	in
// 	prc[pid1]: wait pid2; close a
// 	prc[pid2]: close self
// 	end`

// const program = ` 	/* CLS + FWD rule - problematic*/
// 	let
// 	in
// 	prc[pid1]: wait pid2; close a
//  	prc[pid2]: +fwd self pid3
// 	prc[pid3]: close self
// 	end`

// const program = ` /* CLS rule */
// 		prc[pid1]: wait pid2; close a
// 		prc[pid2]: close self
// `

// const program = ` /* CLS rule */
// prc[pid1]: <a, b> <- recv pid2; wait a; close self
// prc[pid2]: send self<pid3, self>
// prc[pid3]: close self
// `

func main() {
	// Execute from file
	// processes, err := parser.ParseFile("cmd/examples/ex1.txt")

	// Execute directly from string
	processes, globalEnv, err := parser.ParseString(program)
	if err != nil {
		log.Fatal(err)
		// fmt.Println(err.Error())
		return
	}

	re := &process.RuntimeEnvironment{
		GlobalEnvironment: globalEnv,
		Debug:             true,
		Color:             true,
		LogLevels: []process.LogLevel{
			process.LOGINFO,
			process.LOGPROCESSING,
			process.LOGRULE,
			process.LOGRULEDETAILS,
			process.LOGMONITOR,
		},
		ExecutionVersion: process.NORMAL_ASYNC,
	}

	process.Typecheck(processes, globalEnv)

	process.InitializeProcesses(processes, nil, nil, re)

	// Run via API
	// setupAPI()
}

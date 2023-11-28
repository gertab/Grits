package main

import (
	"flag"
	"fmt"
	"log"
	"phi/parser"
	"phi/process"
	"phi/webserver"
)

/* ignore sample programs -- used for development*/

const program = `

type A = 1 -* (1 -* (1 * 1))
prc[x1] : 1 -* (1 * 1) = send z<yy, self>
prc[x2] : 1 * 1 = send x1<xx, self>
prc[z] : A = <x, yyyyyyy> <- recv self; 
			 <xx, yyyyyyy> <- recv yyyyyyy; 
			 send yyyyyyy<x, xx>
prc[xx] : 1 = close self
prc[yy] : 1 = close self

prc[final] : 1 = <g1, g2> <- recv x2;
				 drop g1;
				 drop g2;
				 close self


// assuming pid3 : 1, pid4 : 1

// prc[pid1] : 1 = <pid2_first, pid2_second> <- split pid2; /* split gets its polarity from the types */
// 				k : 1 <- new send pid2_first<pid3, self>;
// 				wait k;
// 				send pid2_second<pid4, self>
// prc[pid2] : 1 -* 1 = <a, b> <- recv self; 
// 					 drop a; 
// 					 close self

// // Positive fwd
// type A = +{label1 : B}
// type B = 1
// prc[y] : 1 = case ff (label1<cont> => print cont; wait cont; close self)
// prc[ff] : A = +fwd self z
// prc[z] : A = self.label1<x>
// prc[x] : B = close self

// // Positive fwd
// type A = +{label1 : B}
// type B = 1
// prc[y1] : 1 = case z1 (label1<cont> => print cont; wait cont; close self)
// prc[y2] : 1 = case z2 (label1<cont> => print cont; wait cont; close self)
// prc[z1, z2] : A = z : A <- new (self.label1<x>); +fwd self z
// prc[x] : B = close self

// // Positive fwd
// type A = &{label1 : B}
// type B = 1
// prc[y1] : 1 = z1.label1<self>
// prc[y2] : 1 = z2.label1<self>
// prc[z1, z2] : A = z : A <- new (case self (label1<cont> => print cont; close self)); +fwd self z
// prc[x] : B = close self



// type A = &{label : +{next : 1}}
// let f1(x : A) : +{next : 1} = x.label<self>
// let f2(y : 1) : A = case self (label<zz> => zz.next<y> )
// prc[x] : +{next : 1} = f1(z)
// prc[z] : A = f2(y)
// prc[y] : 1 = close self
// prc[final] : 1 = case x (next<z> => print z; drop z; close self)


// type A = 1 -* (1 -* (1 * 1))
// prc[x1] : 1 -* (1 * 1) = send z<yy, self>
// prc[x2] : 1 * 1 = send x1<xx, self>
// prc[z] : A = <x, y> <- recv self; 
// 			 <xx, y> <- recv y; 
// 			 send y<x, xx>
// prc[xx] : 1 = close self
// prc[yy] : 1 = close self

// prc[final] : 1 = <g1, g2> <- recv x2;
// 			     print g1;
// 			     print g2;
// 			     drop g1;
// 			     drop g2;
// 			     close self



// type A = 1

// let f() : A = x : A <- new close x; 
// 			wait x; 
// 			close self

// exec f()






// prc[a] : 1 = close self
// prc[b] : 1 = -fwd self a 
// prc[c] : 1 = wait b; close self

// type A = &{label : 1}
// type B = 1 -* 1
// let f(a : A, b : B) : A * B = send self<a, b>
// prc[pid1] : 1 = x <- new f(a, b); 
// 				<u, v> <- recv x;  
// 				drop u; 
// 				drop v; 
// 				close self 			% a : A, b : B







// type A = 1 * 1

// prc[pid1] : 1 = 
// 		<a, b> <- +split pid2; 
// 		<a2, b2> <- recv a; 
// 		<a3, b3> <- recv b; 
// 		wait a2; 
// 		wait b2; 
// 		wait a3; 
// 		wait b3; 
// 		close self   % pid2 : A
// prc[pid2] : A = send self<pid3, pid4>	% pid3 : 1, pid4 : 1
// prc[pid3, pid4] : 1 = close self










// let f() : 1 = close self
// prc[pid1] : 1 = x : 1 <- new f(); wait x; close self

// let f2[w : 1] = close w
// prc[pid2] : 1 = x : 1 <- new f2(); wait x; close self

//////
// type A = 1
// type B = 1
// prc[pid1] : 1 = <a, b> <- +split pid2; <a2, b2> <- recv a; <a3, b3> <- recv b; close self	
// 												% pid2 : A * B
// prc[pid2] : A * B = send self<pid3, pid4>		% pid3 : A, pid4 : B
// prc[pid3] : A = close self
// prc[pid4] : B = close self



// type A = &{label1 : 1, label2 : 1, label3 : 1}
// let f2() : A = 
// 			case self (label1<a> => close a
// 					  |label2<a> => close a
// 					  |label3<a> => close a) 

// let f3(x : &{label1 : 1}) : 1 = x.label1<self>

// prc[b] : A = f2()
// prc[dd , aa] : 1 = send a<b, self>   % a : 1 -* 1, b : 1
// prc[c] : 1 = send a<b, self>   		 % a : 1 -* 1, b : 1  
`

const program_no_errors = `
let f1(a : 1, b : 1) : 1 * 1 = send self<a, b>

type A = +{l : 1, r : 1}
let f2(a : A, b : 1 * A) : A * B = send self<a, b>
`

const program_with_errors = `
let f1(a : 1) : 1 * 1 = send self<a, b>
let f2(a : 1, b : 1, c : 1) : 1 * 1 = send self<a, b>
`

// let f2(x : 1, y : 1) : 1 * 1 = send x<y, self>

// prc[pid1] = <a, b> <- +split pid2; <a2, b2> <- recv a; <a3, b3> <- recv b; close self
// prc[pid2] = send self<pid3, self>

// type C = 1 * 1
// type D = 1 -* 1

// let func3(next_pid : D) : C = send self< next_pid, self>
// let func2(next_pid : &{a : 1, c : 1}) : &{a : 1, c : 1} = send self< next_pid, self>

// prc[pid1] : &{a : 1, c : 1} = <a, b> <- recv pid2; wait a; close self

// undefined label reference
// type B = &{a : unknownlabel, c : 1}
// type E = 1 * X
// type A = +{a : (1 -* &{a : FF * 1}), c : 1}
// let func2(next_pid : B) : &{a : ssss, c : 1} = send self< next_pid, self>
// prc[pid1] : &{a : ssss, c : 1} = <a, b> <- recv pid2; wait a; close self

// contractive
// type A = B

// multiple types with the same name
// type A = 1
// type B = 1
// type A = 1

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
type Receive = c -* b
type Brack = (a)
type Complex = +{a : (x -* &{a : f * g}), c : d}
`
const program33 = `
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

const development = true

func main() {
	// Flags
	typecheck := flag.Bool("typecheck", true, "run typechecker")
	noTypecheck := flag.Bool("notypecheck", false, "skip typechecker")
	execute := flag.Bool("execute", true, "execute processes")
	noExecute := flag.Bool("noexecute", false, "do not execute processes")
	logLevel := flag.Int("verbosity", 3, "verbosity level (1 = least, 3 = most)")
	startWebserver := flag.Bool("webserver", false, "start webserver")

	// todo: add execute synchronous vs asynchronous and with polarities

	flag.Parse()

	typecheckRes := !*noTypecheck && *typecheck
	executeRes := !*noExecute && *execute

	if *logLevel < 1 {
		*logLevel = 1
	} else if *logLevel > 3 {
		*logLevel = 3
	}

	fmt.Printf("typecheck: %t, execute: %t, verbosity: %d, webserver: %t\n", typecheckRes, executeRes, *logLevel, *startWebserver)

	args := flag.Args()

	var processes []*process.Process
	var assumedFreeNames []process.Name
	var globalEnv *process.GlobalEnvironment
	var err error

	if *startWebserver {
		// Run via API
		webserver.SetupAPI()
		return
	}

	if development {
		processes, assumedFreeNames, globalEnv, err = parser.ParseString(program)
	} else {
		if len(args) < 1 {
			fmt.Println("expected name of file to be executed")
			return
		}

		processes, assumedFreeNames, globalEnv, err = parser.ParseFile(args[0])
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	globalEnv.LogLevels = generateLogLevel(*logLevel)

	if typecheckRes {
		err = process.Typecheck(processes, assumedFreeNames, globalEnv)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	if executeRes {
		re := &process.RuntimeEnvironment{
			GlobalEnvironment: globalEnv,
			Debug:             true,
			Color:             true,
			ExecutionVersion:  process.NORMAL_ASYNC,
			Typechecked:       typecheckRes,
		}

		process.InitializeProcesses(processes, nil, nil, re)
	}

}

func generateLogLevel(logLevel int) []process.LogLevel {
	switch logLevel {
	case 1:
		return []process.LogLevel{
			process.LOGINFO,
			// process.LOGPROCESSING,
			// process.LOGRULE,
			// process.LOGRULEDETAILS,
			// process.LOGMONITOR,
		}
	case 2:
		return []process.LogLevel{
			process.LOGINFO,
			// process.LOGPROCESSING,
			process.LOGRULE,
			// process.LOGRULEDETAILS,
			// process.LOGMONITOR,
		}
	case 3:
		return []process.LogLevel{
			process.LOGINFO,
			process.LOGPROCESSING,
			process.LOGRULE,
			process.LOGRULEDETAILS,
			process.LOGMONITOR,
		}
	default:
		return []process.LogLevel{
			process.LOGINFO,
			process.LOGPROCESSING,
			process.LOGRULE,
			process.LOGRULEDETAILS,
			process.LOGMONITOR,
		}
	}
}

// func executionVersion() process.Execution_Version {

// }

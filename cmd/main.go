package main

import (
	"phi/parser"
	"phi/process"
)

// const program = ` 	/* SND rule */

// 	let
// 	in
// 		prc[pid1]: send self<pid3, self>
// 		prc[pid2]: <a, b> <- recv pid1; close self
// 	end`

// const program = ` 	/* RCV rule */

// 	let
// 	in
// 		prc[pid1]: <a, b> <- recv self; close self
// 		prc[pid2]: send pid1<pid3, self>
// 	end`

const program = ` 	/* FWD + RCV rule */
	let
	in
	prc[pid1]: send pid2<pid5, self>
	prc[pid2]: fwd self pid3
	prc[pid3]: fwd self pid4
	prc[pid4]: <a, b> <- recv self; close a
	end`

// const program = ` 	/* FWD + SND rule */

// 	let
// 	in
// 	prc[pid1]: <a, b> <- recv pid2; close sel
// 	prc[pid2]: fwd self pid3
// 	prc[pid3]: fwd self pid4
// 	prc[pid4]: send self<pid5, self>
// 	end`

// const program = ` 	/* CUT + SND rule */
//
//	let
//	in
//	prc[pid1]: x <- new (<a, b> <- recv pid2; close b); close self
//	prc[pid2]: send self<pid5, self>
//	end`
//
// const program = ` 	/* CUT + inner SND + inner RCV rule */
//
//	let
//	in
//		prc[pid1]: x <- new (send x<pid5, x>); <a, b> <- recv x; close self
//		prc[pid2]: x <- new (<a, b> <- recv x; close sel); send x<pid5, self>
//	end`
// const program = ` 	/* CUT + inner blocking SND + FWD + RCV rule */
// 					let
// 					in
// 					prc[pid1]: send pid2<pid5, self>
// 					prc[pid2]: fwd self pid3
// 					prc[pid3]: x <- new (send x<pid5, x>); <a, b> <- recv self; close self
// 					end`

// const program = ` 	/* CUT + RCV rule */
// 	let
// 	in
// 	prc[pid1]: x <- new (send pid2<pid5, self>); close self
// 	prc[pid2]: <a, b> <- recv self; close sel
// 	end`

// const program2 = `let
// 				in
// 					prc[a]: send self<pid3, self>
// 					prc[b]: self.labbbel<pid3>
// 					prc[c]:  < x , y > <- recv r ; send self < pid3 , self >
// 					prc[d]: send self<pid3, self>
// 					prc[e]: <x,y> <- recv r; send self<pid3, self>
// 					prc[f]: send self<pid3, self>
// 					prc[g]: case casename (
// 							label1<payloadc> => send self<pid3, self>
// 							| label2<payloadc> => self.labbbel<pid3>
// 						)
// 					prc[h]: x <- new send self<pid3, self>; < x , y > <- recv r ; send self < pid3 , self >
// 					prc[a]: <pay_c,cont_c> <- recv from_c; close pi5
// 				end
// 				`

func main() {
	// parser.Check()

	// processes := parser.ParseFile("parser/input.test")
	// program := "send to_c<pay_c,cont_c>"
	processes := parser.ParseString(program)
	process.InitializeProcesses(processes)
}

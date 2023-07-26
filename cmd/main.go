package main

import (
	"phi/parser"
	"phi/process"
)

const program = ` 	let
					in 
						prc[pid1]: send self<pid3, self>
						prc[pid2]: <a, b> <- recv pid1; close self
					end`

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

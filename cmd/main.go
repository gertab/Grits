package main

// const program = ` /* SND rule */
// 		prc[pid1]: send self<pid3, self>
// 		prc[pid2]: <a, b> <- recv pid1; close self
// 		`

// const program = ` /* RCV rule */
// 		prc[pid1]: send pid2<pid3, self>
// 		prc[pid2]: <a, b> <- recv self; close self
// 	`

const program = ` 	/* FWD + RCV rule */
	let
	in
	prc[pid1]: send pid2<pid5, self>
	prc[pid2]: fwd self pid3
	// prc[pid3]: fwd self pid4
	prc[pid3]: <a, b> <- recv self; close a
	end`

// const program = ` 	/* FWD + SND rule */
// 	let
// 	in
// 	prc[pid1]: <a, b> <- recv pid2; close a
// 	prc[pid2]: fwd self pid3
// 	// prc[pid3]: fwd self pid4
// 	prc[pid3]: send self<pid5, self>
// 	end`

// const program = ` 	/* Simple SPLIT + SND rule (x 2) */
// 	let
// 	in
// 		prc[pid1]: <a, b> <- split pid2; <a2, b2> <- recv a; <a3, b3> <- recv b; close self
// 		prc[pid2]: send self<pid3, self>
// 	end`

// const program = ` 	/* FWD + SND rule */
// 	let
// 	in
// prc[pid1]: <a, b> <- recv pid2; close a
// prc[pid2]: fwd self pid3
// prc[pid3]: fwd self pid4
// prc[pid4]: fwd self pid5
// prc[pid5]: fwd self pid6
// prc[pid6]: fwd self pid7
// prc[pid7]: fwd self pid8
// prc[pid8]: fwd self pid9
// prc[pid9]: fwd self pid10
// prc[pid10]: send self<pid555, self>
// end`

// const program = ` 	/* FWD + RCV rule */
// 	let
// 	in
// 	prc[pid1]: send pid2<pid555, self>
// 	prc[pid2]: fwd self pid3
// 	prc[pid3]: fwd self pid4
// 	prc[pid4]: fwd self pid5
// 	prc[pid5]: fwd self pid6
// 	prc[pid6]: fwd self pid7
// 	prc[pid7]: fwd self pid8
// 	prc[pid8]: fwd self pid9
// 	prc[pid9]: fwd self pid10
// 	prc[pid10]: <a, b> <- recv self; close a
// end`

// const program = ` /* SND rule with process having multiple names */
// prc[pa, pb, pc, pd]: send self<pid0, self>
// prc[pid2]: <a, b> <- recv pa; close self
// prc[pid3]: <a, b> <- recv pb; close self
// prc[pid4]: <a, b> <- recv pc; close self
// prc[pid5]: <a, b> <- recv pd; close self
// 		`

// const program = ` /* SND rule - with a bug */
// 		prc[a]: send self<pid3, self>
// 		prc[pid2]: <a, b> <- recv a; close self
// `

// const program = ` 	/* CUT + SND rule */
// 	let
// 	in
// 	prc[pid1]: x <- new (<a, b> <- recv pid2; close b); close self
// 	prc[pid2]: send self<pid5, self>
// 	end`

// const program = ` 	/* CUT + SND rule */
// 	let
// 	in
// 		prc[pid1]: x <- new (<a, b> <- recv pid2; close b); close self
// 		prc[pid2]: send self<pid5, self>
// 	end`

// const program = ` 	/* SPLIT + CUT + SND rule */
// 	let
// 	in
// prc[pid0]: <x1, x2> <- split pid1; close self
// prc[pid1]: x <- new (<a, b> <- recv pid2; close b); close self
// prc[pid2]: send self<pid5, self>
// 	end`

// const program = ` 	/* CUT + inner SND + inner RCV rule */
// 	let
// 	in
// 		prc[pid1]: x <- new (send x<pid5, x>); <a, b> <- recv x; close self
// 		prc[pid2]: x <- new (<a, b> <- recv x; close sel); send x<pid5, self>
// 	end`

// const program = ` /* FWD for client + SND + RCV rule <<- v. cool */
// 	/*
// 	Sometimes we have:
// 	  prc[pid1[3]]: [send, client] starting RCV rule
// 	or
// 	  prc[pid0_fwd[2]]: [send, client] starting RCV rule
// 	depending on whether the FWD rule executed before or after the RCV rule.
// 	*/
// 	let
// 	in
// 	prc[pid0]: <a, b> <- recv pid0_fwd; close a
// 	prc[pid0_fwd]: fwd self pid1
// 	prc[pid1]: send pid2<pid5, self>
// 	prc[pid2]: <a, b> <- recv self; send self<a, g>
// 	end`

// const program = ` 	/* CUT + inner blocking SND + FWD + RCV rule */
// 		let
// 		in
// 		prc[pid1]: send pid2<pid5, self>
// 		prc[pid2]: fwd self pid3
// 		prc[pid3]: ff <- new (send ff<pid5, ff>); <a, b> <- recv self; close self
// 	end`

// const program = ` 	/* CUT + RCV rule */
// 	let
// 	in
// 	prc[pid1]: x <- new (send pid2<pid5, self>); close self
// 	prc[pid2]: <a, b> <- recv self; close sel
// 	end`

// const program = ` 	/* SPLIT + SND rule (x 2) */
// 	let
// 	in
// 		/*prc[pid1]: <a2, b2> <- recv pid2; close abc*/
// 		prc[pid1]: <a, b> <- split pid2; <c, d> <- split a; <a2, b2> <- recv b; <a2, b2> <- recv c; <a2, b2> <- recv d; close abc
// 		prc[pid2]: send pid3<f, self>
// 		prc[pid3]: <a, b> <- recv self; send b<_wwww, _zzzz>
// 	end`

// const program = ` /* SPLIT + CALL rule */
// 		let
// 			D1(c) =  <a, b> <- recv c; close a
// 		in
// 			prc[pid0]: <x1, x2> <- split pid1; close self
// 			prc[pid1]: D1(pid2)
// 			prc[pid2]: send self<pid3, self>
// 		end`

// const program = ` /* Call rule */
// 		let
// 			D1(c) =  <a, b> <- recv c; close a
// 		in
// 			prc[pid1]: D1(self)
// 			prc[pid2]: send pid1<pid3, self>
// 		end`

// const program = `
// 		let
// 		in
// 			prc[pid1]: x <- new (<a, b> <- recv pid2; close b); close self
// 			prc[pid2]: send self<pid5, self>
// 		end`

// const program2 = `let
//
//	in
//		prc[a]: send self<pid3, self>
//		prc[b]: self.labbbel<pid3>
//		prc[c]:  < x , y > <- recv r ; send self < pid3 , self >
//		prc[d]: send self<pid3, self>
//		prc[e]: <x,y> <- recv r; send self<pid3, self>
//		prc[f]: send self<pid3, self>
//		prc[g]: case casename (
//				label1<payloadc> => send self<pid3, self>
//				| label2<payloadc> => self.labbbel<pid3>
//			)
//		prc[h]: x <- new send self<pid3, self>; < x , y > <- recv r ; send self < pid3 , self >
//		prc[a]: <pay_c,cont_c> <- recv from_c; close pi5
//	end
//	`

// const program = `
// 	let
// 	in
// 		prc[pid1]: x <- new (<a, b> <- recv pid2; close b); close self
// 		prc[pid2]: send self<pid5, self>
// 	end`

func main() {
	// Execute from file
	// processes := parser.ParseFile("parser/input.test")
	// program := "send to_c<pay_c,cont_c>"

	// Execute directly from string
	// processes := parser.ParseString(program)
	// process.InitializeProcesses(processes, nil)

	// Run via API
	setupAPI()
}

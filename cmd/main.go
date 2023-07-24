package main

import (
	"fmt"
	"phi/parser"
	"phi/process"
)

//	func main() {
//		fmt.Println("ok")
//	}

func main() {
	// parser.Check()

	processes := parser.ParseFile("parser/input.test")

	process.InitializeProcesses(processes)

	end := process.NewClose(process.Name{Ident: "self"})

	// to_c := process.Name{Ident: "to_c"}
	pay_c := process.Name{Ident: "pay_c"}
	// cont_c := process.Name{Ident: "cont_c"}
	from_c := process.Name{Ident: "from_c"}
	from_c2 := process.Name{Ident: "from_c"}
	// end := process.NewClose(process.Name{Ident: "self"})
	// proc1 := process.NewSend(to_c, pay_c, cont_c)
	// to_c2 := process.Name{Ident: "to_c"}
	pay_c2 := process.Name{Ident: "pay_c"}
	// cont_c2 := process.Name{Ident: "cont_c"}
	// proc2 := process.NewSend(to_c2, pay_c2, cont_c2)

	proc1 := process.NewCase(from_c, []*process.BranchForm{process.NewBranch(process.Label{L: "label1"}, pay_c, end), process.NewBranch(process.Label{L: "label2"}, pay_c, end)})
	proc2 := process.NewCase(from_c2, []*process.BranchForm{process.NewBranch(process.Label{L: "label1"}, pay_c2, end), process.NewBranch(process.Label{L: "label2"}, pay_c2, end)})

	if process.EqualForm(proc1, proc2) {
		fmt.Println("Equal")

	} else {
		fmt.Println("not Equal")
	}
}

package process

import (
	"fmt"
)

// Initiates new processes
func (process *Process) Transition(re *RuntimeEnvironment) {
	re.ProcessCount++

	go TransitionLoop(process, re)
}

// Entry point before a process transitions
func TransitionLoop(process *Process, re *RuntimeEnvironment) {
	fmt.Printf("Process transitioning: %s\n", process.String())
	process.Body.Transition(process, re)
}

// Transition according to the present body form (e.g. send, receive, ...)
func (f *SendForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of send: ")
	fmt.Println(f.String())

	if f.to_c.IsSelf() {
		// SND rule (provider)
		fmt.Println("[send, provider] Send to self, proceeding with SND")
		fmt.Println(process.Provider.String())

		process.Provider.Channel <- Message{Rule: SND, Channel1: f.payload_c, Channel2: f.continuation_c}
		fmt.Println("[send, provider] finished sending ")
		// todo terminate goroutine. inform monitor
	} else {
		// RCV rule (client)
		fmt.Println("[send, client] proceeding with RCV ")
	}
}
func (f *ReceiveForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of receive: ")
	fmt.Println(f.String())

	if f.from_c.IsSelf() {
		// RCV rule (provider)
		fmt.Println("[receive, provider] receive on self, proceeding with RCV ")
	} else {
		// SND rule (client)
		fmt.Println("[receive, client] proceeding with SND ")
		fmt.Println(f.from_c.String())

		message := <-f.from_c.Channel
		close(f.from_c.Channel)

		fmt.Print("Received: message.Rule")
		fmt.Println(message.Rule)

		new_body := f.continuation_e
		new_body.Substitute(f.payload_c, message.Channel1)
		new_body.Substitute(f.continuation_c, message.Channel2)

		process.Body = new_body
		TransitionLoop(process, re)
	}
}
func (f *SelectForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of select: ")
	fmt.Println(f.String())
}
func (f *BranchForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of branch: ")
	fmt.Println(f.String())
}
func (f *CaseForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of case: ")
	fmt.Println(f.String())
}
func (f *NewForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of new: ")
	fmt.Println(f.String())
}
func (f *CloseForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of close: ")
	fmt.Println(f.String())
}

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
	re.logProcessf(LOGPROCESSING, process, "Process transitioning: %s\n", process.Body.String())
	process.Body.Transition(process, re)
}

// Transition according to the present body form (e.g. send, receive, ...)
func (f *SendForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of send: %s\n", f.String())

	if f.to_c.IsSelf() {
		select {
		case m := <-process.Provider.Channel:
			switch m.Rule {
			case FWD:
				re.logProcessf(LOGRULEDETAILS, process, "RECEIVED FWD RULE\n")
				// Close current channel and switch to new one
				close(process.Provider.Channel)
				process.Provider = m.Channel1
				TransitionLoop(process, re)

			default:
				re.logProcessf(LOGRULEDETAILS, process, "RECEIVED something else --- fix\n")
			}
		case process.Provider.Channel <- Message{Rule: SND, Channel1: f.payload_c, Channel2: f.continuation_c}:
			// SND rule (provider)
			// re.logProcess(LOGRULE, process, "[send, provider] starting SND rule")
			// re.logProcessf(LOGRULEDETAILS, process, "[send, provider] Send to self (%s), proceeding with SND\n", process.Provider.String())

			// process.Provider.Channel <- Message{Rule: SND, Channel1: f.payload_c, Channel2: f.continuation_c}
			re.logProcess(LOGRULEDETAILS, process, "[send, provider] finished sending on self")

			re.logProcess(LOGRULE, process, "[send, provider] finished SND rule")

			// todo terminate goroutine. inform monitor
			process.Terminate(re)
		}
	} else {
		// RCV rule (client)
		re.logProcess(LOGRULE, process, "[send, client] starting RCV rule")

		message := <-f.to_c.Channel
		close(f.to_c.Channel)

		re.logProcessf(LOGRULEDETAILS, process, "Received message on channel %s, containing message.Rule %d\n", f.to_c.String(), message.Rule)

		new_body := message.ContinuationBody
		new_body.Substitute(message.Channel1, f.payload_c)
		new_body.Substitute(message.Channel2, f.continuation_c)

		re.logProcess(LOGRULE, process, "[send, client] finished RCV rule")

		process.Body = new_body
		TransitionLoop(process, re)
	}
}

func (f *ReceiveForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of receive: %s\n", f.String())

	if f.from_c.IsSelf() {
		// todo check where the select is random (if both can succeed) -- not sure what happens
		select {
		case m := <-process.Provider.Channel:
			switch m.Rule {
			case FWD:
				re.logProcessf(LOGRULEDETAILS, process, "RECEIVED FWD RULE\n")
				// Close current channel and switch to new one
				close(process.Provider.Channel)
				process.Provider = m.Channel1
				TransitionLoop(process, re)

			default:
				re.logProcessf(LOGRULEDETAILS, process, "RECEIVED something else --- fix\n")
			}
		case process.Provider.Channel <- Message{Rule: RCV, ContinuationBody: f.continuation_e, Channel1: f.payload_c, Channel2: f.continuation_c}:
			// RCV rule (provider)
			// default:
			// re.logProcessf(LOGRULEDETAILS, process, "[receive, provider] Send to self (%s), proceeding with RCV\n", process.Provider.String())

			// process.Provider.Channel <- Message{Rule: RCV, ContinuationBody: f.continuation_e, Channel1: f.payload_c, Channel2: f.continuation_c}
			re.logProcess(LOGRULEDETAILS, process, "[receive, provider] finished sending on self")

			re.logProcess(LOGRULE, process, "[receive, provider] finished RCV rule")
			// todo terminate goroutine. inform monitor
			process.Terminate(re)
		}
	} else {
		// SND rule (client)
		// todo ask for controller permission
		re.logProcess(LOGRULE, process, "[receive, client] starting SND rule")
		re.logProcessf(LOGRULEDETAILS, process, "[receive, client] proceeding with SND, will receive from %s\n", f.from_c.String())

		message := <-f.from_c.Channel
		close(f.from_c.Channel)

		re.logProcessf(LOGRULEDETAILS, process, "Received message on channel %s, containing message.Rule %d\n", f.from_c.String(), message.Rule)

		new_body := f.continuation_e
		new_body.Substitute(f.payload_c, message.Channel1)
		new_body.Substitute(f.continuation_c, message.Channel2)

		re.logProcess(LOGRULE, process, "[receive, client] finished SND rule")

		process.Body = new_body
		TransitionLoop(process, re)
	}
}

func (f *ForwardForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of forward: %s\n", f.String())

	if f.to_c.IsSelf() {
		// todo check that f.to_c == process.Provider
		f.from_c.Channel <- Message{Rule: FWD, Channel1: process.Provider}

		re.logProcess(LOGRULE, process, "[forward, client] finished FWD rule")
		// todo terminate goroutine. inform monitor
		// channel providing on fwd should not be closed, however this process dies
		// process.Terminate(re)
	} else {
		re.logProcess(LOGERROR, process, "should forward on self")
		// todo panic
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

func (process *Process) Terminate(re *RuntimeEnvironment) {
	re.logProcess(LOGRULEDETAILS, process, "process terminated succesfully")
}

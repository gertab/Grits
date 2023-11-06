package process

import (
	"fmt"
	"sync/atomic"
	"time"
)

// todo:
// -> only send a rule finished when the receiving end finished (instead of sending both times)
// -> DUP or copy
// -> check whether using self is ok

// send a <...>    <...> <- recv a; ...
//  |              /|\
//  |               |
//  |               |
//  -ve (ok-ish)    +ve (problematic)
//  | (active)      | (passive)
//  |               |
//  |               |
// \|/              |
// recv self     send self <...>

// Initiates new processes [new processes are spawned here]
func (process *Process) SpawnThenTransition(re *RuntimeEnvironment) {
	if re.debug {
		// ProcessCount is atomic
		atomic.AddUint64(&re.ProcessCount, 1)
	}

	// notify monitor about new process
	re.monitor.MonitorNewProcess(process)

	go process.transitionLoop(re)
}

type Transitionable interface {
	Transition(*Process, *RuntimeEnvironment)
}

// Entry point for each process transition
// todo maybe rename to process.Transition
func (process *Process) transitionLoop(re *RuntimeEnvironment) {
	re.logProcessf(LOGPROCESSING, process, "Process transitioning [%s]: %s\n", PolarityMap[process.Body.Polarity()], process.Body.String())

	// To slow down the execution speed
	time.Sleep(re.delay)

	process.Body.(Transitionable).Transition(process, re)
}

// When a process starts transitioning, a process chooses to transition as one of these forms:
//   (a) a provider     -> tries to send the final result (on the self/provider channel)
//   (b) a client       -> retrieves any pending messages (on the self/provider channel) and consumes them
//   (c) a special form (i.e. forward/split) -> sends a control message on an external channel
//   (d) internally     -> transitions immediately without sending/receiving messages
//
// A process' control channel is checked for incoming messages. If there are any, the execution of (a-d) may be relegated for later on.

func TransitionBySending(process *Process, toChan chan Message, continuationFunc func(), sendingMessage Message, re *RuntimeEnvironment) {

	if len(process.Providers) > 1 {
		// Split process if needed
		process.performDUPrule(re)
	} else {
		// Send message and perform the remaining work defined by continuationFunc
		toChan <- sendingMessage
		continuationFunc()
	}
}

func TransitionByReceiving(process *Process, clientChan chan Message, processMessageFunc func(Message), re *RuntimeEnvironment) {
	if clientChan == nil {
		re.logProcess(LOGERROR, process, "Channel not initialized (attempting to receive on a dead channel)")
		panic("Channel not initialized (attempting to receive on a dead channel)")
	}

	if len(process.Providers) > 1 {
		// Split process if needed
		process.performDUPrule(re)
	} else {

		// Blocks until a message arrives (may be a FWD request)
		receivedMessage := <-clientChan

		// Process acting as a client by consuming a message from some channel
		if receivedMessage.Rule == FWD {
			handleNegativeForwardRequest(process, receivedMessage, re)
		} else {
			processMessageFunc(receivedMessage)
		}
	}
}

func TransitionInternally(process *Process, internalTransition func(), re *RuntimeEnvironment) {
	if len(process.Providers) > 1 {
		// Split process if needed
		process.performDUPrule(re)
	} else {
		// Immediately perform internal transition (does not affect/depend on external processes)
		internalTransition()
	}
}

func handleNegativeForwardRequest(process *Process, message Message, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "Received FWD request. Continuing as %s\n", NamesToString(message.Providers))
	// todo remove Close current channel and switch to new one

	// Notify that the process will change providers (i.e. the process.Providers will die and be replaced by message.Providers)
	process.terminateBeforeRename(process.Providers, message.Providers, re)

	// the process.Providers can no longer be used, so close them
	// todo check if they are being closed anywhere else
	closeProviders(process.Providers)

	// Change the providers to the one being forwarded to
	process.Providers = message.Providers

	process.transitionLoop(re)
}

func closeProviders(providers []Name) {
	for _, p := range providers {
		if p.Channel != nil {
			close(p.Channel)
		}
	}
}

////////////////////////////////////////////////////////////
///////////////// Transition for each form /////////////////
////////////////////////////////////////////////////////////

// Transition according to the present body form (e.g. send, receive, ...)

func (f *SendForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of send: %s\n", f.String())

	if f.to_c.IsSelf {
		// SND rule (provider, +ve)
		//
		//  <...> <- from_c self; ...
		//	 /|\
		//    |
		//    |
		// [send self <...>]

		message := Message{Rule: SND, Channel1: f.payload_c, Channel2: f.continuation_c}

		sndRule := func() {
			re.logProcess(LOGRULEDETAILS, process, "[send, provider] finished sending on self")
			// todo Although here we say the process finished executing (and died),
			// the rule SND is not guaranteed to be done, since it depends on the other side as well
			process.terminate(re)
		}

		TransitionBySending(process, process.Providers[0].Channel, sndRule, message, re)
	} else {
		// RCV rule (client, -ve)
		//
		// [send to_c <payload_c, self>]
		//    |
		//    |
		//   \|/
		// <...> <- recv self; ...

		if !f.continuation_c.IsSelf {
			// todo error
			re.logProcessf(LOGERROR, process, "[send, client] in RCV rule, the continuation channel should be self, but found %s\n", f.continuation_c.String())
			panic("Expected self but found something else")
		}

		message := Message{Rule: RCV, Channel1: f.payload_c, Channel2: process.Providers[0]}
		// Send the provider channel (self) as the continuation channel

		rcvRule := func() {
			// Message is the received message
			re.logProcess(LOGRULE, process, "[send, client] starting RCV rule")
			re.logProcessf(LOGRULEDETAILS, process, "Received message on channel %s, containing message rule RCV\n", f.to_c.String())

			// Although the process dies, its provider will be used as the client's provider
			process.renamed(process.Providers, []Name{f.to_c}, re)
		}

		TransitionBySending(process, f.to_c.Channel, rcvRule, message, re)
	}
}

func (f *ReceiveForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of receive: %s\n", f.String())

	if f.from_c.IsSelf {
		// RCV rule (provider, -ve)
		//
		// [send to_c <payload_c, self>]
		//    |
		//    |
		//   \|/
		// <...> <- recv self; ...

		rcvRule := func(message Message) {
			re.logProcess(LOGRULEDETAILS, process, "[receive, provider] finished sending on self")

			if message.Rule != RCV {
				re.logProcessf(LOGERROR, process, "expected RCV, found %s\n", RuleString[message.Rule])
				panic("expected RCV")
			}

			new_body := f.continuation_e
			new_body.Substitute(f.payload_c, message.Channel1)
			new_body.Substitute(f.continuation_c, Name{IsSelf: true})

			process.finishedRule(RCV, "[receive, provider]", "(p)", re)
			// Terminate the current provider to replace them with the one being received
			process.terminateBeforeRename(process.Providers, []Name{message.Channel2}, re)

			process.Body = new_body
			process.Providers = []Name{message.Channel2}
			// process.finishedRule(RCV, "[receive, provider]", "(p)", re)
			process.processRenamed(re)

			process.transitionLoop(re)
		}

		TransitionByReceiving(process, process.Providers[0].Channel, rcvRule, re)
	} else {
		// SND rule (client, +ve)
		//
		//  [<...> <- recv self; ...]
		//	 /|\
		//    |
		//    |
		// send to_c <...>

		sndRule := func(message Message) {
			re.logProcess(LOGRULE, process, "[receive, client] starting SND rule")
			re.logProcessf(LOGRULEDETAILS, process, "[receive, client] Received message on channel %s, containing rule: %s\n", f.from_c.String(), RuleString[message.Rule])

			if message.Rule != SND {
				re.logProcessHighlight(LOGERROR, process, "expected SND")
				panic("expected SND")
			}

			new_body := f.continuation_e
			new_body.Substitute(f.payload_c, message.Channel1)
			new_body.Substitute(f.continuation_c, message.Channel2)

			process.Body = new_body

			process.finishedRule(SND, "[receive, client]", "(c)", re)
			process.transitionLoop(re)
		}

		TransitionByReceiving(process, f.from_c.Channel, sndRule, re)
	}
}

// CUT rule (Spawn new process) - provider only
func (f *NewForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of new: %s\n", f.String())

	newRule := func() {
		// This name is indicative only (for debugging), since there shouldn't be more than one process with the same channel name
		// Although channels may have an ID, processes (i.e. goroutines) are anonymous
		newChannelIdent := f.continuation_c.Ident
		// newChannelIdent := ""

		// First create fresh channel (with fake identity of the continuation_c name) to link both processes
		newChannel := re.CreateFreshChannel(newChannelIdent)

		// Substitute reference to this new channel by the actual channel in the current process and new process
		currentProcessBody := f.continuation_e
		currentProcessBody.Substitute(f.continuation_c, newChannel)
		process.Body = currentProcessBody

		// Create structure of new process
		newProcessBody := f.body
		newProcessBody.Substitute(f.continuation_c, Name{IsSelf: true})
		newProcess := NewProcess(newProcessBody, []Name{newChannel}, LINEAR, process.FunctionDefinitions, process.Types)

		re.logProcessf(LOGRULEDETAILS, process, "[new] will create new process with channel %s\n", newChannel.String())

		// Spawn and initiate new process
		newProcess.SpawnThenTransition(re)

		process.finishedRule(CUT, "[new]", "", re)
		// re.logProcess(LOGRULE, process, "[new] finished CUT rule")
		// Continue executing current process
		process.transitionLoop(re)
	}

	TransitionInternally(process, newRule, re)
}

// CALL rule
func (f *CallForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of call: %s\n", f.String())

	callRule := func() {
		// Look up function by name and arity
		arity := len(f.parameters)
		functionCall := GetFunctionByNameArity(*process.FunctionDefinitions, f.functionName, arity)

		if functionCall == nil {
			// No function found in the FunctionDefinitions list
			re.logProcessf(LOGERROR, process, "Function %s does not exist.\n", f.String())
			panic("Wrong function call")
		}

		// Function found
		functionCallBody := functionCall.Body

		for i := range functionCall.Parameters {
			functionCallBody.Substitute(functionCall.Parameters[i], f.parameters[i])
		}

		process.Body = functionCallBody

		process.finishedRule(CALL, "[call]", "", re)

		process.transitionLoop(re)
	}

	// TransitionInternally(process, callRule, re)

	// Always perform CALL before DUP
	prioritiseCallRule := true

	if len(process.Providers) > 1 && !prioritiseCallRule {
		// Split process if needed
		process.performDUPrule(re)
	} else {
		callRule()
	}
}

func (f *SelectForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of select: %s\n", f.String())

	if f.to_c.IsSelf {
		// SEL rule (provider, +ve)
		//
		//  case to (...)
		//	 /|\
		//    |
		//    |
		// [self.label<...>]

		message := Message{Rule: SEL, Channel1: f.continuation_c, Label: f.label}

		selRule := func() {
			re.logProcess(LOGRULEDETAILS, process, "[select, provider] finished sending on self")
			process.terminate(re)
		}

		TransitionBySending(process, process.Providers[0].Channel, selRule, message, re)
	} else {
		// CSE rule (client, -ve)
		//
		// [to_c.label<...>]
		//    |
		//    |
		//   \|/
		// case self (...)

		if !f.continuation_c.IsSelf {
			// todo error
			re.logProcessf(LOGERROR, process, "[select, client] in CSE rule, the continuation channel should be self, but found %s\n", f.continuation_c.String())
			panic("Expected self but found something else")
		}

		message := Message{Rule: CSE, Channel1: process.Providers[0], Label: f.label}
		// Send the provider channel (self) as the continuation channel

		cseRule := func() {
			// Message is the received message
			re.logProcess(LOGRULE, process, "[select, client] starting CSE rule")
			re.logProcessf(LOGRULEDETAILS, process, "Received message on channel %s, containing message rule CSE\n", f.to_c.String())

			// Although the process dies, its provider will be used as the client's provider
			process.renamed(process.Providers, []Name{f.to_c}, re)
		}

		TransitionBySending(process, f.to_c.Channel, cseRule, message, re)
	}
}

func (f *BranchForm) Transition(process *Process, re *RuntimeEnvironment) {
	// Should only be referred from within a case
	fmt.Print("cannot transition on branch")
	panic("should never happen")
}

func (f *CaseForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of case: %s\n", f.String())

	if f.from_c.IsSelf {
		// CSE rule (provider, -ve)
		//
		// [to_c.label<self>]
		//    |
		//    |
		//   \|/
		// case self (...)

		cseRule := func(message Message) {
			re.logProcess(LOGRULEDETAILS, process, "[case, provider] finished receiving on self")

			if message.Rule != CSE {
				re.logProcessHighlightf(LOGERROR, process, "expected CSE, found %s\n", RuleString[message.Rule])
				panic("expected CSE")
			}

			// Match received label with the available branches
			var new_body Form
			found := false
			for _, j := range f.branches {
				if j.label.Equal(message.Label) {
					// Found a matching label
					found = true
					new_body = j.continuation_e
					new_body.Substitute(j.payload_c, message.Channel1)
					break
				}
			}

			if !found {
				re.logProcessHighlightf(LOGERROR, process, "no matching labels found for %s\n", message.Label)
				panic("no matching labels found")
			}

			process.finishedRule(CSE, "[case, provider]", "(p)", re)
			// Terminate the current provider to replace them with the one being received
			process.terminateBeforeRename(process.Providers, []Name{message.Channel2}, re)

			process.Body = new_body
			process.Providers = []Name{message.Channel1}
			// process.finishedRule(CSE, "[receive, provider]", "(p)", re)
			process.processRenamed(re)

			process.transitionLoop(re)
		}

		TransitionByReceiving(process, process.Providers[0].Channel, cseRule, re)
	} else {
		// SEL rule (client, +ve)
		//
		//  [case from_c (...)]
		//	 /|\
		//    |
		//    |
		// self.label<..>

		selRule := func(message Message) {
			re.logProcess(LOGRULE, process, "[case, client] starting SEL rule")
			re.logProcessf(LOGRULEDETAILS, process, "[case, client] received select label %s on channel %s, containing rule: %s\n", message.Label, f.from_c.String(), RuleString[message.Rule])

			if message.Rule != SEL {
				re.logProcessHighlight(LOGERROR, process, "expected SEL ")
				panic("expected SEL")
			}

			// Match received label with the available branches
			var new_body Form
			found := false
			for _, j := range f.branches {
				if j.label.Equal(message.Label) {
					// Found a matching label
					found = true
					new_body = j.continuation_e
					new_body.Substitute(j.payload_c, message.Channel1)
					break
				}
			}

			if !found {
				re.logProcessHighlightf(LOGERROR, process, "no matching labels found for %s\n", message.Label)
				panic("no matching labels found")
			}

			process.Body = new_body

			process.finishedRule(SEL, "[case, client]", "(c)", re)
			process.transitionLoop(re)
		}

		TransitionByReceiving(process, f.from_c.Channel, selRule, re)
	}
}

func (f *CloseForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of close: %s\n", f.String())

	if f.from_c.IsSelf {
		// CLS rule (provider)
		// close self

		message := Message{Rule: CLS}

		clsRule := func() {
			re.logProcess(LOGRULEDETAILS, process, "[close, provider] finished sending on self")
			// the rule CLS is not guaranteed to be done, since it depends on the other side as well
			process.terminate(re)
		}

		TransitionBySending(process, process.Providers[0].Channel, clsRule, message, re)
	}
}

func (f *WaitForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of wait: %s\n", f.String())

	if f.to_c.IsSelf {
		re.logProcess(LOGERROR, process, "Found a wait on self. Wait should only wait for other channels.")
		panic("Found a wait on self. Wait should only wait for other channels.")
	}

	clsRule := func(message Message) {
		// CLS rule (client)
		// wait to_c; ...

		re.logProcess(LOGRULE, process, "[wait, client] starting CLS rule")
		re.logProcessf(LOGRULEDETAILS, process, "[wait, client] Received message on channel %s, containing rule: %s\n", f.to_c.String(), RuleString[message.Rule])

		if message.Rule != CLS {
			re.logProcessHighlight(LOGERROR, process, "expected CLS")
			panic("expected CLS")
		}

		process.Body = f.continuation_e

		process.finishedRule(CLS, "[wait, client]", "c", re)
		process.transitionLoop(re)
	}

	TransitionByReceiving(process, f.to_c.Channel, clsRule, re)
}

// func (f *DropForm) Transition(process *Process, re *RuntimeEnvironment) {
// 	fmt.Print("transition of drop: ")
// 	fmt.Println(f.String())
// }

// Special cases: Forward and Split
func (f *ForwardForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of forward: %s\n", f.String())

	if !f.to_c.IsSelf {
		re.logProcessHighlight(LOGERROR, process, "should forward on self")
		// todo panic
		panic("should forward on self")
	}

	if f.Polarity() == NEGATIVE {
		// -ve
		// least problematic
		// ACTIVE

		message := Message{Rule: FWD, Providers: process.Providers}
		f.from_c.Channel <- message
		re.logProcessf(LOGRULE, process, "[forward, client] sent FWD request to control channel %s\n", f.from_c.String())

		// todo check if this is needed: process.finishedRule(FWD, "[forward, client]", "", re)
		process.terminateForward(re)

	} else {
		// +ve
		// problematic
		// PASSIVE: wait before acting

		// Blocks until it received a message
		message := <-f.from_c.Channel
		re.logProcessf(LOGRULE, process, "[forward, +ve] received message on %s. Will become a %s \n", f.from_c.String(), RuleString[message.Rule])

		// todo: maybe instead of recreating each process, what I can do is check how many providers the
		// forwarding process has. If it has exactly 1, then just forward the message directly.
		// If it has >1, then recreate the process -- this allows for DUP to take place.

		// Depending on the message type, recreate a corresponding process
		switch message.Rule {
		case SND:
			process.Body = NewSend(f.to_c, message.Channel1, message.Channel2)
		case CLS:
			process.Body = NewClose(f.to_c)
		case FWD:
			re.logProcessf(LOGINFO, process, "oldProviders: %s, newProviders: %s\n", NamesToString(process.Providers), NamesToString(message.Providers))
			process.Body = NewForward(f.to_c, message.Providers[0], f.Polarity())
			process.Providers = message.Providers
		case SEL:
			process.Body = NewSelect(f.to_c, message.Label, message.Channel1)
			panic("todo check if ok")
		case CST:
			process.Body = NewCast(f.to_c, message.Channel1)
			panic("todo check if ok")

			// The following are not possible: e.g. a receive does not send anything
		case RCV:
			re.logProcess(LOGERROR, process, "a positive forward should never receive RCV messages")
		case CUT:
			re.logProcess(LOGERROR, process, "a positive forward should never receive CUT messages")
		case CALL:
			re.logProcess(LOGERROR, process, "a positive forward should never receive CALL messages")
		case SPLIT:
			re.logProcess(LOGERROR, process, "a positive forward should never receive SPLIT messages")
		case DUP:
			re.logProcess(LOGERROR, process, "a positive forward should never receive DUP messages")
		default:
			re.logProcessHighlightf(LOGERROR, process, "forward should handle message %s", RuleString[message.Rule])
			panic("todo forward should handle message")
		}

		process.finishedRule(FWD, "[fwd]", "(+ve)", re)
		process.transitionLoop(re)
	}
}

func (f *SplitForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of split: %s\n", f.String())

	if f.from_c.IsSelf {
		// from_cannot be self -- only split other processes
		re.logProcessHighlight(LOGERROR, process, "should not split on self")
		// todo panic
		panic("should not split on self")
	}

	// Perform SPLIT

	// Prepare new channels
	newSplitNames := []Name{re.CreateFreshChannel(f.channel_one.Ident), re.CreateFreshChannel(f.channel_two.Ident)}

	splitRule := func() {
		// todo check that f.to_c == process.Provider
		re.logProcessf(LOGRULE, process, "[split, client] initiating split for %s into %s\n", f.from_c.String(), NamesToString(newSplitNames))

		currentProcessBody := f.continuation_e
		currentProcessBody.Substitute(f.channel_one, newSplitNames[0])
		currentProcessBody.Substitute(f.channel_two, newSplitNames[1])
		process.Body = currentProcessBody

		process.finishedRule(SPLIT, "[split, client]", "(c)", re)
		// re.logProcess(LOGRULE, process, "[split, client] finished SPLIT rule (c)")

		// Create structure of new forward process
		// todo instead of process.Body.Polarity() maybe we can set a fixed polarity
		newProcessBody := NewForward(Name{IsSelf: true}, f.from_c, process.Body.Polarity())
		newProcess := NewProcess(newProcessBody, newSplitNames, LINEAR, process.FunctionDefinitions, process.Types)
		re.logProcessf(LOGRULEDETAILS, process, "[split, client] will create new forward process providing on %s\n", NamesToString(newSplitNames))
		// Spawn and initiate new forward process
		newProcess.SpawnThenTransition(re)

		process.transitionLoop(re)
	}

	TransitionInternally(process, splitRule, re)
}

func (process *Process) performDUPrule(re *RuntimeEnvironment) {
	if len(process.Providers) == 1 {
		re.logProcessHighlight(LOGERROR, process, "Cannot duplicate this process")
		panic("Cannot duplicate this process")
	}
	// The process needs to be DUPlicated

	newProcessNames := process.Providers

	re.logProcessf(LOGRULE, process, "[DUP] Initiating DUP rule. Will split in %d processes: %s\n", len(newProcessNames), NamesToString(newProcessNames))

	processFreeNames := process.Body.FreeNames()

	// re.logProcessf(LOGRULEDETAILS, process, "[DUP] Free names: %s\n", NamesToString(processFreeNames))

	// Create an array of arrays containing fresh channels
	// E.g. if we are splitting to two channels, and have 1 free name called a, then
	// we create freshChannels containing:
	//   [][]Names{
	//       [Names]{a', a''},
	//   }
	// where a' and a'' are the new fresh channels that will be substituted in place of a
	freshChannels := make([][]Name, len(processFreeNames))

	for i := range freshChannels {
		freshChannels[i] = make([]Name, len(newProcessNames))

		for j := range newProcessNames {
			freshChannels[i][j] = re.CreateFreshChannel(processFreeNames[i].Ident)
		}
	}

	chanString := ""
	for i := range freshChannels {
		chanString += processFreeNames[i].String() + ": {"
		chanString += NamesToString(freshChannels[i])
		chanString += "}; "
	}
	re.logProcessf(LOGRULEDETAILS, process, "[DUP] NEW Free names: %s\n", chanString)

	for i := range newProcessNames {
		// Prepare process to spawn, by substituting all free names with the unique ones just created
		newDuplicatedProcessBody := CopyForm(process.Body)
		for k := range freshChannels {
			newDuplicatedProcessBody.Substitute(processFreeNames[k], freshChannels[k][i])
		}

		// Create and spawn the new processes
		// Set its provider to the channel received in the DUP request
		newDuplicatedProcess := NewProcess(newDuplicatedProcessBody, []Name{newProcessNames[i]}, process.Shape, process.FunctionDefinitions, process.Types)

		re.logProcessf(LOGRULEDETAILS, process, "[DUP] creating new process (%d): %s\n", i, newDuplicatedProcess.String())

		// Need to spawn the new duplicated processes except the first one (since it's already running in its own thread)
		if i > 0 {
			newDuplicatedProcess.SpawnThenTransition(re)
		} else {
			process = newDuplicatedProcess
		}
	}

	// Create and launch the forward processes to connect the free names (which will implicitly force a chain of further duplications)
	for i := range processFreeNames {
		// Create structure of new forward process
		newProcessBody := NewForward(Name{IsSelf: true}, processFreeNames[i], process.Body.Polarity())
		newProcess := NewProcess(newProcessBody, freshChannels[i], LINEAR, process.FunctionDefinitions, process.Types)
		re.logProcessf(LOGRULEDETAILS, process, "[DUP] will create new forward process %s\n", newProcess.String())
		// Spawn and initiate new forward process
		newProcess.SpawnThenTransition(re)
	}

	details := fmt.Sprintf("(Duplicated %s into %s)", process.Providers[0].String(), NamesToString(newProcessNames))
	process.finishedRule(DUP, "[DUP]", details, re)

	// todo remove // Current process has been duplicated, so remove the duplication requirements to continue executing its body [i.e. become interactive]
	// process.OtherProviders = []Name{}
	process.transitionLoop(re)
}

func (f *CastForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of cast: %s\n", f.String())

	if f.to_c.IsSelf {
		// CST rule (provider, +ve)
		//
		//  <...> <- shift to_c; ...
		//	 /|\
		//    |
		//    |
		// [cast self <...>]

		message := Message{Rule: CST, Channel1: f.continuation_c}

		cstRule := func() {
			re.logProcess(LOGRULEDETAILS, process, "[cast, provider] finished casting on self")
			// the rule SND is not guaranteed to be done, since it depends on the other side as well
			process.terminate(re)
		}

		TransitionBySending(process, process.Providers[0].Channel, cstRule, message, re)
	} else {
		// SHF rule (client, -ve)
		//
		// [cast to_c <self>]
		//    |
		//    |
		//   \|/
		// <...> <- shift self; ...

		if !f.continuation_c.IsSelf {
			// todo error
			re.logProcessf(LOGERROR, process, "[cast, client] in SHF rule, the continuation channel should be self, but found %s\n", f.continuation_c.String())
			panic("Expected self but found something else")
		}

		message := Message{Rule: SHF, Channel1: process.Providers[0]}
		// Send the provider channel (self) as the continuation channel

		shfRule := func() {
			// Message is the received message
			re.logProcess(LOGRULE, process, "[cast, client] starting SHF rule")
			re.logProcessf(LOGRULEDETAILS, process, "Received message on channel %s, containing message rule SHF\n", f.to_c.String())

			// Although the process dies, its provider will be used as the client's provider
			process.renamed(process.Providers, []Name{f.to_c}, re)
		}

		TransitionBySending(process, f.to_c.Channel, shfRule, message, re)
	}
}

func (f *ShiftForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of shift: %s\n", f.String())

	if f.from_c.IsSelf {
		// RCV rule (provider, -ve)
		//
		// cast to_c<self>
		//    |
		//    |
		//   \|/
		// [<...> <- shift self; ...]

		shfRule := func(message Message) {
			re.logProcess(LOGRULEDETAILS, process, "[shift, provider] finished sending on self")

			if message.Rule != SHF {
				re.logProcessHighlightf(LOGERROR, process, "expected SHF, found %s\n", RuleString[message.Rule])
				panic("expected SHF")
			}

			new_body := f.continuation_e
			new_body.Substitute(f.continuation_c, Name{IsSelf: true})

			process.finishedRule(SHF, "[shift, provider]", "(p)", re)
			// Terminate the current provider to replace them with the one being received
			process.terminateBeforeRename(process.Providers, []Name{message.Channel1}, re)

			process.Body = new_body
			process.Providers = []Name{message.Channel1}
			process.processRenamed(re)

			process.transitionLoop(re)
		}

		TransitionByReceiving(process, process.Providers[0].Channel, shfRule, re)
	} else {
		// CST rule (client, +ve)
		//
		//  [<..> <- shift from_c; ...]
		//	 /|\
		//    |
		//    |
		// cast self<..>

		cstRule := func(message Message) {
			re.logProcess(LOGRULE, process, "[receive, client] starting CST rule")
			re.logProcessf(LOGRULEDETAILS, process, "[receive, client] Received message on channel %s, containing rule: %s\n", f.from_c.String(), RuleString[message.Rule])

			if message.Rule != CST {
				re.logProcessHighlight(LOGERROR, process, "expected CST")
				panic("expected CST")
			}

			new_body := f.continuation_e
			new_body.Substitute(f.continuation_c, message.Channel1)

			process.Body = new_body

			process.finishedRule(CST, "[receive, client]", "(c)", re)
			process.transitionLoop(re)
		}

		TransitionByReceiving(process, f.from_c.Channel, cstRule, re)
	}
}

// Debug
func (f *PrintForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of print: ")
	fmt.Println(f.String())
}

// To keep the log/monitor update with the currently running processes and the transition rules
// being performed, there are the following functions:
//
//	->  finishedRule/3
//	->  finishedHalfRule/3
//	->  finishedRuleBeforeRenamed/4
//	->  processRenamed/1
//	->  terminate/1
//	->  terminateForward/1
//	->  terminateBeforeRename/21
//	->  renamed/1
func (process *Process) finishedRule(rule Rule, prefix, suffix string, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULE, process, "%s finished %s rule %s\n", prefix, RuleString[rule], suffix)

	if re.debug {
		// Update monitor
		re.monitor.MonitorRuleFinished(process, rule)
	}
}

// Used when a process will be terminated, however its provider will be used (to rename) another process
// E.g. in the case of RCV, when the 'send' client dies, its provider channels are used for the continuation of the other process
func (process *Process) finishedRuleBeforeRenamed(rule Rule, prefix, suffix string, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULE, process, "%s finished %s rule %s\n", prefix, RuleString[rule], suffix)

	if re.debug {
		// Update monitor
		re.monitor.MonitorRuleFinishedBeforeRenamed(process, rule)
	}
}

// Process did not finish executing but will be taken over
func (process *Process) processRenamed(re *RuntimeEnvironment) {
	if re.debug {
		// Update monitor
		re.monitor.MonitorProcessRenamed(process)
	}
}

// Process will terminate
func (process *Process) terminate(re *RuntimeEnvironment) {
	re.logProcess(LOGRULEDETAILS, process, "process terminated successfully")

	if re.debug {
		// Update monitor
		re.monitor.MonitorProcessTerminated(process)
	}
}

// A forward process will terminate, but its providers will be used by other processes being forwarded
func (process *Process) terminateForward(re *RuntimeEnvironment) {
	re.logProcess(LOGRULEDETAILS, process, "process will change by forwarding its provider")

	if re.debug {
		// Update monitor
		re.monitor.MonitorRuleFinished(process, FWD)
	}
}

func (process *Process) terminateBeforeRename(oldProviders, newProviders []Name, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "process renamed from %s to %s\n", NamesToString(oldProviders), NamesToString(newProviders))

	if re.debug {
		// Update monitor
		// todo change
		re.monitor.MonitorProcessForwarded(process)
	}
}

func (process *Process) renamed(oldProviders, newProviders []Name, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "process renamed from %s to %s\n", NamesToString(oldProviders), NamesToString(newProviders))

	// Although the old providers should be closed (i.e. die), the process itself does not die. It lives on using the new provider names.
}

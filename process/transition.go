package process

import (
	"fmt"
	"sync/atomic"
)

// Initiates new processes [new processes are spawned here]
func (process *Process) Transition(re *RuntimeEnvironment) {
	// ProcessCount is atomic
	atomic.AddUint64(&re.ProcessCount, 1)

	// notify monitor about new process
	re.monitor.MonitorNewProcess(process)

	go TransitionLoop(process, re)
}

// Entry point for each process transition
func TransitionLoop(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGPROCESSING, process, "Process transitioning: %s\n", process.Body.String())
	process.Body.Transition(process, re)
}

func (process *Process) finishedRule(rule Rule, prefix, suffix string, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULE, process, "%s finished %s rule %s\n", prefix, RuleString[rule], suffix)

	if re.debug {
		// Update monitor
		re.monitor.MonitorRuleFinished(process, rule)
	}
}

func (process *Process) terminate(re *RuntimeEnvironment) {
	re.logProcess(LOGRULEDETAILS, process, "process terminated successfully")

	if re.debug {
		// Update monitor
		re.monitor.MonitorProcessTerminated(process)
	}
}

// When a process starts transitioning, a process chooses to transition as one of these forms:
//   (a) a provider     -> tries to send the final result (on the self/provider channel)
//   (b) a client       -> retrieves any pending messages (on the self/provider channel) and consumes them
//   (c) a special form (i.e. forward/split) -> sends a priority message on an external channel
//   (d) internally     -> transitions immediately without sending/receiving messages
//
// A process' priority channel is checked for incoming messages. If there are any, the execution of (a-d) may be relegated for later on.

func performDUPrule(process *Process, re *RuntimeEnvironment) {
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
		newDuplicatedProcess := NewProcess(newDuplicatedProcessBody, []Name{newProcessNames[i]}, process.Shape, process.FunctionDefinitions)

		re.logProcessf(LOGRULEDETAILS, process, "[DUP] creating new process (%d): %s\n", i, newDuplicatedProcess.String())

		// Need to spawn the new duplicated processes except the first one (since it's already running in its own thread)
		if i > 0 {
			newDuplicatedProcess.Transition(re)
		} else {
			process = newDuplicatedProcess
		}
	}

	// Create and launch the forward processes to connect the free names (which will implicitly force a chain of further duplications)
	for i := range processFreeNames {
		// Create structure of new forward process
		newProcessBody := NewForward(Name{IsSelf: true}, processFreeNames[i])
		newProcess := NewProcess(newProcessBody, freshChannels[i], LINEAR, process.FunctionDefinitions)
		re.logProcessf(LOGRULEDETAILS, process, "[DUP] will create new forward process %s\n", newProcess.String())
		// Spawn and initiate new forward process
		newProcess.Transition(re)
	}

	details := fmt.Sprintf("(Duplicated %s into %s)", process.Providers[0].String(), NamesToString(newProcessNames))
	process.finishedRule(DUP, "[DUP]", details, re)

	// todo remove // Current process has been duplicated, so remove the duplication requirements to continue executing its body [i.e. become interactive]
	// process.OtherProviders = []Name{}
	TransitionLoop(process, re)
}

func TransitionAsProvider(process *Process, providerFunc func(), sendingMessage Message, re *RuntimeEnvironment) {

	if len(process.Providers) > 1 {
		// Split process if needed
		performDUPrule(process, re)
	} else {
		select {
		case pm := <-process.Providers[0].PriorityChannel:
			handlePriorityMessage(process, pm, re)
		case process.Providers[0].Channel <- sendingMessage:
			// Acting as a provider by sending a message on 'self'
			providerFunc()
		}
	}
}

func TransitionAsClient(process *Process, clientChan chan Message, clientFunc func(Message), re *RuntimeEnvironment) {
	if clientChan == nil {
		re.logProcess(LOGERROR, process, "Channel not initialized (attempting to receive on a dead channel)")
		panic("Channel not initialized (attempting to receive on a dead channel)")
	}

	if len(process.Providers) > 1 {
		// Split process if needed
		performDUPrule(process, re)
	} else {
		select {
		case pm := <-process.Providers[0].PriorityChannel:
			handlePriorityMessage(process, pm, re)
		case receivedMessage := <-clientChan:
			// Acting as a client by consuming a message from some channel
			clientFunc(receivedMessage)
		}
	}
}

// func TransitionAsSpecialForm(process *Process, priorityChannel chan PriorityMessage, specialFormFunc func(), sendingPriorityMessage PriorityMessage, re *RuntimeEnvironment) {
// 	// Certain forms (e.g. forward and split forms) send their command on the priority channel, never utilizing their main provider channel
// 	select {
// 	case pm := <-process.Provider.PriorityChannel:
// 		handlePriorityMessage(process, pm, re)
// 	case priorityChannel <- sendingPriorityMessage:
// 		specialFormFunc()
// 	}
// }

func TransitionInternally(process *Process, internalFunction func(), re *RuntimeEnvironment) {
	if len(process.Providers) > 1 {
		// Split process if needed
		performDUPrule(process, re)
	} else {
		select {
		case pm := <-process.Providers[0].PriorityChannel:
			handlePriorityMessage(process, pm, re)
		default:
			internalFunction()
		}
	}
}

func handlePriorityMessage(process *Process, pm PriorityMessage, re *RuntimeEnvironment) {
	switch pm.Action {
	case FWD_REQUEST:
		fwdHandlePriorityMessage(process, pm, re)
	// case SPLIT_DUP_FWD:
	// 	// todo probably remove
	// 	splitHandlePriorityMessage(process, pm, re)
	default:
		handleInvalidPriorityMessage(process, re)
	}
}

func fwdHandlePriorityMessage(process *Process, pm PriorityMessage, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "[Priority Msg] Received FWD request. Continuing as %s\n", pm.Channels[0].String())
	// Close current channel and switch to new one

	// process.Provider = pm.Channel1
	// Reply with the body and shape so that the other process can continue executing
	process.Providers[0].PriorityChannel <- PriorityMessage{Action: FWD_REPLY, Body: process.Body, Shape: process.Shape}

	// close(process.Provider.Channel)
	// close(process.Provider.PriorityChannel)
	// TransitionLoop(process, re)
	process.terminate(re)
}

func handleInvalidPriorityMessage(process *Process, re *RuntimeEnvironment) {
	re.logProcessHighlight(LOGRULEDETAILS, process, "RECEIVED something else --- fix\n")
	// todo panic
	panic("Received incorrect priority message")
}

////////////////////////////////////////////////////////////
///////////////// Transition for each form /////////////////
////////////////////////////////////////////////////////////
// Transition according to the present body form (e.g. send, receive, ...)

func (f *SendForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of send: %s\n", f.String())

	if f.to_c.IsSelf {
		// SND rule (provider)
		// snd self<...>

		sndRule := func() {
			// SND rule (provider)
			// re.logProcess(LOGRULE, process, "[send, provider] starting SND rule")
			// re.logProcessf(LOGRULEDETAILS, process, "[send, provider] Send to self (%s), proceeding with SND\n", process.Provider.String())

			// process.Provider.Channel <- Message{Rule: SND, Channel1: f.payload_c, Channel2: f.continuation_c}
			re.logProcess(LOGRULEDETAILS, process, "[send, provider] finished sending on self")

			process.finishedRule(SND, "[send, provider]", "(p)", re)
			// re.logProcess(LOGRULE, process, "[send, provider] finished SND rule (p)")

			process.terminate(re)
		}

		message := Message{Rule: SND, Channel1: f.payload_c, Channel2: f.continuation_c}

		TransitionAsProvider(process, sndRule, message, re)
	} else {
		// RCV rule (client)
		// snd to_c<...>

		rcvRule := func(message Message) {
			// Message is the received message

			re.logProcess(LOGRULE, process, "[send, client] starting RCV rule")
			re.logProcessf(LOGRULEDETAILS, process, "Received message on channel %s, containing message rule RCV\n", f.to_c.String())
			// message := <-f.to_c.Channel
			// close(f.to_c.Channel)
			// close(f.to_c.PriorityChannel)

			// todo check that rule matches RCV
			if message.Rule != RCV {
				re.logProcessHighlight(LOGERROR, process, "expected RCV")
				panic("expected RCV")
			}

			new_body := message.ContinuationBody
			new_body.Substitute(message.Channel1, f.payload_c)
			new_body.Substitute(message.Channel2, f.continuation_c)

			if !f.continuation_c.IsSelf {
				// todo error
				re.logProcessf(LOGERROR, process, "[send, client] in RCV rule, the continuation channel should be self, but found %s\n", f.continuation_c.String())
				panic("Expected self but found something else")
			}

			process.finishedRule(RCV, "[send, client]", "(c)", re)
			// re.logProcess(LOGRULE, process, "[send, client] finished RCV rule (c)")

			process.Body = new_body
			TransitionLoop(process, re)
		}

		TransitionAsClient(process, f.to_c.Channel, rcvRule, re)
	}
}

func (f *ReceiveForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of receive: %s\n", f.String())

	if f.from_c.IsSelf {
		// RCV rule (provider)

		rcvRule := func() {
			// default:
			// re.logProcessf(LOGRULEDETAILS, process, "[receive, provider] Send to self (%s), proceeding with RCV\n", process.Provider.String())

			// process.Provider.Channel <- Message{Rule: RCV, ContinuationBody: f.continuation_e, Channel1: f.payload_c, Channel2: f.continuation_c}
			re.logProcess(LOGRULEDETAILS, process, "[receive, provider] finished sending on self")

			process.finishedRule(RCV, "[receive, provider]", "(p)", re)
			// re.logProcess(LOGRULE, process, "[receive, provider] finished RCV rule (p)")

			process.terminate(re)
		}

		message := Message{Rule: RCV, ContinuationBody: f.continuation_e, Channel1: f.payload_c, Channel2: f.continuation_c}

		TransitionAsProvider(process, rcvRule, message, re)
	} else {
		// SND rule (client)
		// todo ask for controller permission

		sndRule := func(message Message) {
			re.logProcess(LOGRULE, process, "[receive, client] starting SND rule")
			re.logProcessf(LOGRULEDETAILS, process, "[receive, client] Received message on channel %s, containing rule: %s\n", f.from_c.String(), RuleString[message.Rule])

			if message.Rule != SND {
				re.logProcessHighlight(LOGERROR, process, "expected RCV")
				panic("expected SND")
			}

			new_body := f.continuation_e
			new_body.Substitute(f.payload_c, message.Channel1)
			new_body.Substitute(f.continuation_c, message.Channel2)

			// re.logProcess(LOGRULE, process, "[receive, client] finished SND rule (c)")
			process.finishedRule(SND, "[receive, client]", "(c)", re)

			process.Body = new_body
			TransitionLoop(process, re)
		}

		TransitionAsClient(process, f.from_c.Channel, sndRule, re)
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
		newProcess := NewProcess(newProcessBody, []Name{newChannel}, LINEAR, process.FunctionDefinitions)

		re.logProcessf(LOGRULEDETAILS, process, "[new] will create new process with channel %s\n", newChannel.String())

		// Spawn and initiate new process
		newProcess.Transition(re)

		process.finishedRule(CUT, "[new]", "", re)
		// re.logProcess(LOGRULE, process, "[new] finished CUT rule")
		// Continue executing current process
		TransitionLoop(process, re)
	}

	TransitionInternally(process, newRule, re)
}

func (f *SelectForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of select: ")
	fmt.Println(f.String())
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

		TransitionLoop(process, re)
	}

	// TransitionInternally(process, callRule, re)

	// Always perform CALL before DUP
	prioritiseCallRule := true

	if len(process.Providers) > 1 && !prioritiseCallRule {
		// Split process if needed
		performDUPrule(process, re)
	} else {
		select {
		case pm := <-process.Providers[0].PriorityChannel:
			handlePriorityMessage(process, pm, re)
		default:
			callRule()
		}
	}
}

func (f *BranchForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of branch: ")
	fmt.Println(f.String())
}
func (f *CaseForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of case: ")
	fmt.Println(f.String())
}

func (f *CloseForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of close: ")
	fmt.Println(f.String())
}

// func (f *WaitForm) Transition(process *Process, re *RuntimeEnvironment) {
// 	fmt.Print("transition of wait: ")
// 	fmt.Println(f.String())
// }

// func (f *DropForm) Transition(process *Process, re *RuntimeEnvironment) {
// 	fmt.Print("transition of drop: ")
// 	fmt.Println(f.String())
// }

// Special cases: Forward and Split [split is not a special case]
func (f *ForwardForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of forward: %s\n", f.String())
	// todo check priority messages

	if f.to_c.IsSelf {

		priorityMessage := PriorityMessage{Action: FWD_REQUEST, Channels: []Name{process.Providers[0]}}

		forwardRule := func() {
			// todo check that f.to_c == process.Provider
			// f.from_c.Channel <- Message{Rule: FWD, Channel1: process.Provider}
			re.logProcessf(LOGRULE, process, "[forward, client] sent FWD request to priority channel %s\n", f.from_c.String())

			// Get reply containing body and shape to continue executing
			pm := <-f.from_c.PriorityChannel
			// todo assert that
			// re.logProcess(LOGRULE, process, "[forward, client] finished FWD rule")
			process.finishedRule(FWD, "[forward, client]", "", re)

			// todo terminate goroutine. inform monitor
			// channel providing on fwd should not be closed, however this process dies
			// process.terminate(re)
			// process.terminate(re)

			process.Body = pm.Body
			process.Shape = pm.Shape

			// todo ensure that action is correct
			if pm.Action != FWD_REPLY {
				re.logProcessHighlight(LOGERROR, process, "expected FWD_REPLY")
				panic("expected FWD_REPLY")
			}

			TransitionLoop(process, re)
		}

		// TransitionAsSpecialForm(process, f.from_c.PriorityChannel, forwardRule, priorityMessage, re)
		select {
		case pm := <-process.Providers[0].PriorityChannel:
			// todo check if this should only happen if len(process.OtherProviders) == 0
			handlePriorityMessage(process, pm, re)
		case f.from_c.PriorityChannel <- priorityMessage:
			forwardRule()
		}
	} else {
		re.logProcessHighlight(LOGERROR, process, "should forward on self")
		// todo panic
	}
}

func (f *SplitForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of split: %s\n", f.String())

	if f.from_c.IsSelf {
		// from_cannot be self -- only split other processes
		re.logProcessHighlight(LOGERROR, process, "should forward on self")
		// todo panic
		// } else if len(process.Provider) > 1 {
		// 	// Perform DUP
		// 	performDUPrule(process, re)
	} else {
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
			newProcessBody := NewForward(Name{IsSelf: true}, f.from_c)
			newProcess := NewProcess(newProcessBody, newSplitNames, LINEAR, process.FunctionDefinitions)
			re.logProcessf(LOGRULEDETAILS, process, "[split, client] will create new forward process providing on %s\n", NamesToString(newSplitNames))
			// Spawn and initiate new forward process
			newProcess.Transition(re)

			TransitionLoop(process, re)
		}

		// priorityMessage := PriorityMessage{Action: SPLIT_DUP_FWD, Channels: newSplitNames}

		// TransitionAsSpecialForm(process, f.from_c.PriorityChannel, splitRule, priorityMessage, re)
		TransitionInternally(process, splitRule, re)

		// if len(process.OtherProviders) > 0 {
		// 	// todo first you need to check whether len(process.OtherProviders) > 0 -- if so, then the process is non-interacting and needs to be split
		// 	re.logProcessf(LOGERROR, process, "process.OtherProviders = %d\n", len(process.OtherProviders))
		// 	panic("needs to split first")
		// }

		// select {
		// case pm := <-process.Provider.PriorityChannel:
		// 	handlePriorityMessage(process, pm, re)
		// default:
		// 	splitRule()
		// }
	}
}

// Debug
func (f *PrintForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of print: ")
	fmt.Println(f.String())
}

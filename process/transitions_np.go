package process

/*
 * This NP version of the transition presents a transition semantics without the need for polarities. I.e. forwards act the same irrelevant of whether they are positive or negative.
 * This works by attaching a control channel to each name (in additional to the normal channel)
 *
 * This is the version which uses parallel channels (i.e. ControlMessage channel) to pass FWD requests.
 * Works in a synchronous setting (using unbuffered channels), but fails in an asynchronous one (buffered channels).
 */

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Initiates new processes [new processes are spawned here]
func (process *Process) SpawnThenTransitionNP(re *RuntimeEnvironment) {
	if re.debug {
		// ProcessCount is atomic
		atomic.AddUint64(&re.ProcessCount, 1)
	}

	// notify monitor about new process
	re.monitor.MonitorNewProcess(process)

	go process.transitionLoopNP(re)
}

type TransitionableNP interface {
	TransitionNP(*Process, *RuntimeEnvironment)
}

// Entry point for each process transition
func (process *Process) transitionLoopNP(re *RuntimeEnvironment) {
	re.logProcessf(LOGPROCESSING, process, "Process transitioning [%s]: %s\n", PolarityMap[process.Body.Polarity()], process.Body.String())

	// To slow down the execution speed
	time.Sleep(re.delay)

	process.Body.(TransitionableNP).TransitionNP(process, re)
}

// When a process starts transitioning, a process chooses to transition as one of these forms:
//   (a) a provider     -> tries to send the final result (on the self/provider channel)
//   (b) a client       -> retrieves any pending messages (on the self/provider channel) and consumes them
//   (c) a special form (i.e. forward/split) -> sends a control message on an external channel
//   (d) internally     -> transitions immediately without sending/receiving messages
//
// A process' control channel is checked for incoming messages. If there are any, the execution of (a-d) may be relegated for later on.

func TransitionBySendingNP(process *Process, toChan chan Message, continuationFunc func(), sendingMessage Message, re *RuntimeEnvironment) {

	if len(process.Providers) > 1 {
		// Split process if needed
		process.performDUPruleNP(re)
	} else {
		select {
		case cm := <-process.Providers[0].ControlChannel:
			handleControlMessageNP(process, cm, re)
		case toChan <- sendingMessage:
			// Sending a message to toChan
			continuationFunc()
		}
	}
}

func TransitionByReceivingNP(process *Process, clientChan chan Message, processMessageFunc func(Message), re *RuntimeEnvironment) {
	if clientChan == nil {
		re.logProcess(LOGERROR, process, "Channel not initialized (attempting to receive on a dead channel)")
		panic("Channel not initialized (attempting to receive on a dead channel)")
	}

	if len(process.Providers) > 1 {
		// Split process if needed
		process.performDUPruleNP(re)
	} else {
		select {
		case cm := <-process.Providers[0].ControlChannel:
			handleControlMessageNP(process, cm, re)
		case receivedMessage := <-clientChan:
			// Acting as a client by consuming a message from some channel
			processMessageFunc(receivedMessage)
		}
	}
}

func TransitionInternallyNP(process *Process, internalFunction func(), re *RuntimeEnvironment) {
	if len(process.Providers) > 1 {
		// Split process if needed
		process.performDUPruleNP(re)
	} else {
		select {
		case cm := <-process.Providers[0].ControlChannel:
			handleControlMessageNP(process, cm, re)
		default:
			internalFunction()
		}
	}
}

func handleControlMessageNP(process *Process, cm ControlMessage, re *RuntimeEnvironment) {
	switch cm.Action {
	case FWD_REQUEST:
		fwdhandleControlMessageNP(process, cm, re)
	default:
		handleInvalidControlMessageNP(process, re)
	}
}

func fwdhandleControlMessageNP(process *Process, cm ControlMessage, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "[Control Msg] Received FWD request. Continuing as %s\n", NamesToString(cm.Providers))
	// todo remove Close current channel and switch to new one

	// todo ensure that action is correct
	if cm.Action != FWD_REQUEST {
		re.logProcessHighlight(LOGERROR, process, "expected FWD_REQUEST")
		panic("expected FWD_REQUEST")
	}

	// Notify that the process will change providers (i.e. the process.Providers will die and be replaced by cm.Providers)
	process.terminateBeforeRename(process.Providers, cm.Providers, re)

	// the process.Providers can no longer be used, so close them
	// todo check if they are being closed anywhere else
	closeProvidersNP(process.Providers)

	// Change the providers to the one being forwarded to
	process.Providers = cm.Providers

	process.transitionLoopNP(re)
}

func handleInvalidControlMessageNP(process *Process, re *RuntimeEnvironment) {
	re.logProcessHighlight(LOGRULEDETAILS, process, "RECEIVED something else --- fix\n")
	// todo panic
	panic("Received incorrect control message")
}

func closeProvidersNP(providers []Name) {
	for _, p := range providers {
		if p.Channel != nil {
			close(p.Channel)
		}
		if p.ControlChannel != nil {
			close(p.ControlChannel)
		}
	}
}

////////////////////////////////////////////////////////////
///////////////// Transition for each form /////////////////
////////////////////////////////////////////////////////////
// Transition according to the present body form (e.g. send, receive, ...)

func (f *SendForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of send: %s\n", f.String())

	if f.to_c.IsSelf {
		// SND rule (provider)
		// snd self<...>

		message := Message{Rule: SND, Channel1: f.payload_c, Channel2: f.continuation_c}

		sndRule := func() {
			re.logProcess(LOGRULEDETAILS, process, "[send, provider] finished sending on self")
			process.terminate(re)
		}

		TransitionBySendingNP(process, process.Providers[0].Channel, sndRule, message, re)
	} else {
		// RCV rule (client)
		// snd to_c<...>

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

		TransitionBySendingNP(process, f.to_c.Channel, rcvRule, message, re)
	}
}

func (f *ReceiveForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of receive: %s\n", f.String())

	if f.from_c.IsSelf {
		// RCV rule (provider)

		rcvRule := func(message Message) {
			re.logProcess(LOGRULEDETAILS, process, "[receive, provider] finished sending on self")

			if message.Rule != RCV {
				re.logProcessHighlight(LOGERROR, process, "expected RCV")
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

			process.transitionLoopNP(re)
		}

		TransitionByReceivingNP(process, process.Providers[0].Channel, rcvRule, re)
	} else {
		// SND rule (client)

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

			process.Body = new_body

			process.finishedRule(SND, "[receive, client]", "(c)", re)
			process.transitionLoopNP(re)
		}

		TransitionByReceivingNP(process, f.from_c.Channel, sndRule, re)
	}
}

// CUT rule (Spawn new process) - provider only
func (f *NewForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
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
		newProcess.SpawnThenTransitionNP(re)

		process.finishedRule(CUT, "[new]", "", re)
		// re.logProcess(LOGRULE, process, "[new] finished CUT rule")
		// Continue executing current process
		process.transitionLoopNP(re)
	}

	TransitionInternallyNP(process, newRule, re)
}

func (f *SelectForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of select: %s\n", f.String())

	if f.to_c.IsSelf {
		// SEL rule (provider, +ve)

		message := Message{Rule: SEL, Channel1: f.continuation_c, Label: f.label}

		selRule := func() {
			re.logProcess(LOGRULEDETAILS, process, "[select, provider] finished sending on self")
			process.terminate(re)
		}

		TransitionBySendingNP(process, process.Providers[0].Channel, selRule, message, re)
	} else {
		// CSE rule (client, -ve)

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

		TransitionBySendingNP(process, f.to_c.Channel, cseRule, message, re)
	}
}

func (f *BranchForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
	// Should only be referred from within a case
	fmt.Print("cannot transition on branch")
	panic("should never happen")
}

func (f *CaseForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of case: %s\n", f.String())

	if f.from_c.IsSelf {
		// CSE rule (provider)

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

			process.transitionLoopNP(re)
		}

		TransitionByReceivingNP(process, process.Providers[0].Channel, cseRule, re)
	} else {
		// SEL rule (client, +ve)

		selRule := func(message Message) {
			re.logProcess(LOGRULE, process, "[case, client] starting SEL rule")
			re.logProcessf(LOGRULEDETAILS, process, "[case, client] received select label %s on channel %s, containing rule: %s\n", message.Label, f.from_c.String(), RuleString[message.Rule])

			if message.Rule != SEL {
				re.logProcess(LOGERROR, process, "expected SEL ")
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
			process.transitionLoopNP(re)
		}

		TransitionByReceivingNP(process, f.from_c.Channel, selRule, re)
	}
}

// CALL rule
func (f *CallForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
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

		process.transitionLoopNP(re)
	}

	// TransitionInternallyNP(process, callRule, re)

	// Always perform CALL before DUP
	prioritiseCallRule := true

	if len(process.Providers) > 1 && !prioritiseCallRule {
		// Split process if needed
		process.performDUPruleNP(re)
	} else {
		select {
		case cm := <-process.Providers[0].ControlChannel:
			handleControlMessageNP(process, cm, re)
		default:
			callRule()
		}
	}
}

func (f *CloseForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
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

		TransitionBySendingNP(process, process.Providers[0].Channel, clsRule, message, re)
	}
}

func (f *WaitForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of wait: %s\n", f.String())

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
		process.transitionLoopNP(re)
	}

	TransitionByReceivingNP(process, f.to_c.Channel, clsRule, re)
}

// func (f *DropForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
// 	fmt.Print("transition of drop: ")
// 	fmt.Println(f.String())
// }

// Special cases: Forward and Split [split is not a special case]
func (f *ForwardForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of forward: %s\n", f.String())

	if !f.to_c.IsSelf {
		re.logProcessHighlight(LOGERROR, process, "should forward on self")
		// todo panic
		panic("should forward on self")
	}

	controlMessage := ControlMessage{Action: FWD_REQUEST, Providers: process.Providers}

	forwardRule := func() {
		re.logProcessf(LOGRULE, process, "[forward, client] sent FWD request to control channel %s\n", f.from_c.String())
		process.terminateForward(re)
	}

	// TransitionAsSpecialForm(process, f.from_c.ControlChannel, forwardRule, controlMessage, re)
	select {
	case cm := <-process.Providers[0].ControlChannel:
		// todo check if this should only happen if len(process.OtherProviders) == 0
		handleControlMessageNP(process, cm, re)
	case f.from_c.ControlChannel <- controlMessage:
		forwardRule()
	}
}

func (f *SplitForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
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
		newProcessBody := NewForward(Name{IsSelf: true}, f.from_c, process.Body.Polarity())
		newProcess := NewProcess(newProcessBody, newSplitNames, LINEAR, process.FunctionDefinitions, process.Types)
		re.logProcessf(LOGRULEDETAILS, process, "[split, client] will create new forward process providing on %s\n", NamesToString(newSplitNames))
		// Spawn and initiate new forward process
		newProcess.SpawnThenTransitionNP(re)

		process.transitionLoopNP(re)
	}

	TransitionInternallyNP(process, splitRule, re)
}

func (process *Process) performDUPruleNP(re *RuntimeEnvironment) {
	if len(process.Providers) == 1 {
		re.logProcessHighlight(LOGERROR, process, "Cannot duplicate this process")
		panic("Cannot duplicate this process")
	}
	// The process needs to be DUPlicated

	newProcessNames := process.Providers

	re.logProcessf(LOGRULE, process, "[DUP] Initiating DUP rule. Will split in %d processes: %s\n", len(newProcessNames), NamesToString(newProcessNames))

	processFreeNames := process.Body.FreeNames()

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
			newDuplicatedProcess.SpawnThenTransitionNP(re)
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
		newProcess.SpawnThenTransitionNP(re)
	}

	details := fmt.Sprintf("(Duplicated %s into %s)", process.Providers[0].String(), NamesToString(newProcessNames))
	process.finishedRule(DUP, "[DUP]", details, re)

	// todo remove // Current process has been duplicated, so remove the duplication requirements to continue executing its body [i.e. become interactive]
	// process.OtherProviders = []Name{}
	process.transitionLoopNP(re)
}

func (f *CastForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of cast: %s\n", f.String())

	if f.to_c.IsSelf {
		// CST rule (provider)
		// cast self<...>

		message := Message{Rule: CST, Channel1: f.continuation_c}

		cstRule := func() {
			re.logProcess(LOGRULEDETAILS, process, "[cast, provider] finished sending on self")
			process.terminate(re)
		}

		TransitionBySendingNP(process, process.Providers[0].Channel, cstRule, message, re)
	} else {
		// SHF rule (client)
		// cast to_c<...>

		if !f.continuation_c.IsSelf {
			// todo error
			re.logProcessf(LOGERROR, process, "[send, client] in SHF rule, the continuation channel should be self, but found %s\n", f.continuation_c.String())
			panic("Expected self but found something else")
		}

		message := Message{Rule: SHF, Channel1: process.Providers[0]}
		// Send the provider channel (self) as the continuation channel

		shfRule := func() {
			// Message is the received message
			re.logProcess(LOGRULE, process, "[send, client] starting SHF rule")
			re.logProcessf(LOGRULEDETAILS, process, "Received message on channel %s, containing message rule SHF\n", f.to_c.String())

			// Although the process dies, its provider will be used as the client's provider
			process.renamed(process.Providers, []Name{f.to_c}, re)
		}

		TransitionBySendingNP(process, f.to_c.Channel, shfRule, message, re)
	}
}

func (f *ShiftForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of shift: %s\n", f.String())

	if f.from_c.IsSelf {
		// SHF rule (provider)

		shfRule := func(message Message) {
			re.logProcess(LOGRULEDETAILS, process, "[shift, provider] finished sending on self")

			if message.Rule != SHF {
				re.logProcessHighlight(LOGERROR, process, "expected SHF")
				panic("expected SHF")
			}

			new_body := f.continuation_e
			new_body.Substitute(f.continuation_c, Name{IsSelf: true})

			process.finishedRule(SHF, "[shift, provider]", "(p)", re)
			// Terminate the current provider to replace them with the one being shifted
			process.terminateBeforeRename(process.Providers, []Name{message.Channel2}, re)

			process.Body = new_body
			process.Providers = []Name{message.Channel1}
			// process.finishedRule(SHF, "[shift, provider]", "(p)", re)
			process.processRenamed(re)

			process.transitionLoopNP(re)
		}

		TransitionByReceivingNP(process, process.Providers[0].Channel, shfRule, re)
	} else {
		// CST rule (client)

		shfRule := func(message Message) {
			re.logProcess(LOGRULE, process, "[receive, client] starting CST rule")
			re.logProcessf(LOGRULEDETAILS, process, "[receive, client] Received message on channel %s, containing rule: %s\n", f.from_c.String(), RuleString[message.Rule])

			if message.Rule != CST {
				re.logProcessHighlight(LOGERROR, process, "expected RCV")
				panic("expected CST")
			}

			new_body := f.continuation_e
			new_body.Substitute(f.continuation_c, message.Channel1)

			// re.logProcess(LOGRULE, process, "[receive, client] finished CST rule (c)")

			process.Body = new_body

			process.finishedRule(CST, "[receive, client]", "(c)", re)
			process.transitionLoopNP(re)
		}

		TransitionByReceivingNP(process, f.from_c.Channel, shfRule, re)
	}
}

// Debug
func (f *PrintForm) TransitionNP(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of print: ")
	fmt.Println(f.String())
}

// // To keep the log/monitor update with the currently running processes and the transition rules
// // being performed, there are the following functions:
// //
// //	->  finishedRule/3
// //	->  finishedRuleBeforeRenamed/4
// //	->  processRenamed/1
// //	->  terminate/1
// //	->  terminateForward/1
// //	->  terminateBeforeRename/21
// //	->  renamed/1

// These can be found in transition.go

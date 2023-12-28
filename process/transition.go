package process

import (
	"fmt"
	"phi/types"
	"sync/atomic"
	"time"
)

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
	// Increment ProcessCount atomically
	atomic.AddUint64(&re.processCount, 1)

	// notify monitor about new process
	re.monitor.MonitorNewProcess(process)

	go process.transitionLoop(re)
}

// Entry point for each process transition
// todo maybe rename to process.Transition
func (process *Process) transitionLoop(re *RuntimeEnvironment) {
	re.logProcessf(LOGPROCESSING, process, "Process transitioning: %s\n", process.Body.String())

	// Send heartbeat
	re.heartbeat <- struct{}{}

	// To slow down the execution speed
	time.Sleep(re.Delay)

	process.Body.Transition(process, re)
}

// When a process starts transitioning, a process chooses to transition as one of these forms:
//   (a) a provider     -> tries to send the final result
//   (b) a client       -> retrieves any pending messages and consumes them
//   (c) a special form (i.e. forward)
//   (d) internally     -> transitions immediately without sending/receiving messages

func TransitionBySending(process *Process, toChan chan Message, continuationFunc func(), sendingMessage Message, re *RuntimeEnvironment) {

	if len(process.Providers) > 1 {
		// Split process if needed
		process.performDUPrule(re)
	} else {
		// Send message and perform the remaining work defined by continuationFunc
		// If received cancellation request, then stop
		select {
		case <-re.ctx.Done():
			return
		default:
			toChan <- sendingMessage
			continuationFunc()
		}
	}
}

func TransitionByReceiving(process *Process, clientChan chan Message, processMessageFunc func(Message), re *RuntimeEnvironment) {
	if clientChan == nil {
		re.error(process, "Channel not initialized (attempting to receive on a dead channel)")
	}

	if len(process.Providers) > 1 {
		// Split process if needed
		process.performDUPrule(re)
	} else {
		select {
		case <-re.ctx.Done():
			// Received cancellation request, then stop
			return
		case receivedMessage := <-clientChan:
			// Blocks until a message arrives (may be a FWD request)

			// Process acting as a client by consuming a message from some channel
			if receivedMessage.Rule == FWD {
				handleNegativeForwardRequest(process, receivedMessage, re)
			} else if receivedMessage.Rule == FWD_DROP {
				handleNegativeDropRequest(process, re)
			} else {
				processMessageFunc(receivedMessage)
			}
		}
	}
}

func TransitionInternally(process *Process, internalTransition func(), re *RuntimeEnvironment) {
	select {
	case <-re.ctx.Done():
		// If received cancellation request, then stop
		return
	default:
	}

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

func handleNegativeDropRequest(process *Process, re *RuntimeEnvironment) {
	re.logProcess(LOGRULEDETAILS, process, "Received FWD_DROP request.")

	if len(NamesToString(process.Body.FreeNames())) > 0 {
		re.logProcessf(LOGRULEDETAILS, process, "FWD_DROP will extend to %s.\n", NamesToString(process.Body.FreeNames()))
	}

	// Propagate the drop the the process' clients, before terminating
	for _, fn := range process.Body.FreeNames() {
		p := createDroppableForwardFromClient(process, re, fn)
		p.SpawnThenTransition(re)
	}

	process.terminate(re)
}

func closeProviders(providers []Name) {
	for _, p := range providers {
		if p.Channel != nil {
			close(p.Channel)
		}
	}
}

// Create a forward process with the to_drop flag set to true
func createDroppableForwardFromClient(process *Process, re *RuntimeEnvironment, client Name) *Process {
	clientType := types.CopyType(client.Type)
	newProcessBody := NewDroppableForward(Name{IsSelf: true, Ident: client.Ident, Type: clientType, ExplicitPolarity: client.ExplicitPolarity}, client)
	// First create fresh channel (with fake identity of the continuation_c name) to link both processes
	newChannel := re.CreateFreshChannel(client.Ident)
	newChannel.Type = types.CopyType(client.Type)
	newChannel.ExplicitPolarity = client.ExplicitPolarity
	return NewProcess(newProcessBody, []Name{newChannel}, clientType, LINEAR, process.Position)
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
			re.errorf(process, "[send, client] in RCV rule, the continuation channel should be self, but found %s\n", f.continuation_c.String())
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
		// send to_c <payload_c, self>
		//    |
		//    |
		//   \|/
		// [<...> <- recv self; ...]

		rcvRule := func(message Message) {
			re.logProcess(LOGRULEDETAILS, process, "[receive, provider] finished receiving on self")

			if message.Rule != RCV {
				re.errorf(process, "expected RCV, found %s\n", RuleString[message.Rule])
			}

			new_body := f.continuation_e
			new_body.Substitute(f.payload_c, message.Channel1)
			new_body.Substitute(f.continuation_c, NewSelf(message.Channel2.Ident))

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
				re.error(process, "expected SND")
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
		newChannel.Type = types.CopyType(f.continuation_c.Type)
		newChannel.ExplicitPolarity = f.continuation_c.ExplicitPolarity

		// Substitute reference to this new channel by the actual channel in the current process and new process
		currentProcessBody := f.continuation_e
		currentProcessBody.Substitute(f.continuation_c, newChannel)
		process.Body = currentProcessBody

		// Create structure of new process
		newProcessBody := f.body
		innerSessionType := types.CopyType(f.continuation_c.Type)
		newProcessBody.Substitute(f.continuation_c, Name{IsSelf: true, Ident: f.continuation_c.Ident, Type: innerSessionType}) // todo include polarity in name f.continuation_c.Polarity()
		newProcess := NewProcess(newProcessBody, []Name{newChannel}, innerSessionType, LINEAR, process.Position)

		re.logProcessf(LOGRULEDETAILS, process, "[new] will create new process with channel %s\n", newChannel.String())

		// Spawn and initiate new process
		newProcess.SpawnThenTransition(re)

		process.finishedRule(CUT, "[new]", "", re)

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
		functionCall := GetFunctionByNameArity(*re.GlobalEnvironment.FunctionDefinitions, f.functionName, arity)

		if functionCall == nil {
			re.errorf(process, "Function %s does not exist.\n", f.String())
		}

		// Function found. Important to copy the body, to keep the original untouched
		functionCallBody := CopyForm(functionCall.Body)

		if functionCall.UsesExplicitProvider && arity == functionCall.Arity() {
			// let f[ExplicitProvider, ...] = body    <- called using f(...)
			// No need to modify the ExplicitProvider, since it is already set as IsSelf = true (using SetProviderNameAsSelf in the function definition)

			// Substitute parameters as needed
			for i := range functionCall.Parameters {
				functionCallBody.Substitute(functionCall.Parameters[i], f.parameters[i])
			}

		} else if functionCall.UsesExplicitProvider && arity-1 == functionCall.Arity() {
			// let f[ExplicitProvider, ...] = body    <- called using f(w, ...)

			// Since the function that uses an explicit provider is called using explicit self,
			// e.g. f(self, x1, x2) or f(w, x1, x2) where w has IsSelf true,
			// then w has to be replaced by the new provider
			functionCallBody.Substitute(functionCall.ExplicitProvider, f.parameters[0])

			for i := 1; i < len(f.parameters); i++ {
				functionCallBody.Substitute(functionCall.Parameters[i-1], f.parameters[i])
			}
		} else if !functionCall.UsesExplicitProvider && arity == functionCall.Arity() {
			// let f(...) = body    <- called using f(...)

			// just substitute the parameters
			for i := range functionCall.Parameters {
				functionCallBody.Substitute(functionCall.Parameters[i], f.parameters[i])
			}
		} else if !functionCall.UsesExplicitProvider && arity-1 == functionCall.Arity() {
			// let f(...) = body    <- called using f(w, ...), then w is ignored since there is no explicit provider used

			for i := 1; i < len(f.parameters); i++ {
				functionCallBody.Substitute(functionCall.Parameters[i-1], f.parameters[i])
			}
		} else {
			// problem
			re.errorf(process, "Function %s could be not initialized.\n", f.String())
		}

		process.Body = functionCallBody

		process.finishedRule(CALL, "[call]", "", re)

		process.transitionLoop(re)
	}

	// Always perform DUP before CALL, so that if self is passed as the first parameter, then we can safely substitute the first provider
	TransitionInternally(process, callRule, re)

	// // Always perform CALL before DUP
	// prioritiseCallRule := true

	// if len(process.Providers) > 1 && !prioritiseCallRule {
	// 	// Split process if needed
	// 	process.performDUPrule(re)
	// } else {
	// 	callRule()
	// }
}

// select: to_c.label<continuation_c>
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
	} else if f.continuation_c.IsSelf {
		// CSE rule (client, -ve)
		//
		// [to_c.label<...>]
		//    |
		//    |
		//   \|/
		// case self (...)

		if !f.continuation_c.IsSelf {
			re.errorf(process, "[select, client] in CSE rule, the continuation channel should be self, but found %s\n", f.continuation_c.String())
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
	} else {
		re.errorf(process, "in %s, neither the sender ('%s') or continuation ('%s') is self", f.String(), f.to_c.String(), f.continuation_c.String())
	}
}

func (f *BranchForm) Transition(process *Process, re *RuntimeEnvironment) {
	// Should only be referred from within a case
	re.error(process, "Should never transition on branch")
}

func (f *CaseForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of case: %s\n", f.String())

	if f.from_c.IsSelf {
		// CSE rule (provider, -ve)
		//
		// to_c.label<self>
		//    |
		//    |
		//   \|/
		// [case self (...)]

		cseRule := func(message Message) {
			re.logProcessf(LOGRULEDETAILS, process, "[case, provider] finished receiving on self. Received label '%s'\n", message.Label.String())

			if message.Rule != CSE {
				re.errorf(process, "expected CSE, found %s\n", RuleString[message.Rule])
			}

			// Match received label with the available branches
			var new_body Form
			found := false
			for _, j := range f.branches {
				if j.label.Equal(message.Label) {
					// Found a matching label
					found = true
					new_body = j.continuation_e
					// Substitute the payload with 'self'
					new_body.Substitute(j.payload_c, NewSelf(message.Channel1.Ident))
					break
				}
			}

			if !found {
				re.errorf(process, "no matching labels found for %s\n", message.Label)
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
				re.error(process, "expected SEL")
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
				re.errorf(process, "no matching labels found for %s\n", message.Label)
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
	} else {
		re.error(process, "Found a close on a client. A process can only close itself.")
	}
}

func (f *WaitForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of wait: %s\n", f.String())

	if f.to_c.IsSelf {
		re.error(process, "Found a wait on self. Wait should only wait for other channels.")
	}

	clsRule := func(message Message) {
		// CLS rule (client)
		// wait to_c; ...

		re.logProcess(LOGRULE, process, "[wait, client] starting CLS rule")
		re.logProcessf(LOGRULEDETAILS, process, "[wait, client] Received message on channel %s, containing rule: %s\n", f.to_c.String(), RuleString[message.Rule])

		if message.Rule != CLS {
			re.error(process, "expected CLS")
		}

		process.Body = f.continuation_e

		process.finishedRule(CLS, "[wait, client]", "c", re)
		process.transitionLoop(re)
	}

	TransitionByReceiving(process, f.to_c.Channel, clsRule, re)
}

func (f *DropForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of drop: %s\n", f.String())

	if f.client_c.IsSelf {
		re.error(process, "Found a drop on self. Drop only works with other channels.")
	}

	dropRule := func() {
		// Drop does not need to notify the clients being dropped
		// The new [droppable] forward with no providers will take care of this
		// Create structure of new forward process
		newProcess := createDroppableForwardFromClient(process, re, f.client_c)
		re.logProcessf(LOGRULEDETAILS, process, "[drop, client] will create new forward process to drop %s\n", f.client_c.String())
		// Spawn and initiate new forward process
		newProcess.SpawnThenTransition(re)

		process.finishedRule(DROP, "[drop]", "", re)

		process.Body = f.continuation_e
		process.transitionLoop(re)
	}

	TransitionInternally(process, dropRule, re)
}

// Special case: Forward
func (f *ForwardForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of forward: %s\n", f.String())

	if !f.to_c.IsSelf {
		re.error(process, "should forward on self")
	}

	polarity := f.Polarity(re.Typechecked, re.GlobalEnvironment)

	if polarity == types.NEGATIVE && !f.to_drop {
		// -ve
		// least problematic
		// ACTIVE

		message := Message{Rule: FWD, Providers: process.Providers}
		f.from_c.Channel <- message
		re.logProcessf(LOGRULE, process, "[forward, client] sent FWD request to client %s\n", f.from_c.String())

		// todo check if this is needed: process.finishedRule(FWD, "[forward, client]", "", re)
		process.terminateForward(re)

	} else if polarity == types.POSITIVE && !f.to_drop {
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
			process.Body = NewForward(f.to_c, message.Providers[0])
			process.Providers = message.Providers
		case SEL:
			process.Body = NewSelect(f.to_c, message.Label, message.Channel1)
		case CST:
			process.Body = NewCast(f.to_c, message.Channel1)
			// The following are not possible: e.g. a receive does not send anything
		case RCV:
			re.error(process, "a positive forward should never receive RCV messages")
		case CUT:
			re.error(process, "a positive forward should never receive CUT messages")
		case CALL:
			re.error(process, "a positive forward should never receive CALL messages")
		case SPLIT:
			re.error(process, "a positive forward should never receive SPLIT messages")
		case DUP:
			re.error(process, "a positive forward should never receive DUP messages")
		default:
			re.errorf(process, "forward should handle message %s", RuleString[message.Rule])
		}

		process.finishedRule(FWD, "[fwd]", "(+ve)", re)
		process.transitionLoop(re)
	} else if polarity == types.UNKNOWN && !f.to_drop {
		re.error(process, "forward has an unknown polarity")
	}

	// If f.to_drop is true, then the forward should propagate the drops
	if polarity == types.NEGATIVE && f.to_drop {
		// -ve
		// least problematic
		// ACTIVE

		message := Message{Rule: FWD_DROP}
		f.from_c.Channel <- message
		re.logProcessf(LOGRULE, process, "[droppable forward, client] sent FWD_DROP request to client %s\n", f.from_c.String())

		process.terminateForward(re)

	} else if polarity == types.POSITIVE && f.to_drop {
		// +ve
		// PASSIVE: wait before acting

		// Blocks until it received a message. Then this message will be dropped
		message := <-f.from_c.Channel
		re.logProcessf(LOGRULE, process, "[droppable forward, +ve] received message on %s [%s]. This message will be dropped \n", f.from_c.String(), RuleString[message.Rule])

		// Need to handle any clients (aka free names) that will be dropped as a result,
		// e.g. if dropping message <a, b> then you need to cancel a and b as well

		if message.Channel1.Initialized() {
			p := createDroppableForwardFromClient(process, re, message.Channel1)
			p.SpawnThenTransition(re)
		}

		if message.Channel2.Initialized() {
			p := createDroppableForwardFromClient(process, re, message.Channel2)
			p.SpawnThenTransition(re)
		}

		process.terminate(re)

	} else if polarity == types.UNKNOWN && f.to_drop {
		re.error(process, "forward has an unknown polarity")
	}
}

func (f *SplitForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of split: %s\n", f.String())

	if f.from_c.IsSelf {
		// from_cannot be self -- only split other processes
		re.error(process, "should not split on self")
	}

	// Perform SPLIT

	// Prepare new channels
	newSplitNames := []Name{re.CreateFreshChannel(f.channel_one.Ident), re.CreateFreshChannel(f.channel_two.Ident)}

	// Pass through any explicit polarities and types
	for k := range newSplitNames {
		if f.from_c.ExplicitPolarity != nil {
			pol := *f.from_c.ExplicitPolarity
			newSplitNames[k].ExplicitPolarity = &pol
		}
		newSplitNames[k].Type = types.CopyType(f.from_c.Type)
	}

	splitRule := func() {
		// todo check that f.to_c == process.Provider
		re.logProcessf(LOGRULE, process, "[split, client] initiating split for %s into %s\n", f.from_c.String(), NamesToString(newSplitNames))

		currentProcessBody := f.continuation_e
		currentProcessBody.Substitute(f.channel_one, newSplitNames[0])
		currentProcessBody.Substitute(f.channel_two, newSplitNames[1])
		process.Body = currentProcessBody

		process.finishedRule(SPLIT, "[split, client]", "(c)", re)

		// Create structure of new forward process
		fwdSessionType := types.CopyType(f.from_c.Type)
		// TODO maybe use NewSelf
		newProcessBody := NewForward(Name{IsSelf: true, Ident: f.from_c.Ident, Type: fwdSessionType, ExplicitPolarity: f.from_c.ExplicitPolarity}, f.from_c)
		newProcess := NewProcess(newProcessBody, newSplitNames, fwdSessionType, LINEAR, process.Position)
		re.logProcessf(LOGRULEDETAILS, process, "[split, client] will create new forward process providing on %s\n", NamesToString(newSplitNames))
		// Spawn and initiate new forward process
		newProcess.SpawnThenTransition(re)

		process.transitionLoop(re)
	}

	TransitionInternally(process, splitRule, re)
}

func (process *Process) performDUPrule(re *RuntimeEnvironment) {
	if len(process.Providers) == 1 {
		re.error(process, "Cannot duplicate this process")
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
			freshChannels[i][j].Type = types.CopyType(processFreeNames[i].Type)
			freshChannels[i][j].ExplicitPolarity = processFreeNames[i].ExplicitPolarity
			// todo check if the freshChannels have the correct types
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
		dupSessionType := types.CopyType(process.Type) // todo may be wrong type (maybe )
		newDuplicatedProcess := NewProcess(newDuplicatedProcessBody, []Name{newProcessNames[i]}, dupSessionType, process.Shape, process.Position)

		re.logProcessf(LOGRULEDETAILS, process, "[DUP] creating new process (%d): %s\n", i, newDuplicatedProcess.String())

		// Need to spawn the new duplicated processes except the first one (since it's already running in its own thread) -- this was changed... since keeping the first one alive seems to be causing issues
		// if i > 0 {
		newDuplicatedProcess.SpawnThenTransition(re)
		// } else {
		// 	process = newDuplicatedProcess
		// }
	}

	// Create and launch the forward processes to connect the free names (which will implicitly force a chain of further duplications)
	for i := range processFreeNames {
		// Create structure of new forward process
		fwdProcessType := types.CopyType(processFreeNames[i].Type)
		// polarity := processFreeNames[i].Polarity(re.Typechecked, re.GlobalEnvironment)
		newProcessBody := NewForward(Name{IsSelf: true, Ident: processFreeNames[i].Ident, Type: fwdProcessType}, processFreeNames[i]) // includes the polarity of the name not the process
		newProcess := NewProcess(newProcessBody, freshChannels[i], fwdProcessType, LINEAR, process.Position)
		re.logProcessf(LOGRULEDETAILS, process, "[DUP] will create new forward process %s\n", newProcess.String())
		// Spawn and initiate new forward process
		newProcess.SpawnThenTransition(re)
	}

	details := fmt.Sprintf("(Duplicated %s into %s)", process.Providers[0].String(), NamesToString(newProcessNames))
	process.finishedRule(DUP, "[DUP]", details, re)

	// todo remove // Current process has been duplicated, so remove the duplication requirements to continue executing its body [i.e. become interactive]
	// process.OtherProviders = []Name{}
	// process.transitionLoop(re)
	process.terminate(re)
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
			re.errorf(process, "[cast, client] in SHF rule, the continuation channel should be self, but found %s\n", f.continuation_c.String())
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
		// SHF rule (provider, -ve)
		//
		// cast to_c<self>
		//    |
		//    |
		//   \|/
		// [<...> <- shift self; ...]

		shfRule := func(message Message) {
			re.logProcess(LOGRULEDETAILS, process, "[shift, provider] finished sending on self")

			if message.Rule != SHF {
				re.errorf(process, "expected SHF, found %s\n", RuleString[message.Rule])
			}

			new_body := f.continuation_e
			new_body.Substitute(f.continuation_c, NewSelf(message.Channel1.Ident))

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
			re.logProcess(LOGRULE, process, "[shift, client] starting CST rule")
			re.logProcessf(LOGRULEDETAILS, process, "[shift, client] Received message on channel %s, containing rule: %s\n", f.from_c.String(), RuleString[message.Rule])

			if message.Rule != CST {
				re.error(process, "expected CST")
			}

			new_body := f.continuation_e
			new_body.Substitute(f.continuation_c, message.Channel1)

			process.Body = new_body

			process.finishedRule(CST, "[shift, client]", "(c)", re)
			process.transitionLoop(re)
		}

		TransitionByReceiving(process, f.from_c.Channel, cstRule, re)
	}
}

// Debug
func (f *PrintForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of print: %s\n", f.String())

	printRule := func() {
		fmt.Printf("Output from %s: %s\n", NamesToString(process.Providers), f.name_c.String())
		process.finishedRule(PRINT, "[print]", "", re)

		process.Body = f.continuation_e
		process.transitionLoop(re)
	}

	TransitionInternally(process, printRule, re)
}

func (f *PrintLForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of printl: %s\n", f.String())

	printRule := func() {
		fmt.Printf("> %s\n", f.label.String())
		process.finishedRule(PRINTL, "[printl]", "", re)

		process.Body = f.continuation_e
		process.transitionLoop(re)
	}

	TransitionInternally(process, printRule, re)
}

// To keep the log/monitor update with the currently running processes and the transition rules
// being performed, there are the following functions:
//
//	->  finishedRule/3
//	->  finishedHalfRule/3
//	->  processRenamed/1
//	->  terminate/1
//	->  terminateForward/1
//	->  terminateBeforeRename/21
//	->  renamed/1
func (process *Process) finishedRule(rule Rule, prefix, suffix string, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULE, process, "%s finished %s rule %s\n", prefix, RuleString[rule], suffix)

	re.heartbeat <- struct{}{}

	if re.UseMonitor {
		// Update monitor
		re.monitor.MonitorRuleFinished(process, rule)
	}
}

// Process did not finish executing but will be taken over
func (process *Process) processRenamed(re *RuntimeEnvironment) {
	if re.UseMonitor {
		// Update monitor
		re.monitor.MonitorProcessRenamed(process)
	}
}

// Process will terminate
func (process *Process) terminate(re *RuntimeEnvironment) {
	re.logProcess(LOGRULEDETAILS, process, "process terminated successfully")

	// Send heartbeat
	re.heartbeat <- struct{}{}

	if re.UseMonitor {
		// Update monitor
		re.monitor.MonitorProcessTerminated(process)
	}

	// Update dead process count
	atomic.AddUint64(&re.deadProcessCount, 1)
}

// A forward process will terminate, but its providers will be used by other processes being forwarded
func (process *Process) terminateForward(re *RuntimeEnvironment) {
	re.logProcess(LOGRULEDETAILS, process, "process will change by forwarding its provider")

	// Send heartbeat
	re.heartbeat <- struct{}{}

	if re.UseMonitor {
		// Update monitor
		re.monitor.MonitorRuleFinished(process, FWD)
	}

	// Will not update the dead process count since the process' provider names will 'live' on
}

// // A forward process will terminate, and its providers will be be dropped so the process will die as well
// func (process *Process) terminateForwardDropped(re *RuntimeEnvironment) {
// 	re.logProcess(LOGRULEDETAILS, process, "process will change by forwarding its provider")

// // Send heartbeat
// re.heartbeat <- struct{}{}

// 	if re.Debug {
// 		// Update monitor
// 		re.monitor.MonitorRuleFinished(process, FWD)

// 		// Update dead process count
// 		atomic.AddUint64(&re.deadProcessCount, 1)
// 	}
// }

func (process *Process) terminateBeforeRename(oldProviders, newProviders []Name, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "process renamed from %s to %s\n", NamesToString(oldProviders), NamesToString(newProviders))

	// Send heartbeat
	re.heartbeat <- struct{}{}

	if re.UseMonitor {
		// Update monitor
		// todo change
		re.monitor.MonitorProcessForwarded(process)
	}

	// Update dead process count
	atomic.AddUint64(&re.deadProcessCount, 1)
}

func (process *Process) renamed(oldProviders, newProviders []Name, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "process renamed from %s to %s\n", NamesToString(oldProviders), NamesToString(newProviders))

	// Send heartbeat
	re.heartbeat <- struct{}{}

	// Although the old providers should be closed (i.e. die), the process itself does not die. It lives on using the new provider names.
}

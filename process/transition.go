package process

import (
	"fmt"
	"sync/atomic"
)

// Initiates new processes
func (process *Process) Transition(re *RuntimeEnvironment) {
	// todo make this atomic
	atomic.AddUint64(&re.ProcessCount, 1)

	go TransitionLoop(process, re)
}

// Entry point before a process transitions
func TransitionLoop(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGPROCESSING, process, "Process transitioning: %s\n", process.Body.String())
	process.Body.Transition(process, re)
}

func finishedRule(process *Process, rule Rule, prefix, suffix string, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULE, process, "%s finished %s rule %s", prefix, rule, suffix)

	if re.debug {
		// Update monitor
		re.monitor.monitorChan <- MonitorUpdate{process: *process, rule: rule}
	}
}

//
//
// TransitionAsClient
// TransitionAsProvider
//

// Transition according to the present body form (e.g. send, receive, ...)
func (f *SendForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of send: %s\n", f.String())

	if f.to_c.IsSelf {
		select {
		case pm := <-process.Provider.PriorityChannel:
			switch pm.Action {
			case FWD:
				re.logProcessf(LOGRULEDETAILS, process, "[Priority Msg] Received FWD request. Continuing as %s\n", pm.Channel1.String())
				// Close current channel and switch to new one
				close(process.Provider.Channel)
				close(process.Provider.PriorityChannel)
				process.Provider = pm.Channel1
				TransitionLoop(process, re)
			case SPLIT_DUP_FWD:
				re.logProcessf(LOGRULE, process, "[Priority Msg] Received SPLIT_DUP_FWD request. Will split in %d processes: %s\n", len(pm.Channels), NamesToString(pm.Channels))

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
					freshChannels[i] = make([]Name, len(pm.Channels))

					for j := range pm.Channels {
						freshChannels[i][j] = re.CreateFreshChannel(processFreeNames[i].Ident)
					}
				}

				chanString := ""
				for i := range freshChannels {
					chanString += processFreeNames[i].String() + ": {"
					chanString += NamesToString(freshChannels[i])
					chanString += "}; "
				}
				re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] NEW Free names: %s\n", chanString)

				// for i := range freshChannels {

				for i := range pm.Channels {
					// Prepare process to spawn, by substituting all free names with the unique ones just created
					newProcessBody := CopyForm(f)
					for k := range freshChannels {
						newProcessBody.Substitute(processFreeNames[k], freshChannels[k][i])
					}

					// Create and spawn the new processes
					// Set its provider to the channel received in the SPLIT_DUP_FWD request
					newSplitProcess := NewProcess(newProcessBody, pm.Channels[i], process.Shape, process.FunctionDefinitions)
					newSplitProcess.Transition(re)
				}
				// }

				for i := range processFreeNames {
					// Implicitly we are doing SPLIT -> DUP -> FWD procedure in one move
					// send split to channel processFreeNames[i] containing the new channels freshChannels[i]
					re.logProcessf(LOGRULE, process, "[SPLIT_DUP_FWD] Asking %s to split into %s\n", processFreeNames[i].String(), NamesToString(freshChannels[i]))
					processFreeNames[i].PriorityChannel <- PriorityMessage{Action: SPLIT_DUP_FWD, Channels: freshChannels[i]}
				}

				re.logProcess(LOGRULE, process, "[SPLIT_DUP_FWD] Rule SPLIT finished (SPLIT -> DUP -> FWD) (p)")

				// The process body has been delegated to the freshly spawned processes, so die
				// todo die
				process.Terminate(re)
			default:
				re.logProcessHighlight(LOGRULEDETAILS, process, "RECEIVED something else --- fix\n")
				// todo panic
			}
		// case <-process.Provider.Channel:
		case process.Provider.Channel <- Message{Rule: SND, Channel1: f.payload_c, Channel2: f.continuation_c}:
			// SND rule (provider)
			// re.logProcess(LOGRULE, process, "[send, provider] starting SND rule")
			// re.logProcessf(LOGRULEDETAILS, process, "[send, provider] Send to self (%s), proceeding with SND\n", process.Provider.String())

			// process.Provider.Channel <- Message{Rule: SND, Channel1: f.payload_c, Channel2: f.continuation_c}
			re.logProcess(LOGRULEDETAILS, process, "[send, provider] finished sending on self")

			// finishedRule(process, SND, "[send, provider]", "(p)", re)
			re.logProcess(LOGRULE, process, "[send, provider] finished SND rule (p)")

			// todo terminate goroutine. inform monitor
			process.Terminate(re)
		}
	} else {
		// RCV rule (client)

		select {
		case pm := <-process.Provider.PriorityChannel:
			switch pm.Action {
			case FWD:
				re.logProcessf(LOGRULEDETAILS, process, "[Priority Msg] Received FWD request. Continuing as %s\n", pm.Channel1.String())
				// Close current channel and switch to new one
				close(process.Provider.Channel)
				close(process.Provider.PriorityChannel)
				process.Provider = pm.Channel1
				TransitionLoop(process, re)

			case SPLIT_DUP_FWD:
				re.logProcessf(LOGRULE, process, "[Priority Msg] Received SPLIT_DUP_FWD request. Will split in %d processes: %s\n", len(pm.Channels), NamesToString(pm.Channels))

				processFreeNames := process.Body.FreeNames()
				// re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] Free names: %s\n", NamesToString(processFreeNames))

				// Create an array of arrays containing fresh channels
				// E.g. if we are splitting to two channels, and have 1 free name called a, then
				// we create freshChannels containing:
				//   [][]Names{
				//       [Names]{a', a''},
				//   }
				// where a' and a'' are the new fresh channels that will be substituted in place of a
				freshChannels := make([][]Name, len(processFreeNames))

				for i := range freshChannels {
					freshChannels[i] = make([]Name, len(pm.Channels))

					for j := range pm.Channels {
						freshChannels[i][j] = re.CreateFreshChannel(processFreeNames[i].Ident)
					}
				}

				chanString := ""
				for i := range freshChannels {
					chanString += processFreeNames[i].String() + ": {"
					chanString += NamesToString(freshChannels[i])
					chanString += "}; "
				}
				re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] NEW Free names: %s\n", chanString)

				// for i := range freshChannels {

				for i := range pm.Channels {
					// Prepare process to spawn, by substituting all free names with the unique ones just created
					newProcessBody := CopyForm(f)
					for k := range freshChannels {
						newProcessBody.Substitute(processFreeNames[k], freshChannels[k][i])
					}

					// Create and spawn the new processes
					// Set its provider to the channel received in the SPLIT_DUP_FWD request
					newSplitProcess := NewProcess(newProcessBody, pm.Channels[i], process.Shape, process.FunctionDefinitions)
					newSplitProcess.Transition(re)
				}
				// }

				for i := range processFreeNames {
					// Implicitly we are doing SPLIT -> DUP -> FWD procedure in one move
					// send split to channel processFreeNames[i] containing the new channels freshChannels[i]
					re.logProcessf(LOGRULE, process, "[SPLIT_DUP_FWD] Asking %s to split into %s\n", processFreeNames[i].String(), NamesToString(freshChannels[i]))
					processFreeNames[i].PriorityChannel <- PriorityMessage{Action: SPLIT_DUP_FWD, Channels: freshChannels[i]}
				}

				re.logProcess(LOGRULE, process, "[SPLIT_DUP_FWD] Rule SPLIT finished (SPLIT -> DUP -> FWD) (p)")

				// The process body has been delegated to the freshly spawned processes, so die
				// todo die
				process.Terminate(re)
			default:
				re.logProcessHighlight(LOGRULEDETAILS, process, "RECEIVED something else --- fix\n")
				// todo panic
			}
		// case f.to_c.Channel <- Message{}:
		// 	re.logProcessHighlightf(LOGRULEDETAILS, process, "AAAA")
		case message := <-f.to_c.Channel:
			re.logProcess(LOGRULE, process, "[send, client] starting RCV rule")
			re.logProcessf(LOGRULEDETAILS, process, "Should received message on channel %s, containing message rule RCV\n", f.to_c.String())
			// message := <-f.to_c.Channel
			close(f.to_c.Channel)
			close(f.to_c.PriorityChannel)

			// todo check that rule matches RCV

			new_body := message.ContinuationBody
			new_body.Substitute(message.Channel1, f.payload_c)
			new_body.Substitute(message.Channel2, f.continuation_c)

			if !f.continuation_c.IsSelf {
				// todo error
				re.logProcessf(LOGERROR, process, "[send, client] in RCV rule, the continuation channel should be self, but found %s\n", f.continuation_c.String())
				panic("Expected self but found something else")
			}

			re.logProcess(LOGRULE, process, "[send, client] finished RCV rule (c)")

			process.Body = new_body
			TransitionLoop(process, re)
		}
	}
}

func (f *ReceiveForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of receive: %s\n", f.String())

	if f.from_c.IsSelf {
		// todo check where the select is random (if both can succeed) -- not sure what happens
		select {
		case pm := <-process.Provider.PriorityChannel:
			switch pm.Action {
			case FWD:
				re.logProcessf(LOGRULEDETAILS, process, "[Priority Msg] Received FWD request. Continuing as %s\n", pm.Channel1.String())
				// Close current channel and switch to new one
				close(process.Provider.Channel)
				close(process.Provider.PriorityChannel)
				process.Provider = pm.Channel1
				TransitionLoop(process, re)

			case SPLIT_DUP_FWD:
				re.logProcessf(LOGRULE, process, "[Priority Msg] Received SPLIT_DUP_FWD request. Will split in %d processes: %s\n", len(pm.Channels), NamesToString(pm.Channels))

				processFreeNames := process.Body.FreeNames()
				// re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] Free names: %s\n", NamesToString(processFreeNames))

				// Create an array of arrays containing fresh channels
				// E.g. if we are splitting to two channels, and have 1 free name called a, then
				// we create freshChannels containing:
				//   [][]Names{
				//       [Names]{a', a''},
				//   }
				// where a' and a'' are the new fresh channels that will be substituted in place of a
				freshChannels := make([][]Name, len(processFreeNames))

				for i := range freshChannels {
					freshChannels[i] = make([]Name, len(pm.Channels))

					for j := range pm.Channels {
						freshChannels[i][j] = re.CreateFreshChannel(processFreeNames[i].Ident)
					}
				}

				chanString := ""
				for i := range freshChannels {
					chanString += processFreeNames[i].String() + ": {"
					chanString += NamesToString(freshChannels[i])
					chanString += "}; "
				}
				re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] NEW Free names: %s\n", chanString)

				// for i := range freshChannels {

				for i := range pm.Channels {
					// Prepare process to spawn, by substituting all free names with the unique ones just created
					newProcessBody := CopyForm(f)
					for k := range freshChannels {
						newProcessBody.Substitute(processFreeNames[k], freshChannels[k][i])
					}

					// Create and spawn the new processes
					// Set its provider to the channel received in the SPLIT_DUP_FWD request
					newSplitProcess := NewProcess(newProcessBody, pm.Channels[i], process.Shape, process.FunctionDefinitions)
					newSplitProcess.Transition(re)
				}
				// }

				for i := range processFreeNames {
					// Implicitly we are doing SPLIT -> DUP -> FWD procedure in one move
					// send split to channel processFreeNames[i] containing the new channels freshChannels[i]
					re.logProcessf(LOGRULE, process, "[SPLIT_DUP_FWD] Asking %s to split into %s\n", processFreeNames[i].String(), NamesToString(freshChannels[i]))
					processFreeNames[i].PriorityChannel <- PriorityMessage{Action: SPLIT_DUP_FWD, Channels: freshChannels[i]}
				}

				re.logProcess(LOGRULE, process, "[SPLIT_DUP_FWD] Rule SPLIT finished (SPLIT -> DUP -> FWD) (p)")

				// The process body has been delegated to the freshly spawned processes, so die
				// todo die
				process.Terminate(re)
			default:
				re.logProcessHighlight(LOGRULEDETAILS, process, "RECEIVED something else --- fix")
				// todo panic
			}
		// case <-process.Provider.Channel:
		case process.Provider.Channel <- Message{Rule: RCV, ContinuationBody: f.continuation_e, Channel1: f.payload_c, Channel2: f.continuation_c}:
			// RCV rule (provider)
			// default:
			// re.logProcessf(LOGRULEDETAILS, process, "[receive, provider] Send to self (%s), proceeding with RCV\n", process.Provider.String())

			// process.Provider.Channel <- Message{Rule: RCV, ContinuationBody: f.continuation_e, Channel1: f.payload_c, Channel2: f.continuation_c}
			re.logProcess(LOGRULEDETAILS, process, "[receive, provider] finished sending on self")

			re.logProcess(LOGRULE, process, "[receive, provider] finished RCV rule (p)")
			// todo terminate goroutine. inform monitor
			process.Terminate(re)
		}
	} else {
		// SND rule (client)
		// todo ask for controller permission

		select {
		case pm := <-process.Provider.PriorityChannel:
			switch pm.Action {
			case FWD:
				re.logProcessf(LOGRULEDETAILS, process, "[Priority Msg] Received FWD request. Continuing as %s\n", pm.Channel1.String())
				// Close current channel and switch to new one
				close(process.Provider.Channel)
				close(process.Provider.PriorityChannel)
				process.Provider = pm.Channel1
				TransitionLoop(process, re)

			case SPLIT_DUP_FWD:
				re.logProcessf(LOGRULE, process, "[Priority Msg] Received SPLIT_DUP_FWD request. Will split in %d processes: %s\n", len(pm.Channels), NamesToString(pm.Channels))

				processFreeNames := process.Body.FreeNames()
				// re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] Free names: %s\n", NamesToString(processFreeNames))

				// Create an array of arrays containing fresh channels
				// E.g. if we are splitting to two channels, and have 1 free name called a, then
				// we create freshChannels containing:
				//   [][]Names{
				//       [Names]{a', a''},
				//   }
				// where a' and a'' are the new fresh channels that will be substituted in place of a
				freshChannels := make([][]Name, len(processFreeNames))

				for i := range freshChannels {
					freshChannels[i] = make([]Name, len(pm.Channels))

					for j := range pm.Channels {
						freshChannels[i][j] = re.CreateFreshChannel(processFreeNames[i].Ident)
					}
				}

				chanString := ""
				for i := range freshChannels {
					chanString += processFreeNames[i].String() + ": {"
					chanString += NamesToString(freshChannels[i])
					chanString += "}; "
				}
				re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] NEW Free names: %s\n", chanString)

				// for i := range freshChannels {

				for i := range pm.Channels {
					// Prepare process to spawn, by substituting all free names with the unique ones just created
					newProcessBody := CopyForm(f)
					for k := range freshChannels {
						newProcessBody.Substitute(processFreeNames[k], freshChannels[k][i])
					}

					// Create and spawn the new processes
					// Set its provider to the channel received in the SPLIT_DUP_FWD request
					newSplitProcess := NewProcess(newProcessBody, pm.Channels[i], process.Shape, process.FunctionDefinitions)
					newSplitProcess.Transition(re)
				}
				// }

				for i := range processFreeNames {
					// Implicitly we are doing SPLIT -> DUP -> FWD procedure in one move
					// send split to channel processFreeNames[i] containing the new channels freshChannels[i]
					re.logProcessf(LOGRULE, process, "[SPLIT_DUP_FWD] Asking %s to split into %s\n", processFreeNames[i].String(), NamesToString(freshChannels[i]))
					processFreeNames[i].PriorityChannel <- PriorityMessage{Action: SPLIT_DUP_FWD, Channels: freshChannels[i]}
				}

				re.logProcess(LOGRULE, process, "[SPLIT_DUP_FWD] Rule SPLIT finished (SPLIT -> DUP -> FWD) (p)")

				// The process body has been delegated to the freshly spawned processes, so die
				// todo die
				process.Terminate(re)
			default:
				re.logProcessHighlight(LOGRULEDETAILS, process, "RECEIVED something else --- fix\n")
				// todo panic
			}
		// case f.from_c.Channel <- Message{}:
		case message := <-f.from_c.Channel:
			re.logProcess(LOGRULE, process, "[receive, client] starting SND rule")
			re.logProcessf(LOGRULEDETAILS, process, "[receive, client] proceeding with SND, will receive from %s\n", f.from_c.String())

			// message := <-f.from_c.Channel
			// close(f.from_c.Channel)
			// close(f.from_c.PriorityChannel)

			re.logProcessf(LOGRULEDETAILS, process, "Received message on channel %s, containing message.Rule %d\n", f.from_c.String(), message.Rule)

			new_body := f.continuation_e
			new_body.Substitute(f.payload_c, message.Channel1)
			new_body.Substitute(f.continuation_c, message.Channel2)

			re.logProcess(LOGRULE, process, "[receive, client] finished SND rule (c)")

			process.Body = new_body
			TransitionLoop(process, re)
		}
	}
}

func (f *ForwardForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of forward: %s\n", f.String())
	// todo check priority messages

	if f.to_c.IsSelf {
		select {
		case pm := <-process.Provider.PriorityChannel:
			switch pm.Action {
			case FWD:
				re.logProcessf(LOGRULEDETAILS, process, "[Priority Msg] Received FWD request. Continuing as %s\n", pm.Channel1.String())
				// Close current channel and switch to new one
				close(process.Provider.Channel)
				close(process.Provider.PriorityChannel)
				process.Provider = pm.Channel1
				TransitionLoop(process, re)

			case SPLIT_DUP_FWD:
				re.logProcessf(LOGRULE, process, "[Priority Msg] Received SPLIT_DUP_FWD request. Will split in %d processes: %s\n", len(pm.Channels), NamesToString(pm.Channels))

				processFreeNames := process.Body.FreeNames()
				// re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] Free names: %s\n", NamesToString(processFreeNames))

				// Create an array of arrays containing fresh channels
				// E.g. if we are splitting to two channels, and have 1 free name called a, then
				// we create freshChannels containing:
				//   [][]Names{
				//       [Names]{a', a''},
				//   }
				// where a' and a'' are the new fresh channels that will be substituted in place of a
				freshChannels := make([][]Name, len(processFreeNames))

				for i := range freshChannels {
					freshChannels[i] = make([]Name, len(pm.Channels))

					for j := range pm.Channels {
						freshChannels[i][j] = re.CreateFreshChannel(processFreeNames[i].Ident)
					}
				}

				chanString := ""
				for i := range freshChannels {
					chanString += processFreeNames[i].String() + ": {"
					chanString += NamesToString(freshChannels[i])
					chanString += "}; "
				}
				re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] NEW Free names: %s\n", chanString)

				// for i := range freshChannels {

				for i := range pm.Channels {
					// Prepare process to spawn, by substituting all free names with the unique ones just created
					newProcessBody := CopyForm(f)
					for k := range freshChannels {
						newProcessBody.Substitute(processFreeNames[k], freshChannels[k][i])
					}

					// Create and spawn the new processes
					// Set its provider to the channel received in the SPLIT_DUP_FWD request
					newSplitProcess := NewProcess(newProcessBody, pm.Channels[i], process.Shape, process.FunctionDefinitions)
					newSplitProcess.Transition(re)
				}
				// }

				for i := range processFreeNames {
					// Implicitly we are doing SPLIT -> DUP -> FWD procedure in one move
					// send split to channel processFreeNames[i] containing the new channels freshChannels[i]
					re.logProcessf(LOGRULE, process, "[SPLIT_DUP_FWD] Asking %s to split into %s\n", processFreeNames[i].String(), NamesToString(freshChannels[i]))
					processFreeNames[i].PriorityChannel <- PriorityMessage{Action: SPLIT_DUP_FWD, Channels: freshChannels[i]}
				}

				re.logProcess(LOGRULE, process, "[SPLIT_DUP_FWD] Rule SPLIT finished (SPLIT -> DUP -> FWD) (p)")

				// The process body has been delegated to the freshly spawned processes, so die
				// todo die
				process.Terminate(re)
			default:
				re.logProcessHighlight(LOGRULEDETAILS, process, "RECEIVED something else --- fix")
				// todo panic
			}
		case f.from_c.PriorityChannel <- PriorityMessage{Action: FWD, Channel1: process.Provider}:
			// todo check that f.to_c == process.Provider
			// f.from_c.Channel <- Message{Rule: FWD, Channel1: process.Provider}
			re.logProcessf(LOGRULE, process, "[forward, client] sent FWD request to priority channel %s\n", f.from_c.String())

			re.logProcess(LOGRULE, process, "[forward, client] finished FWD rule")
			// todo terminate goroutine. inform monitor
			// channel providing on fwd should not be closed, however this process dies
			// process.Terminate(re)
			process.Terminate(re)
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
	} else {
		// Prepare new channels
		newSplitNames := []Name{re.CreateFreshChannel(f.channel_one.Ident), re.CreateFreshChannel(f.channel_two.Ident)}
		select {
		case pm := <-process.Provider.PriorityChannel:
			switch pm.Action {
			case FWD:
				re.logProcessf(LOGRULEDETAILS, process, "[Priority Msg] Received FWD request. Continuing as %s\n", pm.Channel1.String())
				// Close current channel and switch to new one
				close(process.Provider.Channel)
				close(process.Provider.PriorityChannel)
				process.Provider = pm.Channel1
				TransitionLoop(process, re)

			case SPLIT_DUP_FWD:
				re.logProcessf(LOGRULE, process, "[Priority Msg] Received SPLIT_DUP_FWD request. Will split in %d processes: %s\n", len(pm.Channels), NamesToString(pm.Channels))

				processFreeNames := process.Body.FreeNames()
				// re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] Free names: %s\n", NamesToString(processFreeNames))

				// Create an array of arrays containing fresh channels
				// E.g. if we are splitting to two channels, and have 1 free name called a, then
				// we create freshChannels containing:
				//   [][]Names{
				//       [Names]{a', a''},
				//   }
				// where a' and a'' are the new fresh channels that will be substituted in place of a
				freshChannels := make([][]Name, len(processFreeNames))

				for i := range freshChannels {
					freshChannels[i] = make([]Name, len(pm.Channels))

					for j := range pm.Channels {
						freshChannels[i][j] = re.CreateFreshChannel(processFreeNames[i].Ident)
					}
				}

				chanString := ""
				for i := range freshChannels {
					chanString += processFreeNames[i].String() + ": {"
					chanString += NamesToString(freshChannels[i])
					chanString += "}; "
				}
				re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] NEW Free names: %s\n", chanString)

				// for i := range freshChannels {

				for i := range pm.Channels {
					// Prepare process to spawn, by substituting all free names with the unique ones just created
					newProcessBody := CopyForm(f)
					for k := range freshChannels {
						newProcessBody.Substitute(processFreeNames[k], freshChannels[k][i])
					}

					// Create and spawn the new processes
					// Set its provider to the channel received in the SPLIT_DUP_FWD request
					newSplitProcess := NewProcess(newProcessBody, pm.Channels[i], process.Shape, process.FunctionDefinitions)
					newSplitProcess.Transition(re)
				}
				// }

				for i := range processFreeNames {
					// Implicitly we are doing SPLIT -> DUP -> FWD procedure in one move
					// send split to channel processFreeNames[i] containing the new channels freshChannels[i]
					re.logProcessf(LOGRULE, process, "[SPLIT_DUP_FWD] Asking %s to split into %s\n", processFreeNames[i].String(), NamesToString(freshChannels[i]))
					processFreeNames[i].PriorityChannel <- PriorityMessage{Action: SPLIT_DUP_FWD, Channels: freshChannels[i]}
				}

				re.logProcess(LOGRULE, process, "[SPLIT_DUP_FWD] Rule SPLIT finished (SPLIT -> DUP -> FWD) (p)")

				// The process body has been delegated to the freshly spawned processes, so die
				// todo die
				process.Terminate(re)
			default:
				re.logProcessHighlight(LOGRULEDETAILS, process, "RECEIVED something else --- fix")
				// todo panic
			}
		case f.from_c.PriorityChannel <- PriorityMessage{Action: SPLIT_DUP_FWD, Channels: newSplitNames}:
			// todo check that f.to_c == process.Provider
			re.logProcessf(LOGRULE, process, "[split, client] should send SPLIT to priority channel %s\n", f.from_c.String())

			currentProcessBody := f.continuation_e
			currentProcessBody.Substitute(f.channel_one, newSplitNames[0])
			currentProcessBody.Substitute(f.channel_two, newSplitNames[1])
			process.Body = currentProcessBody

			re.logProcess(LOGRULE, process, "[split, client] finished SPLIT rule (c)")

			TransitionLoop(process, re)
		}
	}
}

// CUT rule (Spawn new process) - provider only
func (f *NewForm) Transition(process *Process, re *RuntimeEnvironment) {
	re.logProcessf(LOGRULEDETAILS, process, "transition of new: %s\n", f.String())

	select {
	case pm := <-process.Provider.PriorityChannel:
		switch pm.Action {
		case FWD:
			re.logProcessf(LOGRULEDETAILS, process, "[Priority Msg] Received FWD request. Continuing as %s\n", pm.Channel1.String())
			// Close current channel and switch to new one
			close(process.Provider.Channel)
			close(process.Provider.PriorityChannel)
			process.Provider = pm.Channel1
			TransitionLoop(process, re)

		case SPLIT_DUP_FWD:
			re.logProcessf(LOGRULE, process, "[Priority Msg] Received SPLIT_DUP_FWD request. Will split in %d processes: %s\n", len(pm.Channels), NamesToString(pm.Channels))

			processFreeNames := process.Body.FreeNames()
			// re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] Free names: %s\n", NamesToString(processFreeNames))

			// Create an array of arrays containing fresh channels
			// E.g. if we are splitting to two channels, and have 1 free name called a, then
			// we create freshChannels containing:
			//   [][]Names{
			//       [Names]{a', a''},
			//   }
			// where a' and a'' are the new fresh channels that will be substituted in place of a
			freshChannels := make([][]Name, len(processFreeNames))

			for i := range freshChannels {
				freshChannels[i] = make([]Name, len(pm.Channels))

				for j := range pm.Channels {
					freshChannels[i][j] = re.CreateFreshChannel(processFreeNames[i].Ident)
				}
			}

			chanString := ""
			for i := range freshChannels {
				chanString += processFreeNames[i].String() + ": {"
				chanString += NamesToString(freshChannels[i])
				chanString += "}; "
			}
			re.logProcessf(LOGRULEDETAILS, process, "[SPLIT_DUP_FWD] NEW Free names: %s\n", chanString)

			// for i := range freshChannels {

			for i := range pm.Channels {
				// Prepare process to spawn, by substituting all free names with the unique ones just created
				newProcessBody := CopyForm(f)
				for k := range freshChannels {
					newProcessBody.Substitute(processFreeNames[k], freshChannels[k][i])
				}

				// Create and spawn the new processes
				// Set its provider to the channel received in the SPLIT_DUP_FWD request
				newSplitProcess := NewProcess(newProcessBody, pm.Channels[i], process.Shape, process.FunctionDefinitions)
				newSplitProcess.Transition(re)
			}
			// }

			for i := range processFreeNames {
				// Implicitly we are doing SPLIT -> DUP -> FWD procedure in one move
				// send split to channel processFreeNames[i] containing the new channels freshChannels[i]
				re.logProcessf(LOGRULE, process, "[SPLIT_DUP_FWD] Asking %s to split into %s\n", processFreeNames[i].String(), NamesToString(freshChannels[i]))
				processFreeNames[i].PriorityChannel <- PriorityMessage{Action: SPLIT_DUP_FWD, Channels: freshChannels[i]}
			}

			re.logProcess(LOGRULE, process, "[SPLIT_DUP_FWD] Rule SPLIT finished (SPLIT -> DUP -> FWD)")

			// The process body has been delegated to the freshly spawned processes, so die
			// todo die
			process.Terminate(re)
		default:
			re.logProcessHighlight(LOGRULEDETAILS, process, "RECEIVED something else --- fix")
			// todo panic
		}
	default:
		re.logProcess(LOGRULEDETAILS, process, "[new] will create new process")

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
		newProcess := NewProcess(newProcessBody, newChannel, LINEAR, process.FunctionDefinitions)

		re.logProcessf(LOGRULEDETAILS, process, "[new] will create new process with channel %s\n", newChannel.String())

		// Spawn and initiate new process
		newProcess.Transition(re)

		re.logProcess(LOGRULE, process, "[new] finished CUT rule")
		// Continue executing current process
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

func (f *CloseForm) Transition(process *Process, re *RuntimeEnvironment) {
	fmt.Print("transition of close: ")
	fmt.Println(f.String())
}

func (process *Process) Terminate(re *RuntimeEnvironment) {
	re.logProcess(LOGRULEDETAILS, process, "process terminated successfully")
}

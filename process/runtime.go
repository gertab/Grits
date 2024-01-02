package process

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slices"
)

// RuntimeEnvironment provides the instance instance for processes to be initialized and executed
type RuntimeEnvironment struct {

	// Keeps a (read only) global environment containing the function definitions and session types
	GlobalEnvironment *GlobalEnvironment

	// Keeps count of how many processes were spawned (only for debug info)
	processCount uint64
	// Keeps count of how many processes died (only for debug info)
	deadProcessCount uint64

	// Keeps context to control the new processes
	ctx context.Context
	// If heartbeats stop, then the processes are killed
	heartbeat chan struct{}

	// Debugging info
	UseMonitor bool
	// Colored output
	Color bool
	// Keeps counter of the number of channels created
	debugChannelCounter uint64
	// Monitor info
	monitor *Monitor
	// Slow execution speed
	Delay time.Duration
	// Errors are sent to this channel
	errorChan chan error
	// Chooses how the transitions are performed ([non-]polarized [a]synchronous)
	ExecutionVersion Execution_Version
	// Flag to see whether the typechecker was used or not (i.e. if true, then all names have types)
	Typechecked bool
}

type Execution_Version int

const (
	/* polarized + async */
	NORMAL_ASYNC Execution_Version = iota
	/* polarized + sync */
	NORMAL_SYNC
	/* non-polarized forwards + sync */
	NON_POLARIZED_SYNC
	// NON_POLARIZED_ASYNC /* problematic */
)

func NewRuntimeEnvironment() (*RuntimeEnvironment, context.CancelFunc) {
	re := &RuntimeEnvironment{
		UseMonitor:          true,
		Color:               true,
		ExecutionVersion:    NORMAL_ASYNC,
		debugChannelCounter: 0,
		processCount:        0,
		deadProcessCount:    0,
		errorChan:           make(chan error),
		Typechecked:         false,
		heartbeat:           make(chan struct{}, 1),
	}

	// Prepare context with cancellation
	var cancel context.CancelFunc
	re.ctx, cancel = context.WithCancel(context.Background())

	return re, cancel
}

// Entry point for execution
func InitializeProcesses(processes []*Process, globalEnv *GlobalEnvironment, subscriber *SubscriberInfo, re *RuntimeEnvironment) *RuntimeEnvironment {

	if re == nil {
		re = &RuntimeEnvironment{
			UseMonitor:       true,
			Color:            true,
			Delay:            1000 * time.Millisecond,
			ExecutionVersion: NORMAL_ASYNC,
			Typechecked:      false,
		}
	}

	if globalEnv != nil {
		re.GlobalEnvironment = globalEnv
	}

	if re.GlobalEnvironment.LogLevels == nil {
		l := []LogLevel{
			LOGINFO,
			LOGPROCESSING,
			LOGRULE,
			LOGRULEDETAILS,
			LOGMONITOR,
		}

		re.GlobalEnvironment.LogLevels = l
	}

	re.processCount = 0
	re.deadProcessCount = 0
	re.debugChannelCounter = 0
	re.errorChan = make(chan error)

	// Prepare context with cancellation
	var cancel context.CancelFunc
	re.ctx, cancel = context.WithCancel(context.Background())
	re.heartbeat = make(chan struct{}, 1)
	defer cancel()

	re.logf(LOGINFO, "\nInitializing %d processes\n", len(processes))

	channels := re.CreateChannelForEachProcess(processes)

	re.SubstituteNameInitialization(processes, channels)

	if re.UseMonitor {
		startedWg := new(sync.WaitGroup)
		startedWg.Add(1)

		re.InitializeMonitor(startedWg, subscriber)

		// Ensure that both servers are running
		startedWg.Wait()
	}

	go re.HeartbeatReceiver(50*time.Millisecond, cancel)

	re.StartTransitions(processes)

	select {
	case <-re.ctx.Done():
	// case <-re.monitor.monitorFinished:
	// Monitor terminated successfully
	case err := <-re.errorChan:
		log.Fatal(err)
	}

	re.logf(LOGINFO, "End process count: %d (%d)\n", re.ProcessCount(), re.DeadProcessCount())

	return re
}

func (re *RuntimeEnvironment) ProcessCount() uint64 {
	return re.processCount
}

func (re *RuntimeEnvironment) DeadProcessCount() uint64 {
	return re.deadProcessCount
}

// Create the initial channels required. E.g. for a process prc[c1], a channel with Ident: c1 is created
func (re *RuntimeEnvironment) CreateChannelForEachProcess(processes []*Process) []NameInitialization {

	var channels []NameInitialization

	for i := 0; i < len(processes); i++ {
		// todo ensure that len(old_provider) >= 0
		for j := 0; j < len(processes[i].Providers); j++ {
			old_provider := processes[i].Providers[j]
			new_provider := re.CreateFreshChannel(old_provider.Ident)

			// Set new channel as the providing channel for the process
			processes[i].Providers[j] = new_provider
			channels = append(channels, NameInitialization{old_provider, new_provider})
		}
	}

	return channels
}

// Create new channel
func (re *RuntimeEnvironment) CreateFreshChannel(ident string) Name {
	// The channel ID is used for debugging
	atomic.AddUint64(&re.debugChannelCounter, 1)

	// Create new channel and assign a name to it
	var mChan chan Message
	switch re.ExecutionVersion {

	case NORMAL_ASYNC:
		mChan = make(chan Message, 1)
	case NORMAL_SYNC:
		mChan = make(chan Message)
	case NON_POLARIZED_SYNC:
		mChan = make(chan Message)
	}

	// Control channel is only used in the non-polarized version
	var cmChan chan ControlMessage

	if re.ExecutionVersion == NON_POLARIZED_SYNC {
		cmChan = make(chan ControlMessage)
		// Not needed in the case of NORMAL_ASYNC or NORMAL_SYNC
	}

	return Name{
		Ident:          ident,
		Channel:        mChan,
		ChannelID:      re.debugChannelCounter,
		ControlChannel: cmChan,
		IsSelf:         false,
	}
}

func (re *RuntimeEnvironment) InitializeMonitor(startedWg *sync.WaitGroup, subscriber *SubscriberInfo) {
	// Declare new monitor
	re.monitor = NewMonitor(re, subscriber)

	// Start monitor on new thread
	go re.monitor.startMonitor(startedWg)
}

func (re *RuntimeEnvironment) InitializeGivenMonitor(startedWg *sync.WaitGroup, monitor *Monitor, subscriber *SubscriberInfo) {
	// Declare new monitor
	re.monitor = monitor

	// Start monitor on new thread
	go re.monitor.startMonitor(startedWg)
}

// Used after initialization to substitute known names to the actual channel
func (re *RuntimeEnvironment) SubstituteNameInitialization(processes []*Process, channels []NameInitialization) {
	for i := 0; i < len(processes); i++ {
		for _, c := range channels {
			// Substitute all free names in a body
			processes[i].Body.Substitute(c.old, c.new)
		}
	}
}

func (re *RuntimeEnvironment) StartTransitions(processes []*Process) {
	for _, p := range processes {
		p_uniq := p

		switch re.ExecutionVersion {
		case NORMAL_ASYNC:
			p_uniq.SpawnThenTransition(re)
		case NORMAL_SYNC:
			p_uniq.SpawnThenTransition(re)
		case NON_POLARIZED_SYNC:
			p_uniq.SpawnThenTransitionNP(re)
		}
	}
}

func (re *RuntimeEnvironment) StopMonitor() ([]Process, []MonitorRulesLog) {
	re.monitor.stopMonitor()
	return re.monitor.deadProcesses, re.monitor.rulesLog
}

func (re *RuntimeEnvironment) HeartbeatReceiver(timeout time.Duration, cancel context.CancelFunc) {
	fullTimeout := re.Delay + timeout

	for {
		select {
		case <-time.After(fullTimeout):
			// todo do proper restart
			cancel()
			return
		case <-re.heartbeat:
			// Restart timer
		}
	}
}

func (re *RuntimeEnvironment) Ctx() context.Context {
	return re.ctx
}

func (re *RuntimeEnvironment) ErrorChan() chan error {
	return re.errorChan
}

type Message struct {
	Rule Rule
	// Possible payload types, depending on the rule
	Channel1         Name
	Channel2         Name
	Providers        []Name
	ContinuationBody Form
	Label            Label
}

type Rule int

const (
	// These can happen when a process is 'interactive' by sending messages
	SND Rule = iota // uses Channel1 and Channel2 of the Message
	RCV             // uses ContinuationBody, Channel1 and Channel2 of the Message
	CLS             // does not use message payloads
	CST             // uses Channel1 (continuation_c) of the Message
	SHF             // uses Channel1 of the Message
	SEL             // uses Channel1 and Label of the Message
	CSE             // uses Channel1 and Label of the Message
	DROP

	// These can happen when a process is 'interactive' by transitioning internally
	CUT
	SPLIT
	CALL // (maybe can happen when interactive or not)

	// When a process is 'non-interactive', either the FWD or DUP rules take place
	// Special rules for control messages
	// FWD // uses Channel1 of the ControlMessage struct
	DUP

	// Other actions
	FWD
	FWD_DROP
	PRINT
)

var RuleString = map[Rule]string{
	SND:  "SND",
	RCV:  "RCV",
	CLS:  "CLS",
	CST:  "CST",
	SHF:  "SHF",
	SEL:  "SEL",
	CSE:  "CSE",
	DROP: "DROP",

	CUT:   "CUT",
	CALL:  "CALL",
	SPLIT: "SPLIT",

	DUP: "DUP",

	FWD:      "FWD",
	FWD_DROP: "FWD_DROP",
	PRINT:    "PRINT",
}

type ControlMessage struct {
	Action Action
	// Possible payload types, depending on the action
	Providers []Name
	Body      Form
	Shape     Shape
}

type Action int

const (
	FWD_ACTION      Action = 100 // uses Channel1 of the ControlMessage struct
	FWD_ACTION_DROP Action = 101 // Forward with no providers (i.e. perform drop action)
)

type NameInitialization struct {
	old Name
	new Name
}

////// Logger details

type LogLevel int

const (
	LOGRULE        LogLevel = iota // rule started/finished
	LOGINFO                        // information
	LOGRULEDETAILS                 // rule while processing
	LOGPROCESSING                  // process info
	LOGMONITOR
)

// Similar to Println
func (re *RuntimeEnvironment) log(level LogLevel, message string) {
	if slices.Contains(re.GlobalEnvironment.LogLevels, level) {
		fmt.Println(message)
	}
}

// Similar to Printf
func (re *RuntimeEnvironment) logf(level LogLevel, message string, args ...interface{}) {
	if slices.Contains(re.GlobalEnvironment.LogLevels, level) {
		fmt.Printf(message, args...)
	}
}

// Color: Red ("\033[31m", "\033[101m"), Green, Yellow, Blue, Purple, Cyan, Gray
var colors = []string{"\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
var colorsHl = []string{"\033[102m", "\033[103m", "\033[104m", "\033[105m", "\033[106m", "\033[107m"}

const colorRed = "\033[31m"
const colorHlRed = "\033[101m"
const colorsLen = 5 // avoiding gray coz it looks like white and red as it resembles error messages
const resetColor = "\033[0m"

// Similar to Println
func (re *RuntimeEnvironment) logProcess(level LogLevel, process *Process, message string) {
	if slices.Contains(re.GlobalEnvironment.LogLevels, level) {
		var buffer bytes.Buffer

		if re.Color {
			colorIndex := 0
			if len(process.Providers) > 0 {
				colorIndex = int(process.Providers[0].ChannelID) % colorsLen
			}

			buffer.WriteString(colors[colorIndex])
		}

		buffer.WriteString(process.OutlineString())
		buffer.WriteString(": ")
		buffer.WriteString(message)

		if re.Color {
			buffer.WriteString(resetColor)
		}

		buffer.WriteString("\n")

		fmt.Print(buffer.String())

	}
}

// Similar to Printf
func (re *RuntimeEnvironment) logProcessf(level LogLevel, process *Process, message string, args ...interface{}) {
	if slices.Contains(re.GlobalEnvironment.LogLevels, level) {
		var buffer bytes.Buffer

		if re.Color {
			colorIndex := 0
			if len(process.Providers) > 0 {
				colorIndex = int(process.Providers[0].ChannelID) % colorsLen
			}
			buffer.WriteString(colors[colorIndex])
		}
		buffer.WriteString(process.OutlineString())
		buffer.WriteString(": ")
		buffer.WriteString(fmt.Sprintf(message, args...))
		if re.Color {
			buffer.WriteString(resetColor)
		}

		fmt.Print(buffer.String())
	}
}

// // Similar to logProcessf but adds highlighted text
// func (re *RuntimeEnvironment) logProcessHighlight(level LogLevel, process *Process, message string) {
// 	if slices.Contains(re.LogLevels, level) {
// 		colorIndex := 0
// 		if len(process.Providers) > 0 {
// 			colorIndex = int(process.Providers[0].ChannelID) % colorsLen
// 		}

// 		var buf bytes.Buffer
// 		buf.WriteString(colorsHl[colorIndex])
// 		buf.WriteString(process.OutlineString())
// 		buf.WriteString(": ")
// 		buf.WriteString(message)
// 		buf.WriteString(resetColor)

// 		fmt.Print(buf.String())
// 	}
// }

// func (re *RuntimeEnvironment) logProcessHighlightf(level LogLevel, process *Process, message string, args ...interface{}) {
// 	if slices.Contains(re.LogLevels, level) {

// 		data := append([]interface{}{process.OutlineString()}, args...)

// 		if re.Color {
// 			// todo fix: remove /n (if needed) from message and add it at the end

// 			colorIndex := 0
// 			if len(process.Providers) > 0 {
// 				colorIndex = int(process.Providers[0].ChannelID) % colorsLen
// 			}
// 			var buf bytes.Buffer
// 			// buf.WriteString(colors[colorIndex])
// 			buf.WriteString(colorsHl[colorIndex])
// 			buf.WriteString(fmt.Sprintf("%s: "+message, data...))
// 			buf.WriteString(resetColor)

// 			fmt.Print(buf.String())
// 		} else {
// 			fmt.Printf("%s: "+message, data...)
// 		}
// 		// fmt.Printf("%s", resetColor)
// 	}
// }

func (re *RuntimeEnvironment) logMonitorf(message string, args ...interface{}) {
	if slices.Contains(re.GlobalEnvironment.LogLevels, LOGMONITOR) {

		var buf bytes.Buffer
		if re.monitor.re.Color {
			buf.WriteString(colorsHl[0])
		}
		buf.WriteString("[monitor]:")
		if re.monitor.re.Color {
			buf.WriteString(resetColor)
		}
		buf.WriteString(" ")
		buf.WriteString(fmt.Sprintf(message, args...))

		fmt.Print(buf.String())
	}
}

func (re *RuntimeEnvironment) error(process *Process, message string) {
	var buffer bytes.Buffer

	if re.Color {
		buffer.WriteString(colorHlRed)
		buffer.WriteString("Error in ")
		buffer.WriteString(process.OutlineString())
		buffer.WriteString(resetColor)
		buffer.WriteString("\n")
		buffer.WriteString(colorRed)
		buffer.WriteString(message)
		buffer.WriteString(resetColor)
		buffer.WriteString("\n")
	} else {
		fmt.Fprintf(&buffer, "Error %s: "+message+"\n", process.OutlineString())
	}

	// fmt.Println(buffer.String())
	// todo change to use ctx Err
	panic(buffer.String())
}

// Similar to Printf
func (re *RuntimeEnvironment) errorf(process *Process, message string, args ...interface{}) {
	data := append([]interface{}{process.OutlineString()}, args...)

	var buffer bytes.Buffer

	if re.Color {
		buffer.WriteString(colorHlRed)
		buffer.WriteString("Error in ")
		buffer.WriteString(process.OutlineString())
		buffer.WriteString("\n")
		buffer.WriteString(resetColor)
		buffer.WriteString(colorRed)
		buffer.WriteString(fmt.Sprintf(message, args...))
		buffer.WriteString(resetColor)
	} else {
		fmt.Fprintf(&buffer, "%s: "+message, data...)
	}

	// fmt.Println(buffer.String())
	panic(buffer.String())
}

package process

import (
	"fmt"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slices"
)

// A RuntimeEnvironmentOption sets an option on a RuntimeEnvironment.
type RuntimeEnvironmentOption func(*RuntimeEnvironment)

// RuntimeEnvironment provides the instance instance for processes to be initialized and executed
type RuntimeEnvironment struct {
	ProcessCount int

	// Logging levels
	logLevels []LogLevel
	// Debugging info
	debug bool
	// Keeps counter of the number of channels created
	debugChannelCounter uint64
	// Monitor info
	monitor *Monitor
	// Controller info
	controller *Controller
}

// Entry point for execution
func InitializeProcesses(processes []Process) {
	l := []LogLevel{
		LOGINFO,
		LOGRULE,
		LOGRULEDETAILS,
		LOGPROCESSING,
	}

	re := &RuntimeEnvironment{ProcessCount: 0, debugChannelCounter: 0, debug: true, logLevels: l}

	re.logf(LOGINFO, "Initializing %d processes\n", len(processes))

	channels := re.CreateChannelForEachProcess(processes)

	re.SubstituteNameInitialization(processes, channels)

	// re.log(LOGINFO, "After Substitutions")
	// for _, p := range processes {
	// 	re.log(LOGINFO, p.String())
	// }

	if re.debug {
		started := make(chan bool)

		re.InitializeMonitor(started)
		re.InitializeController(started)

		// Ensure that both servers are running
		<-started
		<-started
	}

	re.StartTransitions(processes)
}

// Create the initial channels required. E.g. for a process prc[c1], a channel with Ident: c1 is created
func (re *RuntimeEnvironment) CreateChannelForEachProcess(processes []Process) []NameInitialization {

	var channels []NameInitialization

	for i := 0; i < len(processes); i++ {
		old_provider := processes[i].Provider
		new_provider := re.CreateFreshChannel(old_provider.Ident)

		// Set new channel as the providing channel for the process
		processes[i].Provider = new_provider
		channels = append(channels, NameInitialization{old_provider, new_provider})
	}

	return channels
}

// Create new channel
func (re *RuntimeEnvironment) CreateFreshChannel(ident string) Name {
	// The channel ID is used for debugging
	atomic.AddUint64(&re.debugChannelCounter, 1)

	// Create new channel and assign a name to it
	c := make(chan Message)
	return Name{Ident: ident, Channel: c, ChannelID: re.debugChannelCounter}
}

func (re *RuntimeEnvironment) InitializeMonitor(started chan bool) {
	// Declare new monitor
	re.monitor = NewMonitor(re)

	// Start monitor on new thread
	go re.monitor.StartMonitor(started)
}

func (re *RuntimeEnvironment) InitializeController(started chan bool) {
	// Declare new controller
	re.controller = NewController(re)

	// Start controller on new thread
	go re.controller.StartController(started)
}

// Used after initialization to substitute known names to the actual channel
func (re *RuntimeEnvironment) SubstituteNameInitialization(processes []Process, channels []NameInitialization) {
	for i := 0; i < len(processes); i++ {
		for _, c := range channels {
			// Substitute all free names in a body
			processes[i].Body.Substitute(c.old, c.new)
		}
	}
}

func (re *RuntimeEnvironment) StartTransitions(processes []Process) {
	for _, p := range processes {
		p_uniq := p
		p_uniq.Transition(re)
	}
	time.Sleep(1 * time.Second)

	re.logf(LOGINFO, "End process count: %d\n", re.ProcessCount)
}

type Message struct {
	Rule Rule
	// Possible payload types, depending on the rule
	Channel1         Name
	Channel2         Name
	ContinuationBody Form
	Label            Label
}

type Rule int

const (
	SND Rule = iota // Channel1 and Channel2
	RCV             // ContinuationBody, Channel1 and Channel2
	FWD             // Channel1
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
	LOGERROR
)

// Similar to Println
func (re *RuntimeEnvironment) log(level LogLevel, message string) {
	if slices.Contains(re.logLevels, level) {
		fmt.Println(message)
	}
}

// Similar to Printf
func (re *RuntimeEnvironment) logf(level LogLevel, message string, args ...interface{}) {
	if slices.Contains(re.logLevels, level) {
		fmt.Printf(message, args...)
	}
}

// Similar to Println
func (re *RuntimeEnvironment) logProcess(level LogLevel, process *Process, message string) {
	if slices.Contains(re.logLevels, level) {
		fmt.Printf("%s: "+message+"\n", process.OutlineString())
	}
}

// Similar to Printf
func (re *RuntimeEnvironment) logProcessf(level LogLevel, process *Process, message string, args ...interface{}) {
	if slices.Contains(re.logLevels, level) {
		data := append([]interface{}{process.OutlineString()}, args...)

		fmt.Printf("%s: "+message, data...)
	}
}

package process

import (
	"fmt"
	"sync/atomic"
	"time"
)

// A RuntimeEnvironmentOption sets an option on a RuntimeEnvironment.
type RuntimeEnvironmentOption func(*RuntimeEnvironment)

// RuntimeEnvironment provides the instance instance for processes to be initialized and executed
type RuntimeEnvironment struct {
	ProcessCount int

	// Debugging info
	debug bool
	// Keeps counter of the number of channels created
	debugChannelCounter uint64
	// Monitor info
	monitor *Monitor
	// Controller info
	controller *Controller
}

func InitializeProcesses(processes []Process) {
	fmt.Printf("Initializing %d processes\n", len(processes))

	re := &RuntimeEnvironment{ProcessCount: 0, debugChannelCounter: 0, debug: true}

	channels := re.CreateChannelForEachProcess(processes)

	// fmt.Println(len(channels))

	re.SubstituteNameInitialization(processes, channels)

	fmt.Println("After Substitutions")
	for _, p := range processes {
		fmt.Print(p.String())
		fmt.Println(": ")
	}

	re.InitializeMonitor()
	re.InitializeController()

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

func (re *RuntimeEnvironment) InitializeMonitor() {
	// Declare new monitor
	re.monitor = NewMonitor()

	// Start monitor on new thread
	go re.monitor.StartMonitor()
}

func (re *RuntimeEnvironment) InitializeController() {
	// Declare new controller
	re.controller = NewController()

	// Start controller on new thread
	go re.controller.StartController()
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

	fmt.Print("End process count: ")
	fmt.Println(re.ProcessCount)
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
)

type NameInitialization struct {
	old Name
	new Name
}

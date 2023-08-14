package process

import (
	"bytes"
	"fmt"
	"sync/atomic"

	"golang.org/x/exp/slices"
)

// A RuntimeEnvironmentOption sets an option on a RuntimeEnvironment.
type RuntimeEnvironmentOption func(*RuntimeEnvironment)

// RuntimeEnvironment provides the instance instance for processes to be initialized and executed
type RuntimeEnvironment struct {
	// Keeps count of how many processes were spawned (only for debug info)
	ProcessCount uint64

	// Logging levels
	logLevels []LogLevel
	// Debugging info
	debug bool
	// Colored output
	color bool
	// Keeps counter of the number of channels created
	debugChannelCounter uint64
	// Monitor info
	monitor *Monitor
	// Controller info
	controller *Controller
}

func NewRuntimeEnvironment(l []LogLevel, debug, coloredOutput bool) *RuntimeEnvironment {
	return &RuntimeEnvironment{ProcessCount: 0, debugChannelCounter: 0, debug: true, color: true, logLevels: l}
}

// Entry point for execution
func InitializeProcesses(processes []Process) {
	l := []LogLevel{
		LOGERROR,
		LOGINFO,
		LOGPROCESSING,
		LOGRULE,
		// LOGRULEDETAILS,
		LOGMONITOR,
	}

	re := &RuntimeEnvironment{ProcessCount: 0, debugChannelCounter: 0, debug: true, color: true, logLevels: l}

	re.logf(LOGINFO, "Initializing %d processes\n", len(processes))

	channels := re.CreateChannelForEachProcess(processes)

	re.SubstituteNameInitialization(processes, channels)

	if re.debug {
		started := make(chan bool)

		re.InitializeMonitor(started)
		re.InitializeController(started)

		// Ensure that both servers are running
		<-started
		<-started
	}

	re.StartTransitions(processes)

	re.WaitForMonitorToFinish()

	// time.Sleep(5 * time.Second)

	re.logf(LOGINFO, "End process count: %d\n", re.ProcessCount)
}

func (re *RuntimeEnvironment) WaitForMonitorToFinish() ([]Process, []MonitorRulesLog) {
	<-re.monitor.monitorFinished
	return re.monitor.deadProcesses, re.monitor.rulesLog
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
	mChan := make(chan Message)
	pmChan := make(chan PriorityMessage)
	return Name{Ident: ident, Channel: mChan, ChannelID: re.debugChannelCounter, PriorityChannel: pmChan, IsSelf: false}
}

func (re *RuntimeEnvironment) InitializeMonitor(started chan bool) {
	// Declare new monitor
	re.monitor = NewMonitor(re)

	// Start monitor on new thread
	go re.monitor.startMonitor(started)
}

func (re *RuntimeEnvironment) InitializeController(started chan bool) {
	// Declare new controller
	re.controller = NewController(re)

	// Start controller on new thread
	go re.controller.startController(started)
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
	// These can happen when a process is 'interactive' by sending messages
	SND Rule = iota // uses Channel1 and Channel2 of the Message struct
	RCV             // uses ContinuationBody, Channel1 and Channel2 of the Message struct

	// These can happen when a process is 'interactive' by transitioning internally
	CUT
	SPLIT
	CALL // (maybe can happen when interactive or not)

	// When a process is 'non-interactive', either the FWD or DUP rules take place
	// Special rules for priority messages
	FWD       // uses Channel1 of the PriorityMessage struct
	FWD_REPLY // uses Body, Shape of the PriorityMessage struct
	DUP
)

var RuleString = map[Rule]string{
	SND:   "SND",
	RCV:   "RCV",
	CUT:   "CUT",
	CALL:  "CALL",
	SPLIT: "SPLIT",

	FWD:       "FWD",
	FWD_REPLY: "FWD_REPLY",
	DUP:       "DUP",
}

type PriorityMessage struct {
	Action Rule
	// Possible payload types, depending on the action
	Channels []Name
	Body     Form
	Shape    Shape
}

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
	LOGMONITOR
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

// Color: Red ("\033[31m", "\033[101m"), Green, Yellow, Blue, Purple, Cyan, Gray
var colors = []string{"\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
var colorsHl = []string{"\033[102m", "\033[103m", "\033[104m", "\033[105m", "\033[106m", "\033[107m"}

const colorsLen = 5 // avoiding gray coz it looks like white and red as it resembles error messages

var resetColor = "\033[0m"

// Similar to Println
func (re *RuntimeEnvironment) logProcess(level LogLevel, process *Process, message string) {
	if slices.Contains(re.logLevels, level) {
		if re.color {
			colorIndex := int(process.Provider.ChannelID) % colorsLen
			fmt.Printf("%s%s: "+message+"\n%s", colors[colorIndex], process.OutlineString(), resetColor)
		} else {
			fmt.Printf("%s: "+message+"\n", process.OutlineString())
		}
	}
}

// Similar to Printf
func (re *RuntimeEnvironment) logProcessf(level LogLevel, process *Process, message string, args ...interface{}) {
	if slices.Contains(re.logLevels, level) {
		data := append([]interface{}{process.OutlineString()}, args...)

		if re.color {
			colorIndex := int(process.Provider.ChannelID) % colorsLen
			var buf bytes.Buffer
			buf.WriteString(colors[colorIndex])
			buf.WriteString(fmt.Sprintf("%s: "+message, data...))
			buf.WriteString(resetColor)

			fmt.Print(buf.String())
		} else {
			fmt.Printf("%s: "+message, data...)
		}
	}
}

// Similar to logProcessf but adds highlighted text
func (re *RuntimeEnvironment) logProcessHighlight(level LogLevel, process *Process, message string) {
	if slices.Contains(re.logLevels, level) {

		colorIndex := int(process.Provider.ChannelID) % colorsLen

		var buf bytes.Buffer
		buf.WriteString(colorsHl[colorIndex])
		buf.WriteString(process.OutlineString())
		buf.WriteString(": ")
		buf.WriteString(message)
		buf.WriteString(resetColor)

		fmt.Print(buf.String())
	}
}

func (re *RuntimeEnvironment) logProcessHighlightf(level LogLevel, process *Process, message string, args ...interface{}) {
	if slices.Contains(re.logLevels, level) {

		data := append([]interface{}{process.OutlineString()}, args...)

		if re.color {
			// todo fix: remove /n (if needed) from message and add it at the end

			colorIndex := int(process.Provider.ChannelID) % colorsLen
			var buf bytes.Buffer
			// buf.WriteString(colors[colorIndex])
			buf.WriteString(colorsHl[colorIndex])
			buf.WriteString(fmt.Sprintf("%s: "+message, data...))
			buf.WriteString(resetColor)

			fmt.Print(buf.String())
		} else {
			fmt.Printf("%s: "+message, data...)
		}
		// fmt.Printf("%s", resetColor)
	}
}

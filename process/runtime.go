package process

import (
	"fmt"
	"time"
)

type RuntimeEnvironment struct {
	processCount int
}

var runtimeEnvironment RuntimeEnvironment

func InitializeProcesses(processes []Process) {
	fmt.Printf("Initializing %d processes\n", len(processes))

	channels := CreateChannelForEachProcess(processes)

	// fmt.Println(len(channels))

	SubstituteNameInitialization(processes, channels)

	fmt.Println("After Substitutions")
	for _, p := range processes {
		fmt.Print(p.String())
		fmt.Println(": ")
	}

	StartTransitions(processes)
}

// Create the initial channels required. E.g. for a process prc[c1], a channel with Ident: c1 is created
func CreateChannelForEachProcess(processes []Process) []NameInitialization {

	var channels []NameInitialization

	for i := 0; i < len(processes); i++ {
		old_provider := processes[i].Provider
		c := make(chan Message)
		new_provider := Name{Ident: old_provider.Ident, Channel: c}

		// Set new channel as the providing channel for the process
		processes[i].Provider = new_provider
		channels = append(channels, NameInitialization{old_provider, new_provider})
	}

	return channels
}

// Used after initialization to substitute known names to the actual channel
func SubstituteNameInitialization(processes []Process, channels []NameInitialization) {
	for i := 0; i < len(processes); i++ {
		for _, c := range channels {
			// Substitute all free names in a body
			processes[i].Body.Substitute(c.old, c.new)
		}
	}
}

func StartTransitions(processes []Process) {
	for _, p := range processes {
		p.Transition()
	}
	time.Sleep(1 * time.Second)

	fmt.Print("End process count: ")
	fmt.Println(runtimeEnvironment.processCount)
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

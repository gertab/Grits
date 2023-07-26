package process

import "fmt"

type NameInitialization struct {
	old Name
	new Name
}

func InitializeProcesses(processes []Process) {
	fmt.Printf("Initializing %d processes\n", len(processes))

	channels := CreateChannelForEachProcess(processes)

	SubstituteNameInitialization(processes, channels)

	fmt.Println("After Substitutions")
	for _, p := range processes {
		fmt.Print(p.String())
		fmt.Println(": ")
	}
}

// Create the initial channels required. E.g. for a process prc[c1], the channel c1 is created
func CreateChannelForEachProcess(processes []Process) []NameInitialization {

	var channels []NameInitialization

	for _, p := range processes {
		fmt.Println(p.String())
		old_provider := p.Provider
		c := make(chan Message)
		new_provider := Name{Ident: old_provider.Ident, Channel: c}

		channels = append(channels, NameInitialization{old_provider, new_provider})
	}

	return channels
}

// Used after initialization to substitute known names to the actual channel
func SubstituteNameInitialization(processes []Process, channels []NameInitialization) {
	// fmt.Println("SubstituteNameInitialization")
	for _, p := range processes {
		// fmt.Print(p.String())
		// fmt.Println(": ")

		for _, c := range channels {
			// fmt.Print(c.old.String())
			// fmt.Print(" -> ")
			// fmt.Println(c.new.String())

			p.Body.Substitute(c.old, c.new)
		}
	}
}

type Message interface {
	// Payload()
	// Type()
}

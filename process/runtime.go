package process

import "fmt"

func InitializeProcesses(processes []Process) {
	fmt.Printf("Initializing %d processes\n", len(processes))

	for _, p := range processes {
		fmt.Println(p.String())
	}
}

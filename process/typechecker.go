package process

import (
	"fmt"
	"log"
	"phi/types"
)

func Typecheck(processes []*Process, globalEnv *GlobalEnvironment) error {
	errorChan := make(chan error)
	doneChan := make(chan bool)

	// Running in a separate process allows us to break the typechecking part as soon as the first
	// error is found
	go typecheckProcesses(processes, globalEnv, errorChan, doneChan)

	select {
	case err := <-errorChan:
		log.Fatal(err)
	case <-doneChan:
		fmt.Println("Typecheck successful")
	}

	i := types.NewLabelType("hello")

	k := types.NewReceiveType(i, types.NewUnitType())

	j := types.CopyType(k)

	i.Label = "iiii"

	fmt.Println(j)
	fmt.Println(k)

	// k.

	return nil
}

func typecheckProcesses(processes []*Process, globalEnv *GlobalEnvironment, errorChan chan error, doneChan chan bool) {
	defer func() {
		doneChan <- true
	}()

	// Start with some preliminary check on the types
	err := lightweightChecks(processes, globalEnv)
	if err != nil {
		errorChan <- err
	}

	// We can initiate the more heavyweight typechecking

	fmt.Println("ok")
}

// Perform some preliminary checks about the type type definitions
// Ensures that types only referred to existing labelled types
func lightweightChecks(processes []*Process, globalEnv *GlobalEnvironment) error {

	// First check the labelled types (i.e. type A = ...)
	err := types.SanityChecksTypeDefinitions(*globalEnv.Types)

	if err != nil {
		return err
	}

	// Check the types for the function declarations (i.e. let f() : A = ...)
	var all_types []types.SessionType
	for _, i := range *globalEnv.FunctionDefinitions {
		if i.Type != nil {
			all_types = append(all_types, i.Type)
		}

		for _, i := range i.Parameters {
			if i.Type != nil {
				all_types = append(all_types, i.Type)
			}
		}
	}

	// Check the types for the processes (i.e. prc[a] : A = ...)
	for _, i := range processes {
		if i.Type != nil {
			all_types = append(all_types, i.Type)
		}
	}

	err = types.SanityChecksType(all_types, *globalEnv.Types)
	if err != nil {
		return err
	}

	return nil
}

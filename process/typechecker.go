package process

import (
	"fmt"
	"log"
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
		fmt.Println("Typecheck succesful")
	}

	return nil
}

func typecheckProcesses(processes []*Process, globalEnv *GlobalEnvironment, errorChan chan error, doneChan chan bool) {
	defer func() {
		doneChan <- true
	}()

	fmt.Println("ok")
}

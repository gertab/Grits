package process

import "fmt"

type Monitor struct {
	i int
	// Monitor receives info from processes on monitorChan
	monitorChan chan Process
	// Processes report error to monitor on monitorChan
	errorChan chan error
}

func NewMonitor() *Monitor {

	monitorChan := make(chan Process)
	errorChan := make(chan error)

	return &Monitor{i: 0, monitorChan: monitorChan, errorChan: errorChan}
}

func (m *Monitor) StartMonitor() {
	fmt.Println("Monitor alive, waiting to receive...")

}

type Controller struct {
	i int
	// Controller receives new action permission requests on this channel
	controllerNewActionChan chan int
}

func NewController() *Controller {

	controllerChan := make(chan int)

	return &Controller{i: 0, controllerNewActionChan: controllerChan}
}

func (m *Controller) StartController() {
	fmt.Println("Controller alive, waiting to receive...")

}

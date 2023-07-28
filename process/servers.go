package process

type Monitor struct {
	i int
	// Monitor receives info from processes on monitorChan
	monitorChan chan Process
	// Processes report error to monitor on monitorChan
	errorChan chan error
	// Runtime environment contains log info
	re *RuntimeEnvironment
}

func NewMonitor(re *RuntimeEnvironment) *Monitor {

	monitorChan := make(chan Process)
	errorChan := make(chan error)

	return &Monitor{i: 0, monitorChan: monitorChan, errorChan: errorChan, re: re}
}

func (m *Monitor) StartMonitor(started chan bool) {
	m.re.log(LOGINFO, "Monitor alive, waiting to receive...")

	started <- true

}

type Controller struct {
	i int
	// Controller receives new action permission requests on this channel
	controllerNewActionChan chan int
	// Runtime environment contains log info
	re *RuntimeEnvironment
}

func NewController(re *RuntimeEnvironment) *Controller {

	controllerChan := make(chan int)

	return &Controller{i: 0, controllerNewActionChan: controllerChan, re: re}
}

func (m *Controller) StartController(started chan bool) {
	m.re.log(LOGINFO, "Controller alive, waiting to receive...")

	started <- true
}

package process

import "fmt"

type Monitor struct {
	i int
	// Monitor receives info from processes on monitorChan
	monitorChan (chan MonitorUpdate)
	// Processes report error to monitor on monitorChan
	errorChan chan error
	// Runtime environment contains log info
	re *RuntimeEnvironment
}

type MonitorUpdate struct {
	process Process
	rule    Rule
	isDead  bool
}

func NewMonitor(re *RuntimeEnvironment) *Monitor {

	// todo make these buffered
	monitorChan := make(chan MonitorUpdate)
	errorChan := make(chan error)

	return &Monitor{i: 0, monitorChan: monitorChan, errorChan: errorChan, re: re}
}

func (m *Monitor) startMonitor(started chan bool) {
	m.re.log(LOGINFO, "Monitor alive, waiting to receive...")

	started <- true

	m.monitorLoop()

}

func (m *Monitor) monitorLoop() {
	select {
	case processUpdate := <-m.monitorChan:
		if processUpdate.isDead {
			fmt.Println("[monitor] Process died")
		} else {
			fmt.Println("[monitor] ", processUpdate.rule, processUpdate.process.String())
		}

	case error := <-m.errorChan:
		fmt.Println(error)
	}
}

// Monitor: User API
func (m *Monitor) MonitorRuleFinished(process *Process, rule Rule) {

	body := CopyForm(process.Body)
	provider := process.Provider
	shape := process.Shape

	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, shape, nil), rule: rule}
}

func (m *Monitor) MonitorProcessTerminated(process *Process) {

	// body := CopyForm(process.Body)
	provider := process.Provider
	shape := process.Shape

	m.monitorChan <- MonitorUpdate{process: *NewProcess(nil, provider, shape, nil), isDead: true}
}

// Controller
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

func (m *Controller) startController(started chan bool) {
	m.re.log(LOGINFO, "Controller alive, waiting to receive...")

	started <- true
}

package process

import (
	"fmt"
	"time"
)

type Monitor struct {
	i int
	// Monitor receives info from processes on monitorChan
	monitorChan (chan MonitorUpdate)
	// Processes report error to monitor on monitorChan
	errorChan chan error
	// Runtime environment contains log info
	re *RuntimeEnvironment
	// Keeps a log of all the rules that took place
	rulesLog      []monitorRulesLog
	deadProcesses []Process
	// Inactive after _ milliseconds
	inactiveTimer   time.Duration
	monitorFinished chan bool
}

type MonitorUpdate struct {
	process     Process
	rule        Rule
	isDead      bool
	stopMonitor bool
}

type monitorRulesLog struct {
	process Process
	rule    Rule
}

func NewMonitor(re *RuntimeEnvironment) *Monitor {

	// todo make these buffered
	monitorChan := make(chan MonitorUpdate)
	errorChan := make(chan error)

	return &Monitor{i: 0, monitorChan: monitorChan, errorChan: errorChan, re: re, inactiveTimer: 50 * time.Millisecond}
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
			// Process is terminated
			fmt.Printf("[monitor] Process %s died\n", processUpdate.process.Provider.String())
			m.deadProcesses = append(m.deadProcesses, processUpdate.process)
		} else if processUpdate.stopMonitor {
			// Stops monitoring
			return
		} else {
			// Process finished rule
			fmt.Println("[monitor] finished", ruleString[processUpdate.rule], processUpdate.process.String())
			m.rulesLog = append(m.rulesLog, monitorRulesLog{process: processUpdate.process, rule: processUpdate.rule})
		}

	case error := <-m.errorChan:
		fmt.Println(error)

	case <-time.After(m.inactiveTimer):
		fmt.Println("Timer terminated after", m.inactiveTimer)
		// m.monitorFinished <- true
		return
	}

	m.monitorLoop()
}

func (m *Monitor) GetRulesLog() []monitorRulesLog {
	return m.rulesLog
}

// Monitor: User API
func (m *Monitor) MonitorRuleFinished(process *Process, rule Rule) {

	body := CopyForm(process.Body)
	provider := process.Provider
	shape := process.Shape

	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, shape, nil), rule: rule, isDead: false, stopMonitor: false}
}

func (m *Monitor) MonitorProcessTerminated(process *Process) {
	// body := CopyForm(process.Body)
	provider := process.Provider
	shape := process.Shape

	m.monitorChan <- MonitorUpdate{process: *NewProcess(nil, provider, shape, nil), isDead: true, stopMonitor: false}
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

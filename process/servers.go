package process

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/exp/slices"
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
	rulesLog      []MonitorRulesLog
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

type MonitorRulesLog struct {
	Process Process
	Rule    Rule
}

func NewMonitor(re *RuntimeEnvironment) *Monitor {

	// todo make these buffered
	monitorChan := make(chan MonitorUpdate)
	errorChan := make(chan error)
	monitorFinishedChan := make(chan bool)

	return &Monitor{i: 0, monitorChan: monitorChan, errorChan: errorChan, monitorFinished: monitorFinishedChan, re: re, inactiveTimer: 200 * time.Millisecond}
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
			m.re.logMonitorf("Process %s died\n", processUpdate.process.Provider.String())
			m.deadProcesses = append(m.deadProcesses, processUpdate.process)
		} else if processUpdate.stopMonitor {
			// Stops monitoring
			return
		} else {
			// Process finished rule
			m.re.logMonitorf("Finished %s, %s\n", RuleString[processUpdate.rule], processUpdate.process.String())
			m.rulesLog = append(m.rulesLog, MonitorRulesLog{Process: processUpdate.process, Rule: processUpdate.rule})
		}

	case error := <-m.errorChan:
		fmt.Println(error)

	case <-time.After(m.inactiveTimer):
		m.re.logMonitorf("Monitor inactive, terminating\n")
		m.monitorFinished <- true
		return
	}

	m.monitorLoop()
}

func (m *Monitor) GetRulesLog() []MonitorRulesLog {
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

func (re *RuntimeEnvironment) logMonitorf(message string, args ...interface{}) {
	if slices.Contains(re.logLevels, LOGMONITOR) {

		data := append([]interface{}{"[monitor]"}, args...)

		colorIndex := 1

		var buf bytes.Buffer
		buf.WriteString(colorsHl[colorIndex])
		buf.WriteString(fmt.Sprintf("%s:\033[0m "+message, data...))
		buf.WriteString(resetColor)

		fmt.Print(buf.String())
	}
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

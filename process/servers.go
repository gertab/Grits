package process

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/exp/slices"
)

type Monitor struct {
	i int
	// Monitor receives info from processes on monitorChan
	monitorChan chan MonitorUpdate
	// Processes report errors to monitor on monitorChan
	errorChan chan error
	// Processes report errors to monitor on monitorChan
	stopMonitorChan chan bool
	// Send updates to subscriber chan
	subscriber *SubscriberInfo
	// Runtime environment contains log info
	re *RuntimeEnvironment
	// Keeps a log of all the rules that took place
	rulesLog []MonitorRulesLog
	// Keeps a log of all the dead/present processes
	deadProcesses []Process
	// presentProcesses     []Process
	processID            int
	providersToProcessID map[Name]int
	processIDToProcess   map[int]*Process
	// Inactive after _ milliseconds
	inactiveTimer   time.Duration
	monitorFinished chan bool
}

type MonitorUpdate struct {
	process                Process
	rule                   Rule
	isDead                 bool
	isRuleDone             bool
	isRuleDoneBeforeRename bool
	updatedProcess         bool
	newProcess             bool
}

type MonitorRulesLog struct {
	Process Process
	Rule    Rule
}

func NewMonitor(re *RuntimeEnvironment, subscriberInfo *SubscriberInfo) *Monitor {

	// todo make these buffered
	monitorChan := make(chan MonitorUpdate)
	errorChan := make(chan error)
	stopMonitorChan := make(chan bool)
	monitorFinishedChan := make(chan bool)

	// Set the maps to store the processes
	providersToProcessID := make(map[Name]int)
	processIDToProcess := make(map[int]*Process)

	return &Monitor{i: 0, monitorChan: monitorChan, errorChan: errorChan, stopMonitorChan: stopMonitorChan, subscriber: subscriberInfo, monitorFinished: monitorFinishedChan, re: re, inactiveTimer: 300 * time.Millisecond, providersToProcessID: providersToProcessID, processIDToProcess: processIDToProcess, processID: 0}
}

func (m *Monitor) startMonitor(started chan bool) {
	m.re.log(LOGINFO, "Monitor alive, waiting to receive...")

	started <- true

	m.monitorLoop()
}

// func (m *Monitor) stopMonitor() {
// 	m.stopMonitorChan <- true
// }

func (m *Monitor) monitorLoop() {
	select {
	case processUpdate := <-m.monitorChan:
		if processUpdate.isDead {
			// Process is terminated
			m.re.logMonitorf("Process %s died\n", processUpdate.process.Providers[0].String())
			m.deadProcesses = append(m.deadProcesses, processUpdate.process)
			m.removeProcessFromList(processUpdate.process)
		} else if processUpdate.isRuleDone {
			// Process finished rule
			m.re.logMonitorf("Finished %s, %s\n", RuleString[processUpdate.rule], processUpdate.process.String())
			m.rulesLog = append(m.rulesLog, MonitorRulesLog{Process: processUpdate.process, Rule: processUpdate.rule})
			m.updateProcessToList(&processUpdate.process)
		} else if processUpdate.isRuleDoneBeforeRename {
			// Process finished rule but do not update the process list
			m.re.logMonitorf("Finished %s, %s\n", RuleString[processUpdate.rule], processUpdate.process.String())
			m.rulesLog = append(m.rulesLog, MonitorRulesLog{Process: processUpdate.process, Rule: processUpdate.rule})
			// m.updateProcessToList(&processUpdate.process)
		} else if processUpdate.updatedProcess {
			// Process updated form
			m.re.logMonitorf("Updated %s, %s\n", processUpdate.process.String())
			// m.rulesLog = append(m.rulesLog, MonitorRulesLog{Process: processUpdate.process, Rule: processUpdate.rule})
			m.updateProcessToList(&processUpdate.process)
		} else if processUpdate.newProcess {
			// Process finished rule
			m.re.logMonitorf("New process %s\n", processUpdate.process.String())
			// m.rulesLog = append(m.rulesLog, MonitorRulesLog{Process: processUpdate.process, Rule: processUpdate.rule})

			m.addProcessToList(&processUpdate.process)
		}

		// Send updated structure to subscriber
		m.updateSubscriber()

	case <-m.stopMonitorChan:
		m.re.logMonitorf("Monitor terminating\n")
		return

	case error := <-m.errorChan:
		fmt.Println(error)

	case <-time.After(m.inactiveTimer):
		m.re.logMonitorf("Monitor inactive, terminating\n")
		m.monitorFinished <- true
		return
	}

	m.monitorLoop()
}

func (m *Monitor) updateSubscriber() {
	// Send list of processes to the subscriber
	if m.subscriber != nil {

		v := make([]ProcessInfo, 0, len(m.processIDToProcess))

		for id, value := range m.processIDToProcess {
			providers := make([]string, 0, len(value.Providers))

			for _, provider := range value.Providers {
				providers = append(providers, provider.String())
			}

			v = append(v, ProcessInfo{ID: strconv.Itoa(id), Providers: providers, Body: value.String()})
			// v = append(v, ProcessInfo{ID: strconv.Itoa(id), Providers: providers, Body: value.Body.String()})
		}
		m.subscriber.ProcessesStringChan <- v
	}
}

func (m *Monitor) addProcessToList(process *Process) {
	m.processID++
	processID := m.processID
	for _, p := range process.Providers {
		m.providersToProcessID[p] = processID
	}
	m.processIDToProcess[processID] = process
}

func (m *Monitor) updateProcessToList(process *Process) {
	if len(process.Providers) > 0 {
		processID, providerExists := m.providersToProcessID[process.Providers[0]]

		if providerExists {
			// Already Linked To Process
			m.processIDToProcess[processID] = process
			return
		}
	}

	// m.processID++
	// processID := m.processID
	// for _, p := range process.Providers {
	// 	m.providersToProcessID[p] = processID
	// }
	// m.processIDToProcess[processID] = process
}

func (m *Monitor) removeProcessFromList(process Process) {
	for _, p := range process.Providers {
		// m.providersToProcessID[p] = processID
		i, providerExists := m.providersToProcessID[p]

		if providerExists {
			delete(m.providersToProcessID, p)
			delete(m.processIDToProcess, i)
		}
	}
}

func (m *Monitor) GetRulesLog() []MonitorRulesLog {
	return m.rulesLog
}

// Monitor: User API
func (m *Monitor) MonitorRuleFinished(process *Process, rule Rule) {
	body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape

	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, shape, nil), rule: rule, isRuleDone: true}
}

func (m *Monitor) MonitorRuleFinishedBeforeRenamed(process *Process, rule Rule) {
	// Since the provider will be stolen and used by a different process, the monitor should only note that a rule
	// has finished, without updating the list of current processes
	body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape

	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, shape, nil), rule: rule, isRuleDoneBeforeRename: true}
}

func (m *Monitor) MonitorNewProcess(process *Process) {
	body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape

	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, shape, nil), newProcess: true}
}

func (m *Monitor) MonitorProcessRenamed(process *Process) {
	body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape

	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, shape, nil), updatedProcess: true}
}
func (m *Monitor) MonitorProcessTerminated(process *Process) {
	// body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape

	m.monitorChan <- MonitorUpdate{process: *NewProcess(nil, provider, shape, nil), isDead: true}
}

func (m *Monitor) MonitorProcessForwarded(process *Process) {
	// body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape

	m.monitorChan <- MonitorUpdate{process: *NewProcess(nil, provider, shape, nil), isDead: true}
}

type SubscriberInfo struct {
	ProcessesStringChan chan []ProcessInfo
}

type ProcessInfo struct {
	ID        string   `json:"id"`
	Providers []string `json:"providers"`
	Body      string   `json:"body"`
}

func NewSubscriberInfo() *SubscriberInfo {
	processStringChan := make(chan []ProcessInfo)
	return &SubscriberInfo{ProcessesStringChan: processStringChan}
}

// func (m *Monitor) ForwardProcessFinished(process *Process) {
// 	body := CopyForm(process.Body)
// 	provider := process.Providers
// 	shape := process.Shape

// 	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, shape, nil), rule: FWD, isRuleDone: true}
// }

func (re *RuntimeEnvironment) logMonitorf(message string, args ...interface{}) {
	if slices.Contains(re.logLevels, LOGMONITOR) {

		data := append([]interface{}{"[monitor]"}, args...)

		colorIndex := 0

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

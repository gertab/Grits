package process

import (
	"fmt"
	"phi/types"
	"strconv"
	"sync"
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

	// Set the maps to store the processes
	providersToProcessID := make(map[Name]int)
	processIDToProcess := make(map[int]*Process)

	return &Monitor{
		i:                    0,
		monitorChan:          monitorChan,
		errorChan:            errorChan,
		stopMonitorChan:      stopMonitorChan,
		subscriber:           subscriberInfo,
		re:                   re,
		providersToProcessID: providersToProcessID,
		processIDToProcess:   processIDToProcess,
		processID:            0}
}

func (m *Monitor) startMonitor(startedWg *sync.WaitGroup) {
	m.re.log(LOGINFO, "Monitor alive, waiting to receive...")

	// Notify parent that the monitor has started
	startedWg.Done()

	m.monitorLoop()
}

func (m *Monitor) stopMonitor() {
	m.stopMonitorChan <- true
}

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
			m.updateSubscriberRules()
		} else if processUpdate.isRuleDoneBeforeRename {
			// Process finished rule but do not update the process list
			m.re.logMonitorf("Finished %s, %s\n", RuleString[processUpdate.rule], processUpdate.process.String())
			m.rulesLog = append(m.rulesLog, MonitorRulesLog{Process: processUpdate.process, Rule: processUpdate.rule})
			m.updateSubscriberRules()
		} else if processUpdate.updatedProcess {
			// Process updated form
			m.re.logMonitorf("Updated %s\n", processUpdate.process.String())
			// m.rulesLog = append(m.rulesLog, MonitorRulesLog{Process: processUpdate.process, Rule: processUpdate.rule})
			m.updateProcessToList(&processUpdate.process)
		} else if processUpdate.newProcess {
			// Process finished rule
			m.re.logMonitorf("New process %s\n", processUpdate.process.String())
			// m.rulesLog = append(m.rulesLog, MonitorRulesLog{Process: processUpdate.process, Rule: processUpdate.rule})

			m.addProcessToList(&processUpdate.process)
		}

		// Send updated structure to subscriber
		m.updateSubscriberProcesses()

	case <-m.stopMonitorChan:
		m.re.logMonitorf("Monitor terminating\n")
		return

	case error := <-m.errorChan:
		fmt.Println(error)
	}

	m.monitorLoop()
}

func (m *Monitor) updateSubscriberProcesses() {
	// Send list of processes and links to the subscriber
	if m.subscriber != nil {

		v := make([]ProcessInfo, 0, len(m.processIDToProcess))

		// Prepare list of processes
		for id, value := range m.processIDToProcess {
			providers := make([]string, 0, len(value.Providers))

			for _, provider := range value.Providers {
				providers = append(providers, provider.String())
			}

			v = append(v, ProcessInfo{ID: strconv.Itoa(id), Providers: providers, Body: value.Body.String(), Polarity: types.PolarityMap[value.Body.Polarity(m.re.Typechecked, m.re.GlobalEnvironment)]})
		}

		// Prepare list of links between processes
		links := []Link{}
		for pID_source, process := range m.processIDToProcess {
			for _, freeName := range process.Body.FreeNames() {
				pID_destination, exists := m.providersToProcessID[freeName]
				if exists {
					links = append(links, Link{strconv.Itoa(pID_source), strconv.Itoa(pID_destination)})
				}
			}
		}

		m.subscriber.ProcessesSubscriberChan <- ProcessesStructure{ProcessInfo: v, Links: links}
	}
}

func (m *Monitor) updateSubscriberRules() {
	// Send list of rules to the subscriber
	if m.subscriber != nil {

		v := make([]RuleInfo, 0, len(m.rulesLog))

		for i, value := range m.rulesLog {
			providers := make([]string, 0, len(value.Process.Providers))

			for _, provider := range value.Process.Providers {
				providers = append(providers, provider.String())
			}

			v = append(v, RuleInfo{ID: strconv.Itoa(i), Providers: providers, Rule: RuleString[value.Rule]})
		}
		m.subscriber.RulesSubscriberChan <- v
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
	position := process.Position

	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, nil, shape, position), rule: rule, isRuleDone: true}
}

func (m *Monitor) MonitorRuleFinishedBeforeRenamed(process *Process, rule Rule) {
	// Since the provider will be stolen and used by a different process, the monitor should only note that a rule
	// has finished, without updating the list of current processes
	body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape
	position := process.Position

	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, nil, shape, position), rule: rule, isRuleDoneBeforeRename: true}
}

func (m *Monitor) MonitorNewProcess(process *Process) {
	body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape
	position := process.Position

	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, nil, shape, position), newProcess: true}
}

func (m *Monitor) MonitorProcessRenamed(process *Process) {
	body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape
	position := process.Position

	m.monitorChan <- MonitorUpdate{process: *NewProcess(body, provider, nil, shape, position), updatedProcess: true}
}
func (m *Monitor) MonitorProcessTerminated(process *Process) {
	// body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape
	position := process.Position

	m.monitorChan <- MonitorUpdate{process: *NewProcess(nil, provider, nil, shape, position), isDead: true}
}

func (m *Monitor) MonitorProcessForwarded(process *Process) {
	// body := CopyForm(process.Body)
	provider := process.Providers
	shape := process.Shape
	position := process.Position

	m.monitorChan <- MonitorUpdate{process: *NewProcess(nil, provider, nil, shape, position), isDead: true}
}

// Information to interact with the webserver
type SubscriberInfo struct {
	ProcessesSubscriberChan chan ProcessesStructure
	RulesSubscriberChan     chan []RuleInfo
}

type ProcessesStructure struct {
	ProcessInfo []ProcessInfo `json:"processes"`
	Links       []Link        `json:"links"`
}

type ProcessInfo struct {
	ID        string   `json:"id"`
	Providers []string `json:"providers"`
	Body      string   `json:"body"`
	Polarity  string   `json:"polarity"`
}

type Link struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type RuleInfo struct {
	ID        string   `json:"id"`
	Providers []string `json:"providers"`
	Rule      string   `json:"rule"`
}

func NewSubscriberInfo() *SubscriberInfo {
	processSubscriberChan := make(chan ProcessesStructure, 100)
	ruleSubscriberChan := make(chan []RuleInfo, 100)
	return &SubscriberInfo{ProcessesSubscriberChan: processSubscriberChan, RulesSubscriberChan: ruleSubscriberChan}
}

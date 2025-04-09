package modules

type ProcessControlBlock struct {
	PID       int
	Status    string
	Program   []func()
	Pc        int
	Registers map[string]any
}

type ProcessManager struct {
	CurrentProcessPID int
	ReadyProcesses    []ProcessControlBlock
}

func (pm *ProcessManager) Init() {
	pm.CurrentProcessPID = 0
	pm.ReadyProcesses = []ProcessControlBlock{}
}

func (pm *ProcessManager) CreateProcess()  {}
func (pm *ProcessManager) DestroyProcess() {}

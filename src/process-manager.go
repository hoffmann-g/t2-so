package main

var ProgramList = map[string][]func(){
	"notes":     Program1,
	"navigator": Program2,
}

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

	MemoryManager *MemoryManager
}

func (pm *ProcessManager) Init(mm *MemoryManager) {
	pm.CurrentProcessPID = 0
	pm.ReadyProcesses = []ProcessControlBlock{}

	pm.MemoryManager = mm
}

func (pm *ProcessManager) CreateProcess(process_name string) bool {
	program, exists := ProgramList[process_name]
	if !exists {
		panic("Error: Program not found")
	}

	pcb := ProcessControlBlock{
		PID:       1,
		Status:    "READY",
		Program:   program,
		Pc:        0,
		Registers: make(map[string]any),
	}

	if !pm.MemoryManager.AllocateProcess(pcb) {
		panic("Error: Not enough memory to allocate process, please kill yourself")
	}

	pm.ReadyProcesses = append(pm.ReadyProcesses, pcb)

	return true
}

func (pm *ProcessManager) DestroyProcess() {}

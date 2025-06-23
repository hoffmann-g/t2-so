package main

import (
	"errors"
	"fmt"
	"slices"
)

var ContinuousExecution = false

var ProgramList = map[string][]Instruction{
	"notes": Program1,
	"nav":   Program2,
	"AAA":   Program3,
	"13":    Program4,
	"io":    ProgramIO,
}

// PCB - estrutura do processo
// A tabela de páginas agora é mantida no MemoryManager, como map[int][]PageTableEntry
type ProcessControlBlock struct {
	PID           int
	Status        string
	Program       []Instruction
	ProgramLength int
	Pc            int
	Registers     map[string]any
	QuantumUsed   int
}

type ProcessManager struct {
	CurrentProcessPID int
	ReadyProcesses    []ProcessControlBlock
}

func (pm *ProcessManager) Init(mm *MemoryManager) {
	pm.CurrentProcessPID = -1
	pm.ReadyProcesses = []ProcessControlBlock{}
}

func (pm *ProcessManager) CreateProcess(process_name string) (int, error) {
	program, exists := ProgramList[process_name]
	if !exists {
		return 0, errors.New("program not found")
	}

	status := "ADDED"

	if ContinuousExecution {
		status = "READY"
	}

	pcb := ProcessControlBlock{
		PID:           nextPID,
		Program:       program,
		ProgramLength: len(program),
		Pc:            0,
		Registers:     RegistersBase,
		QuantumUsed:   0,
		Status:        status,
	}

	if !kernel.MMU.AllocateProcess(pcb) {
		return 0, errors.New("not enough memory to allocate process")
	}

	pm.ReadyProcesses = append(pm.ReadyProcesses, pcb)

	nextPID++

	if _, _, exists := pm.GetRunningProcess(); !exists && ContinuousExecution {
		kernel.Scheduler.Schedule()
	}

	return pcb.PID, nil
}

func (pm *ProcessManager) DestroyProcess() {
	pid := cpu.Registers["$v0"].(int)

	for i, process := range pm.ReadyProcesses {
		if process.PID == pid {
			LogDebug(fmt.Sprintf("[PMU] Removendo processo PID %d", pid))
			pm.ReadyProcesses = slices.Delete(pm.ReadyProcesses, i, i+1)
			LogTrace("Calling MMU directly")
			kernel.MMU.DeallocateProcess()
			// Após desalocar, chama o escalonador
			kernel.Scheduler.Schedule()
			// Se não houver mais processos, coloca o PC em idle
			if len(pm.ReadyProcesses) == 0 {
				cpu.Pc = IdleStatePc
				LogDebug("[PMU] Nenhum processo restante. PC em idle.")
			}
			LogDebug(fmt.Sprintf("[PMU] PC após destruição: %d", cpu.Pc))
			return
		}
	}
}

func (pm *ProcessManager) Execute(pid int) error {
	LogDebug(fmt.Sprintf("Executing process PID %d", pid))

	for i, process := range pm.ReadyProcesses {
		if process.PID != pid {
			continue
		}
		if process.Status == "RUNNING" {
			return errors.New("process already running")
		}
		if process.Status == "FINISHED" {
			return errors.New("process not found")
		}

		LogDebug(fmt.Sprintf("Setting process PID %d to READY", pid))
		pm.ReadyProcesses[i].Status = "READY"

		LogDebug("Calling scheduler from Execute")
		kernel.Scheduler.Schedule()
		return nil
	}
	return errors.New("process not found")
}

func (pm *ProcessManager) ExecuteAll() error {
	found := false
	for i, process := range pm.ReadyProcesses {
		if process.Status == "RUNNING" {
			return errors.New("process already running")
		}
		if process.Status == "FINISHED" {
			return errors.New("process not found")
		}

		pm.ReadyProcesses[i].Status = "READY"
		found = true

	}

	if !found {
		return errors.New("no processes found")
	}

	kernel.Scheduler.Schedule()
	return nil
}

func (pm *ProcessManager) GetRunningProcess() (ProcessControlBlock, int, bool) {
	for i, process := range pm.ReadyProcesses {
		if process.Status == "RUNNING" {
			return process, i, true
		}
	}
	return ProcessControlBlock{}, -1, false
}

func (pm *ProcessManager) UnblockProcess(pid int) {
	for i, process := range pm.ReadyProcesses {
		if process.PID == pid && process.Status == "BLOCKED" {
			pm.ReadyProcesses[i].Status = "READY"
			return
		}
	}
}

func (pm *ProcessManager) GetProcessByPID(pid int) ProcessControlBlock {
	for _, process := range pm.ReadyProcesses {
		if process.PID == pid {
			return process
		}
	}
	return ProcessControlBlock{PID: -1}
}

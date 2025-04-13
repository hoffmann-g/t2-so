package main

import (
	"errors"
	"slices"
)

var ProgramList = map[string][]any{
	"notes": Program1,
	"nav":   Program2,
	"AAA":   Program3,
	"13":    Program4,
}

type ProcessControlBlock struct {
	PID           int
	Status        string
	Program       []any
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

	pcb := ProcessControlBlock{
		PID:           nextPID,
		Program:       program,
		ProgramLength: len(program),
		Pc:            0,
		Registers:     RegistersBase,
		QuantumUsed:   0,

		Status: "ADDED", // manual execution
		// Status: "READY", // automatic execution
	}

	if !kernel.MMU.AllocateProcess(pcb) {
		return 0, errors.New("not enough memory to allocate process, please kill yourself")
	}

	pm.ReadyProcesses = append(pm.ReadyProcesses, pcb)

	nextPID++

	return pcb.PID, nil
}

func (pm *ProcessManager) DestroyProcess() {
	pid := cpu.Registers["$v0"].(int)

	for i, process := range pm.ReadyProcesses {
		if process.PID == pid {
			pm.ReadyProcesses = slices.Delete(pm.ReadyProcesses, i, i+1)
			LogTrace("Jumping to MMU")
			cpu.jump(DealocateProcessStart)
			return
		}
	}
}

func (pm *ProcessManager) Execute(pid int) error {
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

		pm.ReadyProcesses[i].Status = "READY"

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

func (pm *ProcessManager) GetRunningProcess() (ProcessControlBlock, bool) {
	for _, process := range pm.ReadyProcesses {
		if process.Status == "RUNNING" {
			return process, true
		}
	}
	return ProcessControlBlock{}, false
}

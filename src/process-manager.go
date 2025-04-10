package main

import (
	"errors"
)

var ProgramList = map[string][]func(){
	"notes": Program1,
	"nav":   Program2,
}

type ProcessControlBlock struct {
	PID           int
	Status        string
	Program       []func()
	ProgramLength int
	Pc            int
	Registers     map[string]any
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
		Status:        "ADDED",
		Program:       program,
		ProgramLength: len(program),
		Pc:            0,
		Registers:     make(map[string]any),
	}

	if !kernel.MMU.AllocateProcess(pcb) {
		return 0, errors.New("not enough memory to allocate process, please kill yourself")
	}

	pm.ReadyProcesses = append(pm.ReadyProcesses, pcb)

	nextPID++

	return pcb.PID, nil
}

func (pm *ProcessManager) DestroyProcess(pid int) error {
	for i, process := range pm.ReadyProcesses {
		if process.PID == pid {
			pm.ReadyProcesses[i].Status = "TERMINATED"

			kernel.MMU.DeallocateProcess(process)

			kernel.Scalonator.Scalonate()
			return nil
		}
	}
	return errors.New("process not found")
}

func (pm *ProcessManager) Execute(pid int) error {
	for i, process := range pm.ReadyProcesses {
		if process.PID != pid {
			continue
		}
		if process.Status != "ADDED" {
			return errors.New("process already running")
		}

		pm.ReadyProcesses[i].Status = "READY"

		kernel.Scalonator.Scalonate()
		return nil

	}
	return errors.New("process not found")
}

func (pm *ProcessManager) GetRunningProcess() (ProcessControlBlock, bool) {
	for _, process := range pm.ReadyProcesses {
		if process.Status == "RUNNING" {
			return process, true
		}
	}
	return ProcessControlBlock{}, false
}

package main

import (
	"errors"
)

var ProgramList = map[string][]func(p *Processor){
	"notes":     Program1,
	"navigator": Program2,
}

type ProcessControlBlock struct {
	PID           int
	Status        string
	Program       []func(p *Processor)
	ProgramLength int
	Pc            int
	Registers     map[string]any
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

func (pm *ProcessManager) CreateProcess(process_name string) (int, error) {
	program, exists := ProgramList[process_name]
	if !exists {
		return 0, errors.New("Program not found")
	}

	pcb := ProcessControlBlock{
		PID:           nextPID,
		Status:        "ADDED",
		Program:       program,
		ProgramLength: len(program),
		Pc:            0,
		Registers:     make(map[string]any),
	}

	if !pm.MemoryManager.AllocateProcess(pcb) {
		return 0, errors.New("Not enough memory to allocate process, please kill yourself")
	}

	pm.ReadyProcesses = append(pm.ReadyProcesses, pcb)

	nextPID++

	return pcb.PID, nil
}

func (pm *ProcessManager) DestroyProcess(pid int) {
	for i, process := range pm.ReadyProcesses {
		if process.PID == pid {
			pm.ReadyProcesses[i].Status = "TERMINATED"

			pm.MemoryManager.DeallocateProcess(process)
			pm.CurrentProcessPID = 0

			// k.Scalonator.Scalonate()
			break
		}
	}
}

func (pm *ProcessManager) Execute(pid int) error {
	for i, process := range pm.ReadyProcesses {
		if process.PID == pid {
			pm.ReadyProcesses[i].Status = "READY"

			k.Scalonator.Scalonate()
			return nil
		}
	}
	return errors.New("Process not found")
}

// func (pm *ProcessManager) Execute(pid int) error {
// 	for i, process := range pm.ReadyProcesses {
// 		if process.PID == pid {
// 			pm.ReadyProcesses[i].Status = "READY"

// 			pm.CurrentProcessPID = pid

// 			return nil
// 		}
// 	}
// 	return errors.New("Error: Process not found")
// }

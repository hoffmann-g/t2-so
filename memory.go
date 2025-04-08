package main

import "fmt"

const memorySize = 16000

var Data [memorySize]any

type Kernel struct {
	handlerRoutinesAddresses []func()
	readyProcesses           []ProcessControlBlock

	memoryManager           *MemoryManager
	processManaginterfaceer *ProcessManager
	scalonator              *Scalonator
}

func (k *Kernel) init() {
	k.handlerRoutinesAddresses[0] = k.timeInterruptionRoutine
	k.handlerRoutinesAddresses[1] = k.ioInterruptionRoutine
}

func (k *Kernel) timeInterruptionRoutine() {
	fmt.Println("Time interruption routine")
}

func (k *Kernel) ioInterruptionRoutine() {
	fmt.Println("I/O interruption routine")
}

type ProcessControlBlock struct {
	pid    int
	status string
	data   []any
}

type MemoryManager struct {
	paginationTable []map[int]any
}

//func (mm *MemoryManager) allocateProcess(p Process) {}
//func (mm *MemoryManager) deallocateProcess(p Process) {}

type ProcessManager struct {
	currentProcessPid int
}

//func (pm *ProcessManager) createProcess() {}
//func (pm *ProcessManager) destroyProcess() {}

type Scalonator struct{}

//func (s *Scalonator) scalonate() {}

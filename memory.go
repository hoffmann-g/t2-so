package main

const memorySize = 16000

type Kernel struct {
	handlerRoutinesAddresses []any
	currentProcessPid        int
	readyProcesses           []ProcessControlBlock
}

type ProcessControlBlock struct {
	pid    int
	status string
}

var Data [memorySize]any

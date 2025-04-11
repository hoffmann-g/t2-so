package main

import "fmt"

type Scalonator struct {
}

func (s *Scalonator) Init() {}

func (s *Scalonator) Scalonate() {
	// LogTrace("Scalonating processes 😈")
	// processor.InterruptionBits[2] = false

	// set running process to ready
	for i, process := range kernel.PMU.ReadyProcesses {
		if process.PID != kernel.PMU.CurrentProcessPID {
			continue
		}

		if process.Status != "RUNNING" {
			continue
		}

		kernel.PMU.ReadyProcesses[i].Status = "READY"
		LogDebug("Process " + fmt.Sprint(kernel.PMU.ReadyProcesses[i].PID) + " set from RUNNING to READY")

		break
	}

	// find next process to run
	for i, process := range kernel.PMU.ReadyProcesses {
		if process.Status != "READY" {
			continue
		}

		kernel.PMU.CurrentProcessPID = process.PID
		processor.Pc = process.Pc
		processor.Registers = process.Registers

		kernel.PMU.ReadyProcesses[i].Status = "RUNNING"

		// LogTrace("Process set to run")
		return
	}

	kernel.PMU.CurrentProcessPID = -1
	processor.InterruptionBits[2] = false
	LogDebug("Entering idle state")
}

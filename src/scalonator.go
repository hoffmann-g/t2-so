package main

import "fmt"

type Scalonator struct {
}

func (s *Scalonator) Init() {}

func (s *Scalonator) Scalonate() {
	LogTrace("Scalonating processes 😈")

	for i := range kernel.PMU.ReadyProcesses {
		if kernel.PMU.ReadyProcesses[i].Status != "RUNNING" {
			continue
		}

		kernel.PMU.ReadyProcesses[i].Status = "READY"
		kernel.PMU.ReadyProcesses[i].Pc = proc.Pc
		kernel.PMU.ReadyProcesses[i].Registers = proc.Registers

		break
	}

	for i, process := range kernel.PMU.ReadyProcesses {
		if process.Status != "READY" {
			continue
		}

		LogDebug(fmt.Sprintf("Process %d is ready", process.PID))

		kernel.PMU.CurrentProcessPID = process.PID
		LogDebug(fmt.Sprintf("Current process PID: %d", kernel.PMU.CurrentProcessPID))

		proc.Pc = process.Pc
		proc.Registers = process.Registers

		kernel.PMU.ReadyProcesses[i].Status = "RUNNING"

		return
	}

	kernel.PMU.CurrentProcessPID = -1
	LogDebug("Entering idle state")
}

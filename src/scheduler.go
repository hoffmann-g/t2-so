package main

type Scheduler struct {
}

func (s *Scheduler) Init() {}

func (s *Scheduler) Schedule() {
	LogDebug("Scheduling processes")
	cpu.InterruptionBits[2] = false

	for i, process := range kernel.PMU.ReadyProcesses {
		if process.PID != kernel.PMU.CurrentProcessPID {
			continue
		}

		if process.Status != "RUNNING" {
			continue
		}

		kernel.PMU.ReadyProcesses[i].Status = "READY"

		break
	}

	for i, process := range kernel.PMU.ReadyProcesses {
		if process.Status != "READY" {
			continue
		}

		kernel.PMU.CurrentProcessPID = process.PID
		cpu.Pc = process.Pc
		cpu.Registers = process.Registers

		kernel.PMU.ReadyProcesses[i].Status = "RUNNING"

		return
	}

	kernel.PMU.CurrentProcessPID = -1
	LogDebug("Entering idle state")
}

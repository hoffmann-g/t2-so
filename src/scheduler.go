package main

type Scheduler struct {
}

func (s *Scheduler) Init() {}

func (s *Scheduler) Schedule() {
	LogDebug("Scheduling processes")
	cpu.Registers["$pam"] = 0

	// save current process state
	for i, process := range kernel.PMU.ReadyProcesses {
		if process.PID != kernel.PMU.CurrentProcessPID {
			continue
		}

		if process.Status != "RUNNING" {
			continue
		}

		kernel.PMU.ReadyProcesses[i].Status = "READY"
		kernel.PMU.ReadyProcesses[i].Pc = cpu.Pc
		kernel.PMU.ReadyProcesses[i].Registers = cpu.Registers

		break
	}

	// look for a new process to run
	// CHANGE CODE HERE:
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

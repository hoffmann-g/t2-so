package main

type Scalonator struct {
}

func (s *Scalonator) Init() {}

func (s *Scalonator) Scalonate() {
	// LogTrace("Scalonating processes 😈")
	// kernel.PMU.CurrentProcessPID = -1
	// processor.InterruptionBits[2] = false]
	// THIS DOESNT HAVE TIME TO EXECUTE

	for i := range kernel.PMU.ReadyProcesses {
		if kernel.PMU.ReadyProcesses[i].PID != kernel.PMU.CurrentProcessPID {
			continue
		}

		kernel.PMU.ReadyProcesses[i].Status = "READY"

		break
	}

	// LogTrace("Searching for a process to run")

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

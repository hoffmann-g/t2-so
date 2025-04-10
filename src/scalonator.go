package main

type Scalonator struct {
}

func (s *Scalonator) Init() {}

func (s *Scalonator) Scalonate() {
	// fmt.Println("Scalonating processes 😈")

	// save state of the current process
	for i := 0; i < len(k.ProcessManager.ReadyProcesses); i++ {
		if k.ProcessManager.ReadyProcesses[i].Status == "TERMINATED" {
			k.ProcessManager.ReadyProcesses = append(k.ProcessManager.ReadyProcesses[:i], k.ProcessManager.ReadyProcesses[i+1:]...)
			i--
		}

		if k.ProcessManager.ReadyProcesses[i].Status != "RUNNING" {
			continue
		}

		k.ProcessManager.ReadyProcesses[i].Status = "READY"
		k.ProcessManager.ReadyProcesses[i].Pc = proc.Pc + 1
		k.ProcessManager.ReadyProcesses[i].Registers = proc.Registers

		break
	}

	for i, process := range k.ProcessManager.ReadyProcesses {
		if process.Status != "READY" {
			continue
		}

		k.ProcessManager.CurrentProcessPID = k.ProcessManager.ReadyProcesses[i].PID

		proc.Pc = k.ProcessManager.ReadyProcesses[i].Pc
		proc.Registers = k.ProcessManager.ReadyProcesses[i].Registers

		k.ProcessManager.ReadyProcesses[i].Status = "RUNNING"

		break
	}
}

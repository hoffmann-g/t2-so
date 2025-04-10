package main

type Scalonator struct {
}

func (s *Scalonator) Init() {}

func (s *Scalonator) Scalonate() {
	// fmt.Println("Scalonating processes 😈")

	// save state of the current process
	for i := range k.ProcessManager.ReadyProcesses {
		if k.ProcessManager.ReadyProcesses[i].Status != "RUNNING" {
			continue
		}

		k.ProcessManager.ReadyProcesses[i].Status = "READY"
		k.ProcessManager.ReadyProcesses[i].Pc = proc.Pc
		k.ProcessManager.ReadyProcesses[i].Registers = proc.Registers

		break
	}

	for i, process := range k.ProcessManager.ReadyProcesses {
		if process.Status != "READY" {
			continue
		}

		println("Process", process.PID, "is ready")

		k.ProcessManager.CurrentProcessPID = process.PID
		println("Current process PID:", k.ProcessManager.CurrentProcessPID)

		proc.Pc = process.Pc
		proc.Registers = process.Registers

		k.ProcessManager.ReadyProcesses[i].Status = "RUNNING"

		return
	}

	k.ProcessManager.CurrentProcessPID = -1
}

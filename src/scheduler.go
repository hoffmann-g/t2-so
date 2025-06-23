package main

import (
	"fmt"
)

type Scheduler struct {
	LastIndex   int
	Quantum     int
	quantumLeft int
}

func (s *Scheduler) Init() {
	s.LastIndex = -1
	s.Quantum = 5
	s.quantumLeft = 3

}

func (s *Scheduler) Schedule() {
	LogDebug("Scheduling processes")
	LogDebug(fmt.Sprintf("Current LastIndex: %d", s.LastIndex))
	LogDebug(fmt.Sprintf("Number of processes: %d", len(kernel.PMU.ReadyProcesses)))

	// Salva o estado do processo atual se houver um rodando
	if kernel.PMU.CurrentProcessPID != -1 {
		for i, process := range kernel.PMU.ReadyProcesses {
			if process.PID != kernel.PMU.CurrentProcessPID {
				continue
			}
			if process.Status != "RUNNING" {
				continue
			}

			kernel.PMU.ReadyProcesses[i].Status = "READY"
			LogDebug(fmt.Sprintf("Process PID %d set to READY", process.PID))
			break
		}
	}

	// Log status de todos os processos
	for i, process := range kernel.PMU.ReadyProcesses {
		LogDebug(fmt.Sprintf("Process %d: PID %d, Status %s", i, process.PID, process.Status))
	}

	// round robin
	n := len(kernel.PMU.ReadyProcesses)
	for offset := 1; offset <= n; offset++ {
		i := (s.LastIndex + offset) % n
		process := kernel.PMU.ReadyProcesses[i]

		LogDebug(fmt.Sprintf("Checking process at index %d: PID %d, Status %s", i, process.PID, process.Status))

		if process.Status != "READY" {
			LogDebug(fmt.Sprintf("Skipping process PID %d - not READY", process.PID))
			continue
		}

		// Configura o próximo processo para rodar
		kernel.PMU.CurrentProcessPID = process.PID
		kernel.PMU.ReadyProcesses[i].Status = "RUNNING"

		// Atualiza o índice para o próximo processo na próxima rodada
		s.LastIndex = i
		s.quantumLeft = s.Quantum // Reseta o quantum

		LogDebug(fmt.Sprintf("Scheduled process PID %d at index %d", process.PID, i))
		return
	}

	// Se não houver processos prontos, entra em estado ocioso
	kernel.PMU.CurrentProcessPID = -1
	LogDebug("Entering idle state - no READY processes found")
}

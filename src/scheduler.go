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

	// cpu.Registers = copyRegisters(RegistersBase)
	// NO NEED TO RESET REGISTERS

	// round robin
	n := len(kernel.PMU.ReadyProcesses)
	for offset := 1; offset <= n; offset++ {
		i := (s.LastIndex + offset) % n
		process := kernel.PMU.ReadyProcesses[i]

		if process.Status != "READY" {
			continue
		}

		// Configura o próximo processo para rodar
		kernel.PMU.CurrentProcessPID = process.PID

		cpu.Pc = process.Pc
		cpu.Registers = copyRegisters(process.Registers)

		kernel.PMU.ReadyProcesses[i].Status = "RUNNING"

		// Atualiza o índice para o próximo processo na próxima rodada
		s.LastIndex = i
		s.quantumLeft = s.Quantum // Reseta o quantum

		LogDebug("[SCHED] Processo escolhido para rodar:")
		for _, p := range kernel.PMU.ReadyProcesses {
			LogDebug(fmt.Sprintf("[SCHED] PID %d status: %s", p.PID, p.Status))
		}
		LogDebug(fmt.Sprintf("[SCHED] CurrentProcessPID: %d", kernel.PMU.CurrentProcessPID))

		return
	}

	// Se não houver processos prontos, entra em estado ocioso
	kernel.PMU.CurrentProcessPID = -1
	LogDebug("Entering idle state")
	for _, p := range kernel.PMU.ReadyProcesses {
		LogDebug(fmt.Sprintf("[SCHED] PID %d status: %s", p.PID, p.Status))
	}
	LogDebug(fmt.Sprintf("[SCHED] CurrentProcessPID: %d", kernel.PMU.CurrentProcessPID))
	LogDebug(fmt.Sprintf("[SCHED] PC atual: %d", cpu.Pc))
	LogDebug(fmt.Sprintf("[SCHED] Lista de processos: %v", kernel.PMU.ReadyProcesses))
}

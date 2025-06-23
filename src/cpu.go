package main

import (
	"time"

	"fmt"

	"github.com/gorilla/websocket"
)

var kernel *Kernel = &Kernel{}
var cpu *CPU = &CPU{}

var clockCycleTime = 1 * time.Second

// var clockCycleTime = 300 * time.Millisecond

var RegistersBase = map[string]any{
	"$zero": 0,
	"$jump": 0,
	"$v0":   0,
	"$v1":   0,
	"$a0":   0,
	"$a1":   0,
	"$a2":   0,
	"$a3":   0,
	"$t0":   0,
	"$t1":   0,
	"$t2":   0,
	"$t3":   0,
}

type CPU struct {
	InterruptionBits map[int]bool
	Pc               int
	Registers        map[string]any

	Conn *websocket.Conn
}

func (c *CPU) Init() {
	c.InterruptionBits = make(map[int]bool)
	c.Registers = make(map[string]any)

	c.InterruptionBits[0] = false
	c.InterruptionBits[1] = false
	c.InterruptionBits[2] = false
	c.InterruptionBits[3] = false

	c.Pc = 0

	c.Registers = copyRegisters(RegistersBase)
}

func (c *CPU) Run() {
	for c.Pc < len(Data) {
		time.Sleep(clockCycleTime)
		SendStatusToWS()

		if kernel.PMU.CurrentProcessPID == -1 {
			continue
		}

		// Checa interrupção de tempo ANTES de executar a instrução
		if kernel.Scheduler.quantumLeft <= 0 && c.Pc != 0 {
			c.InterruptionBits[0] = true
		}

		// Trata interrupções ANTES de executar a instrução
		bit, _, isInterrupt := c.getInterruptions()
		if isInterrupt {
			c.handleInterruption(bit)
			continue
		}

		executed := c.executeInstruction()

		if jumpAddress, ok := c.Registers["$jump"].(int); ok && jumpAddress != -1 {
			c.Pc = jumpAddress
			c.Registers["$jump"] = -1

			LogDebug("Jumping to address: " + fmt.Sprint(jumpAddress))

			continue
		}

		process, i, exists := kernel.PMU.GetRunningProcess()

		if exists {
			if c.Pc == (process.ProgramLength - 1) {
				c.InterruptionBits[2] = true
				continue
			}

			if executed {
				kernel.Scheduler.quantumLeft--
				kernel.PMU.ReadyProcesses[i].QuantumUsed++
				c.Pc++
			}
		}
	}
}

func (c *CPU) executeInstruction() bool {
	LogDebug(fmt.Sprintf("[CPU] Executando PID %d, PC=%d, página=%d", kernel.PMU.CurrentProcessPID, c.Pc, c.Pc/FrameSize))
	physAddr := kernel.MMU.GetPhysicalPcAddress(c.Pc)
	if physAddr == -1 {
		LogDebug("Tratando page-fault: chamando HandlePageFault")
		c.saveCurrentProcessStatus()
		kernel.MMU.HandlePageFault(kernel.PMU.CurrentProcessPID, c.Pc/FrameSize)
		c.callScheduler()
		return false // Só não avança o PC em caso de page fault
	}
	code := Data[physAddr].Code
	if instruction, ok := code.(func()); ok {
		instruction()
		return true // Executou função, avança PC
	}
	// Não era função, mas não foi page fault: avança PC
	return true
}

func (c *CPU) getInterruptions() (bit int, address int, occured bool) {
	if c.InterruptionBits[0] {
		return 0, -1, true
	}
	if c.InterruptionBits[1] {
		return 1, -1, true
	}
	if c.InterruptionBits[2] {
		return 2, -1, true
	}
	if c.InterruptionBits[3] {
		return 3, -1, true
	}
	return 0, 0, false
}

func (c *CPU) saveCurrentProcessStatus() {
	for i := range kernel.PMU.ReadyProcesses {
		if kernel.PMU.ReadyProcesses[i].PID != kernel.PMU.CurrentProcessPID {
			continue
		}
		kernel.PMU.ReadyProcesses[i].Pc = cpu.Pc
		kernel.PMU.ReadyProcesses[i].Registers = copyRegisters(cpu.Registers)
		break
	}
}

func (c *CPU) handleInterruption(bit int) {
	c.saveCurrentProcessStatus()
	c.InterruptionBits[bit] = false

	switch bit {
	case 0: // Time interruption
		LogTrace("Time interruption - calling scheduler directly")
		kernel.Scheduler.quantumLeft = kernel.Scheduler.Quantum
		c.callScheduler()
	case 1: // I/O interruption
		LogTrace("I/O interruption - blocking process and calling scheduler")
		c.blockCurrentProcess()
		c.callScheduler()
	case 2: // Stop interruption
		LogTrace("Stop interruption - calling PMU directly")
		LogDebug("[CPU] Tratando interrupção de término: destruindo processo e indo para idle se necessário.")
		cpu.Registers["$v0"] = kernel.PMU.CurrentProcessPID
		kernel.PMU.DestroyProcess()
	case 3: // I/O Response interruption
		LogTrace("I/O Response interruption - unblocking process and calling scheduler")
		c.unblockIOProcess()
		c.callScheduler()
	}
}

func (c *CPU) callScheduler() {
	LogTrace("Calling scheduler directly")
	kernel.Scheduler.Schedule()

	LogDebug(fmt.Sprintf("After scheduler: CurrentProcessPID = %d", kernel.PMU.CurrentProcessPID))

	// Após o scheduler configurar o próximo processo, pula para o PC do processo
	if kernel.PMU.CurrentProcessPID != -1 {
		process, _, exists := kernel.PMU.GetRunningProcess()
		if exists {
			LogDebug(fmt.Sprintf("Found running process: PID %d, PC %d", process.PID, process.Pc))
			cpu.Pc = process.Pc
			cpu.Registers = copyRegisters(process.Registers)
			LogTrace(fmt.Sprintf("Jumping to next scheduled process PID %d at PC %d", process.PID, process.Pc))
		} else {
			LogDebug("No running process found after scheduler")
		}
	} else {
		LogDebug("CurrentProcessPID is -1 after scheduler")
	}
}

func (c *CPU) blockCurrentProcess() {
	for i, process := range kernel.PMU.ReadyProcesses {
		if process.PID == kernel.PMU.CurrentProcessPID && process.Status == "RUNNING" {
			kernel.PMU.ReadyProcesses[i].Pc++
			kernel.PMU.ReadyProcesses[i].Status = "BLOCKED"
			return
		}
	}
}

func (c *CPU) unblockIOProcess() {
	// Procura por processos que estavam bloqueados por IO e cuja resposta foi dada
	for i, process := range kernel.PMU.ReadyProcesses {
		if process.Status == "BLOCKED" && process.Registers["$t0"] != nil {
			kernel.PMU.ReadyProcesses[i].Status = "READY"
			break
		}
	}
}

func copyRegisters(src map[string]any) map[string]any {
	dst := make(map[string]any)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (c *CPU) RequestIO(message string) {
	process, _, exists := kernel.PMU.GetRunningProcess()
	c.blockCurrentProcess()
	if !exists {
		return
	}
	// Adiciona pedido de IO à lista global
	kernel.AddIORequest(process.PID, message)
	// Sinaliza interrupção de IO
	c.Pc++
	c.InterruptionBits[1] = true
}

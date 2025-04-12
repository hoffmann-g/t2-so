package main

import (
	"time"

	"fmt"

	"github.com/gorilla/websocket"
)

var kernel *Kernel = &Kernel{}
var cpu *CPU = &CPU{}

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

	c.Pc = 0

	c.Registers = RegistersBase
}

func (c *CPU) Run() {
	for c.Pc < len(Data) {
		// simulate interruption time each 10 instructions
		if c.Pc%10 == 0 && c.Pc != 0 && !c.InterruptionBits[2] {
			c.InterruptionBits[0] = true
		}

		time.Sleep(1 * time.Second)
		SendStatusToWS()

		if kernel.PMU.CurrentProcessPID == -1 {
			continue
		}

		c.executeInstruction()

		bit, isrAddress, isInterrupt := c.getInterruptions()
		if isInterrupt {
			c.saveCurrentProcessStatus()
			c.Pc = isrAddress
			c.InterruptionBits[2] = true

			c.InterruptionBits[bit] = false

			continue
		}

		if jumpAddress, ok := c.Registers["$jump"].(int); ok && jumpAddress != -1 {
			c.Pc = jumpAddress
			c.Registers["$jump"] = -1

			LogDebug("Jumping to address: " + fmt.Sprint(jumpAddress))

			continue
		}

		process, _ := kernel.PMU.GetRunningProcess()
		if c.Pc == (process.ProgramLength - 1) {
			c.InterruptionBits[2] = true

			c.Registers["$v0"] = kernel.PMU.CurrentProcessPID
			c.Pc = DestroyProcessStart

			continue
		}

		c.Pc++
	}
}

func (c *CPU) executeInstruction() {
	if instruction, ok := (Data)[kernel.MMU.GetPhysicalPcAddress(c.Pc)].(func()); ok {
		instruction()
	} else {
		LogDebug("Invalid instruction at PC: " + fmt.Sprint(c.Pc))
	}
}

func (c *CPU) getInterruptions() (bit int, address int, occured bool) {
	if c.InterruptionBits[0] {
		return 0, TimeISRStart, true
	}

	if c.InterruptionBits[1] {
		return 1, IOISRStart, true
	}

	return 0, 0, false
}

func (c *CPU) jump(address int) {
	LogTrace("Jump to address: " + fmt.Sprint(address))
	c.Registers["$jump"] = address
}

func (c *CPU) saveCurrentProcessStatus() {
	for i := range kernel.PMU.ReadyProcesses {
		if kernel.PMU.ReadyProcesses[i].PID != kernel.PMU.CurrentProcessPID {
			continue
		}

		kernel.PMU.ReadyProcesses[i].Pc = cpu.Pc
		kernel.PMU.ReadyProcesses[i].Registers = cpu.Registers

		break
	}
}

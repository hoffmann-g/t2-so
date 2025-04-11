package main

import (
	"time"

	"fmt"

	"github.com/gorilla/websocket"
)

var kernel *Kernel = &Kernel{}
var processor *Processor = &Processor{}

type Processor struct {
	InterruptionBits map[int]bool
	Pc               int
	Registers        map[string]any

	Conn *websocket.Conn
}

func (p *Processor) Init() {
	p.InterruptionBits = make(map[int]bool)
	p.Registers = make(map[string]any)

	p.InterruptionBits[0] = false
	p.InterruptionBits[1] = false
	p.InterruptionBits[2] = false

	p.Pc = 0

	p.Registers["$jump"] = -1

	p.Registers["$t0"] = 0
	p.Registers["$t1"] = 0
	p.Registers["$t2"] = 0
	p.Registers["$t3"] = 0

	p.Registers["$v0"] = 0
	p.Registers["$v1"] = 0
	p.Registers["$v2"] = 0
	p.Registers["$v3"] = 0
}

func (p *Processor) Run() {
	for p.Pc < len(Data) {
		// simulate interruption time
		if p.Pc%10 == 0 && p.Pc != 0 && !p.InterruptionBits[2] {
			p.InterruptionBits[0] = true
		}

		time.Sleep(1 * time.Second)
		SendStatusToWS()

		if kernel.PMU.CurrentProcessPID == -1 {
			continue
		}

		p.executeInstruction()

		bit, isrAddress, isInterrupt := p.getInterruptions()
		if isInterrupt {
			p.saveCurrentProcessStatus()
			p.Pc = isrAddress
			p.InterruptionBits[2] = true

			p.InterruptionBits[bit] = false

			continue
		}

		if jumpAddress, ok := p.Registers["$jump"].(int); ok && jumpAddress != -1 {
			p.Pc = jumpAddress
			p.Registers["$jump"] = -1

			LogDebug("Jumping to address: " + fmt.Sprint(jumpAddress))

			continue
		}

		process, _ := kernel.PMU.GetRunningProcess()
		if p.Pc == (process.ProgramLength - 1) {
			p.InterruptionBits[2] = true

			p.Registers["$v0"] = kernel.PMU.CurrentProcessPID
			p.Pc = DestroyProcessStart

			continue
		}

		p.Pc++
	}
}

func (p *Processor) executeInstruction() {
	if instruction, ok := (Data)[kernel.MMU.GetPhysicalPcAddress(p.Pc)].(func()); ok {
		instruction()
	} else {
		LogDebug("Invalid instruction at PC: " + fmt.Sprint(p.Pc))
	}
}

func (p *Processor) getInterruptions() (bit int, address int, occured bool) {
	if p.InterruptionBits[0] {
		// p.saveCurrentProcessStatus()
		// p.setNextPc(TimeISRStart)
		// p.InterruptionBits[0] = false
		// p.InterruptionBits[2] = true

		return 0, TimeISRStart, true
	}

	if p.InterruptionBits[1] {
		// p.saveCurrentProcessStatus()
		// p.setNextPc(IOISRStart)
		// p.InterruptionBits[1] = false
		// p.InterruptionBits[2] = true

		return 1, IOISRStart, true
	}

	return 0, 0, false
}

func (p *Processor) jump(address int) {
	// LogTrace("Jump to address: " + fmt.Sprint(address))
	p.Registers["$jump"] = address
}

func (p *Processor) saveCurrentProcessStatus() {
	for i := range kernel.PMU.ReadyProcesses {
		if kernel.PMU.ReadyProcesses[i].PID != kernel.PMU.CurrentProcessPID {
			continue
		}

		kernel.PMU.ReadyProcesses[i].Pc = processor.Pc
		kernel.PMU.ReadyProcesses[i].Registers = processor.Registers

		break
	}
}

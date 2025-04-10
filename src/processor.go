package main

import (
	"time"

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
	for {
		if p.Pc%10 == 0 && p.Pc != 0 {
			p.InterruptionBits[0] = true
		}

		time.Sleep(1 * time.Second)
		SendStatusToWS()

		if kernel.PMU.CurrentProcessPID == -1 {
			continue
		}

		p.executeInstruction()
		p.handlePc()
	}
}

func (p *Processor) executeInstruction() {
	if instruction, ok := (Data)[kernel.MMU.GetPhysicalPcAddress(p.Pc)].(func()); ok {
		instruction()
	}
}

func (p *Processor) handlePc() {
	if p.InterruptionBits[0] {
		p.saveCurrentProcessStatus()
		p.jumpToAddress(TimeISRStart)
		processor.InterruptionBits[0] = false
		processor.InterruptionBits[2] = true
		return
	}

	if p.InterruptionBits[1] {
		p.saveCurrentProcessStatus()
		p.jumpToAddress(IOISRStart)
		processor.InterruptionBits[1] = false
		processor.InterruptionBits[2] = true
		return
	}

	process, _ := kernel.PMU.GetRunningProcess()

	if p.Pc < (process.ProgramLength - 1) {
		p.Pc++
	} else {
		//p.Registers["$v0"] = kernel.PMU.CurrentProcessPID
		//p.jumpToAddress(DestroyProcessStart)
		kernel.PMU.DestroyProcess(kernel.PMU.CurrentProcessPID)
	}
}

func (p *Processor) jumpToAddress(address int) {
	p.Pc = address
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

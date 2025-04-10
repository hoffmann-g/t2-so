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

	p.Pc = 0

	p.Registers["$t0"] = 0
	p.Registers["$t1"] = 0
	p.Registers["$t2"] = 0
	p.Registers["$t3"] = 0
}

func (p *Processor) Run() {
	for {
		// if p.Pc%10 == 0 {
		// 	p.setInterruptionBit(0, true)
		// }

		time.Sleep(1 * time.Second)
		SendMessage()

		if kernel.PMU.CurrentProcessPID == -1 {
			continue
		}

		bit, isrAddress, occured := p.getInterruptionServiceRoutineAddress()
		if occured {
			p.handleInterruption(bit, isrAddress)
		}

		p.executeInstruction()
		p.handlePc()
	}
}

func (p *Processor) handlePc() {
	process, _ := kernel.PMU.GetRunningProcess()

	if p.Pc < (process.ProgramLength - 1) {
		p.Pc++
	} else {
		kernel.PMU.DestroyProcess(kernel.PMU.CurrentProcessPID)
	}
}

func (p *Processor) handleInterruption(bit int, isrAddress int) {
	p.InterruptionBits[bit] = false

	if isr, ok := (Data)[isrAddress].(func(*Processor)); ok {
		isr(p)
	} else {
		LogDebug("Error: ISR not found")
	}
}

func (p *Processor) executeInstruction() {
	if instruction, ok := (Data)[kernel.MMU.GetPhysicalPcAddress(p.Pc)].(func(*Processor)); ok {
		instruction(p)
	}
}

func (p *Processor) getInterruptionServiceRoutineAddress() (int, int, bool) {
	rot_adresses := map[int]int{
		0: 0x000000,
		1: 0x000001,
	}

	if p.InterruptionBits[0] {
		return 0, (rot_adresses[0]), true
	}

	if p.InterruptionBits[1] {
		return 1, (rot_adresses[1]), true
	}

	return -1, -1, false
}

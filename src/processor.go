package main

import (
	"sync"
	"time"

	"fmt"

	"github.com/gorilla/websocket"
)

var kernel *Kernel = &Kernel{}
var processor *Processor = &Processor{}

type Processor struct {
	Mu sync.Mutex

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
		if p.Pc%10 == 0 && p.Pc != 0 && !p.InterruptionBits[2] {
			p.InterruptionBits[0] = true
		}

		time.Sleep(500 * time.Millisecond)
		SendStatusToWS()

		if kernel.PMU.CurrentProcessPID == -1 {
			continue
		}

		p.executeInstruction() // execute the instruction at the current PC, or jump pc
		p.handlePc() // pc ++ or jump to instruction
	}
}

func (p *Processor) executeInstruction() {
	if instruction, ok := (Data)[kernel.MMU.GetPhysicalPcAddress(p.Pc)].(func()); ok {
		// val := reflect.ValueOf(instruction)
		// if val.Kind() == reflect.Func {
		// 	funcName := runtime.FuncForPC(val.Pointer()).Name()
		// 	LogDebug("Executing instruction: " + funcName + " at PC: " + fmt.Sprint(p.Pc))
		// }
		instruction()
		LogDebug(fmt.Sprint(kernel.PMU.CurrentProcessPID))
	} else {
		LogDebug("Invalid instruction at PC: " + fmt.Sprint(p.Pc))
	}
}

func (p *Processor) handlePc() {
	if p.InterruptionBits[0] {
		p.saveCurrentProcessStatus()
		p.setNextPc(TimeISRStart)
		p.InterruptionBits[0] = false
		p.InterruptionBits[2] = true
		return
	}

	if p.InterruptionBits[1] {
		p.saveCurrentProcessStatus()
		p.setNextPc(IOISRStart)
		p.InterruptionBits[1] = false
		p.InterruptionBits[2] = true
		return
	}

	process, exists := kernel.PMU.GetRunningProcess()
	if exists {
		if p.Pc < (process.ProgramLength - 1) {
			p.Pc++
			return
		} else {
			p.InterruptionBits[2] = true
			p.Registers["$v0"] = kernel.PMU.CurrentProcessPID
			p.setNextPc(DestroyProcessStart)
		}
	}
}

func (p *Processor) setNextPc(address int) {
	// LogTrace("Jump to address: " + fmt.Sprint(address))
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

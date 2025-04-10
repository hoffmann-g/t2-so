package main

import (
	"time"

	"github.com/gorilla/websocket"
)

var k *Kernel

type Processor struct {
	Interruption_bits map[int]bool
	Pc                int
	Registers         map[string]any

	conn *websocket.Conn
}

func (p *Processor) Init() {
	p.Interruption_bits = make(map[int]bool)
	p.Registers = make(map[string]any)

	p.Interruption_bits[0] = false
	p.Interruption_bits[1] = false

	p.Pc = 0

	p.Registers["$t0"] = 0
	p.Registers["$t1"] = 0
	p.Registers["$t2"] = 0
	p.Registers["$t3"] = 0

	k = &Kernel{}
}

func (p *Processor) Run() {
	// time.Sleep(5 * time.Second)
	for {
		// if p.Pc%10 == 0 {
		// 	p.setInterruptionBit(0, true)
		// }

		// println("PC: ", p.Pc, "\n")
		if k.ProcessManager.CurrentProcessPID == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		// println("PC: ", p.Pc, "\n")
		// print("Current PID: ", k.ProcessManager.CurrentProcessPID, "\n")

		SendMessage(Data, p)

		interruptionBit, isrAddress := p.getInterruptionServiceRoutineAddress()
		if interruptionBit != -1 {
			p.setInterruptionBit(interruptionBit, false)

			if isr, ok := (Data)[isrAddress].(func(*Processor)); ok {
				isr(p)
			}
		}

		if instruction, ok := (Data)[k.MemoryManager.GetPhysicalPcAddress(p.Pc)].(func(*Processor)); ok {
			// print("Executing instruction at PC: ", p.Pc, "\n")
			instruction(p)
		}

		time.Sleep(1 * time.Second)
		if p.Pc < k.ProcessManager.ReadyProcesses[k.ProcessManager.CurrentProcessPID-1].ProgramLength {
			p.Pc++
		} else {
			p.Pc = 0
			k.ProcessManager.DestroyProcess(k.ProcessManager.CurrentProcessPID)
		}
	}
}

func (p *Processor) setInterruptionBit(bit int, value bool) {
	p.Interruption_bits[bit] = value
}

func (p *Processor) getInterruptionServiceRoutineAddress() (int, int) {
	rot_adresses := map[int]int{
		0: 0x000000,
		1: 0x000001,
	}

	for key, value := range p.Interruption_bits {
		if value {
			return key, rot_adresses[key]
		}
	}

	return -1, -1
}

func (p *Processor) SetConn(conn *websocket.Conn) {
	p.conn = conn
}

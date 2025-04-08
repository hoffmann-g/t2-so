package main

import (
	"fmt"
	"time"
)

type Processor struct {
	interruption_bits map[int]bool
	pc                int
	registers         map[string]any
}

func (p *Processor) init() {
	p.interruption_bits[0] = false
	p.interruption_bits[1] = false

	p.pc = 0

	p.registers["R0"] = 0
	p.registers["R1"] = 0
	p.registers["R2"] = 0
	p.registers["R3"] = 0
}

func main() {
	kernel := Kernel{}
	kernel.init()

	processor := Processor{}
	processor.init()

	processor.run()
}

func (p *Processor) run() {
	for {
		// Check for interruptions
		for key, value := range p.interruption_bits {
			if value {
				if handler, exists := kernel.handlerRoutinesAddresses[key]; exists {
					handler()
					p.interruption_bits[key] = false // Reset the interruption bit after handling
				}
			}
		}

		// Print the program counter
		fmt.Printf("%06x\n", p.pc)

		// Execute the instruction at the current program counter
		if fn, ok := Data[p.pc].(func()); ok {
			fn()
		} else {
			fmt.Println(Data[p.pc])
		}

		// Simulate a delay and increment the program counter
		time.Sleep(1 * time.Second)
		p.pc++
	}
}

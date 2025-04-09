package so

import (
	"fmt"
	"t1-so/src/processor"
	"t1-so/src/so/modules"
)

type Kernel struct {
	ProcessManager *modules.ProcessManager
	MemoryManager  *modules.MemoryManager
	Scalonator     *modules.Scalonator

	TimeIsr func(*processor.Processor)
	IoIsr   func(*processor.Processor)
}

func (k *Kernel) Init() {
	k.MemoryManager = &modules.MemoryManager{}
	k.ProcessManager = &modules.ProcessManager{}
	k.Scalonator = &modules.Scalonator{}

	k.TimeIsr = k.timeInterruptionRoutine
	k.IoIsr = k.ioInterruptionRoutine
}

func (k *Kernel) timeInterruptionRoutine(processor *processor.Processor) {
	fmt.Println("Time interruption routine")

	k.Scalonator.Scalonate(processor, k.ProcessManager)
}

func (k *Kernel) ioInterruptionRoutine(processor *processor.Processor) {
	fmt.Println("I/O interruption routine")
}

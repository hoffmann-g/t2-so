package main

var nextPID = 1

type Kernel struct {
	ProcessManager *ProcessManager
	MemoryManager  *MemoryManager
	Scalonator     *Scalonator

	TimeIsr func(*Processor)
	IoIsr   func(*Processor)
}

func (k *Kernel) Init() {
	k.MemoryManager = &MemoryManager{}
	k.ProcessManager = &ProcessManager{}
	k.Scalonator = &Scalonator{}

	k.MemoryManager.Init()
	k.ProcessManager.Init(k.MemoryManager)
	k.Scalonator.Init()

	k.TimeIsr = k.timeInterruptionRoutine
	k.IoIsr = k.ioInterruptionRoutine

}

func (k *Kernel) timeInterruptionRoutine(processor *Processor) {
	// fmt.Println("Time interruption routine")

	k.Scalonator.Scalonate()
}

func (k *Kernel) ioInterruptionRoutine(processor *Processor) {
	// fmt.Println("I/O interruption routine")
}

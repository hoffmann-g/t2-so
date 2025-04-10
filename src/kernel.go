package main

var nextPID = 1

type Kernel struct {
	PMU        *ProcessManager
	MMU        *MemoryManager
	Scalonator *Scalonator

	TimeISR func(*Processor)
	IOISR   func(*Processor)
}

func (k *Kernel) Init() {
	k.MMU = &MemoryManager{}
	k.PMU = &ProcessManager{}
	k.Scalonator = &Scalonator{}

	k.MMU.Init()
	k.PMU.Init(k.MMU)
	k.Scalonator.Init()

	k.TimeISR = k.timeInterruptionRoutine
	k.IOISR = k.ioInterruptionRoutine
}

func (k *Kernel) timeInterruptionRoutine(processor *Processor) {
	LogTrace("Time interruption routine")

	k.Scalonator.Scalonate()
}

func (k *Kernel) ioInterruptionRoutine(processor *Processor) {
	LogTrace("I/O interruption routine")
}

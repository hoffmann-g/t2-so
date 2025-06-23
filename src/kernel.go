package main

import "fmt"

var nextPID = 1

const (
	MemorySize = 512
	FrameSize  = 8
)

var (
	FrameTableSize = MemorySize / FrameSize
	IdleStatePc    = 0
)

func LoadKernelIntoMemory() {
	for i := range 0 {
		LogDebug(fmt.Sprintf("Allocating frame %d for kernel...", i))
		kernel.MMU.FrameTable[i] = false
	}
}

type Kernel struct {
	PMU       *ProcessManager
	MMU       *MemoryManager
	Scheduler *Scheduler
}

func (k *Kernel) Init() {
	k.MMU = &MemoryManager{}
	k.PMU = &ProcessManager{}
	k.Scheduler = &Scheduler{}

	k.MMU.Init()
	k.PMU.Init(k.MMU)
	k.Scheduler.Init()
}

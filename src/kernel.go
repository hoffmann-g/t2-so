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

	KernelStart = 0

	MMUStart              = KernelStart + 0
	DealocateProcessStart = MMUStart + 1
	AllocateProcessStart  = DealocateProcessStart + 1

	PMUStart            = MMUStart + 2
	CreateProcessStart  = PMUStart + 1
	DestroyProcessStart = CreateProcessStart + 1

	KernelEnd = PMUStart + 2

	KernelSize   = KernelEnd - KernelStart
	KernelFrames = KernelSize / FrameSize
)

func LoadKernelIntoMemory() {
	for i := range KernelFrames {
		LogDebug(fmt.Sprintf("Allocating frame %d for kernel...", i))
		kernel.MMU.FrameTable[i] = false
	}

	Data[DealocateProcessStart] = kernel.MMU.DeallocateProcess
	Data[AllocateProcessStart] = kernel.MMU.AllocateProcess

	Data[CreateProcessStart] = kernel.PMU.CreateProcess
	Data[DestroyProcessStart] = kernel.PMU.DestroyProcess
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

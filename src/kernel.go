package main

var nextPID = 1

const (
	MemorySize = 512
	FrameSize  = 4
)

var (
	FrameTableSize = MemorySize / FrameSize
	IdleStatePc    = 0

	KernelStart = 4

	ISRStart = KernelStart
	TimeISR  = ISRStart
	IOISR    = TimeISR + 1

	MMUStart              = ISRStart + 4
	DealocateProcessStart = MMUStart
	AllocateProcessStart  = DealocateProcessStart + 1

	PMUStart            = MMUStart + 4
	CreateProcessStart  = PMUStart
	DestroyProcessStart = CreateProcessStart + 1

	ScalonatorStart         = PMUStart + 4
	ScalonateProcessesStart = ScalonatorStart

	KernelEnd = ScalonatorStart + 4

	KernelSize   = KernelEnd - KernelStart
	KernelFrames = KernelSize / FrameSize
)

func LoadKernelIntoMemory() {
	for i := range KernelFrames {
		LogDebug("Allocating frames for kernel")
		kernel.MMU.FrameTable[i] = false
	}

	Data[TimeISR] = kernel.TimeISR
	Data[IOISR] = kernel.IOISR

	Data[DealocateProcessStart] = kernel.MMU.DeallocateProcess
	Data[AllocateProcessStart] = kernel.MMU.AllocateProcess

	Data[CreateProcessStart] = kernel.PMU.CreateProcess
	Data[DestroyProcessStart] = kernel.PMU.DestroyProcess

	Data[ScalonateProcessesStart] = kernel.Scalonator.Scalonate
}

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

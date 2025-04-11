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

	ISRStart     = KernelStart
	TimeISRStart = ISRStart
	IOISRStart   = TimeISRStart + 1

	MMUStart              = ISRStart + 4
	DealocateProcessStart = MMUStart
	AllocateProcessStart  = DealocateProcessStart + 1

	PMUStart            = MMUStart + 4
	CreateProcessStart  = PMUStart
	DestroyProcessStart = CreateProcessStart + 1

	ScalonatorStart         = PMUStart + 3
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

	Data[TimeISRStart] = kernel.TimeISR
	Data[IOISRStart] = kernel.IOISR

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

	TimeISR func()
	IOISR   func()
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

func (k *Kernel) timeInterruptionRoutine() {
	LogTrace("Time interruption routine")

	processor.jump(ScalonateProcessesStart)
	// k.Scalonator.Scalonate()
}

func (k *Kernel) ioInterruptionRoutine() {
	LogTrace("I/O interruption routine")
}

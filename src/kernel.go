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

	ISRStart     = KernelStart + 0
	TimeISRStart = ISRStart + 1
	IOISRStart   = TimeISRStart + 1
	StopISRStart = IOISRStart + 1

	MMUStart              = ISRStart + 3
	DealocateProcessStart = MMUStart + 1
	AllocateProcessStart  = DealocateProcessStart + 1

	PMUStart            = MMUStart + 2
	CreateProcessStart  = PMUStart + 1
	DestroyProcessStart = CreateProcessStart + 1

	SchedulerStart         = PMUStart + 2
	ScheduleProcessesStart = SchedulerStart + 1

	KernelEnd = SchedulerStart + 2

	KernelSize   = KernelEnd - KernelStart
	KernelFrames = KernelSize / FrameSize
)

func LoadKernelIntoMemory() {
	for i := range KernelFrames {
		LogDebug(fmt.Sprintf("Allocating frame %d for kernel...", i))
		kernel.MMU.FrameTable[i] = false
	}

	Data[TimeISRStart] = kernel.timeInterruptionRoutine
	Data[IOISRStart] = kernel.ioInterruptionRoutine
	Data[StopISRStart] = kernel.stopInterruptionRoutine

	Data[DealocateProcessStart] = kernel.MMU.DeallocateProcess
	Data[AllocateProcessStart] = kernel.MMU.AllocateProcess

	Data[CreateProcessStart] = kernel.PMU.CreateProcess
	Data[DestroyProcessStart] = kernel.PMU.DestroyProcess

	Data[ScheduleProcessesStart] = kernel.Scheduler.Schedule
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

func (k *Kernel) timeInterruptionRoutine() {
	LogTrace("Time interruption routine")

	kernel.Scheduler.quantumLeft = kernel.Scheduler.Quantum

	LogTrace("Jumping to Scheduler...")
	cpu.jump(ScheduleProcessesStart)
}

func (k *Kernel) ioInterruptionRoutine() {
	LogTrace("I/O interruption routine")

	LogTrace("Jumping to Scheduler...")
	cpu.jump(ScheduleProcessesStart)
}

func (k *Kernel) stopInterruptionRoutine() {
	LogTrace("Stop interruption routine")

	cpu.Registers["$v0"] = kernel.PMU.CurrentProcessPID

	LogTrace("Jumping to PMU...")
	cpu.jump(DestroyProcessStart)
}

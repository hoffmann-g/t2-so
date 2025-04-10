package main

import "fmt"

var FrameSize = 8

type Frame struct {
	IsFree bool
}

type MemoryManager struct {
	ProcessPageTables map[int][]int
	FrameTable        []Frame
}

func (mm *MemoryManager) Init() {
	mm.ProcessPageTables = make(map[int][]int)
	mm.FrameTable = make([]Frame, MemorySize/FrameSize)

	for i := range MemorySize / FrameSize {
		mm.FrameTable[i] = Frame{
			IsFree: true,
		}
	}

	mm.FrameTable[0].IsFree = false
	mm.FrameTable[1].IsFree = false
	mm.FrameTable[2].IsFree = false
	mm.FrameTable[3].IsFree = false
	mm.FrameTable[4].IsFree = false
}

func (mm *MemoryManager) AllocateProcess(p ProcessControlBlock) bool {
	programSize := len(p.Program)
	pageCount := (programSize + FrameSize - 1) / FrameSize // round up

	pageTable := make([]int, pageCount)
	frameUsed := 0
	programIndex := 0

	for i := 0; i < len(mm.FrameTable) && frameUsed < pageCount; i++ {
		if mm.FrameTable[i].IsFree {
			// mark frame as used
			mm.FrameTable[i].IsFree = false

			// store mapping: page N -> frame i
			pageTable[frameUsed] = i

			// copy program instructions into memory
			start := i * FrameSize
			end := start + FrameSize
			for j := start; j < end && programIndex < programSize; j++ {
				Data[j] = p.Program[programIndex]
				programIndex++
			}

			frameUsed++
		}
	}

	if frameUsed < pageCount {
		return false
	}

	// save page table for the process
	mm.ProcessPageTables[p.PID] = pageTable

	LogDebug(fmt.Sprintf("Process PID: %d", p.PID))
	LogDebug("Page Table:")
	for page, frame := range kernel.MMU.ProcessPageTables[p.PID] {
		LogDebug(fmt.Sprintf("Page: %d, Frame: %d", page, frame))
	}
	LogDebug("Frame Table:")
	for i, frame := range kernel.MMU.FrameTable {
		if !frame.IsFree {
			LogDebug(fmt.Sprintf("Frame: %d Is used", i))
		}
	}

	return true
}

func (mm *MemoryManager) DeallocateProcess(p ProcessControlBlock) {
	pageTable := mm.ProcessPageTables[p.PID]
	for _, frame := range pageTable {
		mm.FrameTable[frame].IsFree = true
	}

	delete(mm.ProcessPageTables, p.PID)

	if _, exists := mm.ProcessPageTables[p.PID]; !exists {
		LogDebug("Page is not present")
	}

	for i, frame := range kernel.MMU.FrameTable {
		if !frame.IsFree {
			LogDebug(fmt.Sprintf("Frame: %d Is used", i))
		}
	}

	// NO NEED FOR ITERATING THROUGH DATA AND SETTING IT TO NIL
}

func (mm *MemoryManager) GetPhysicalPcAddress(pc int) int {
	if kernel.PMU.CurrentProcessPID < 1 {
		LogDebug("No process is running")
		return 10
	}

	currentProcessPid := kernel.PMU.CurrentProcessPID
	pageTable := kernel.MMU.ProcessPageTables[currentProcessPid]

	pageSize := FrameSize
	currentPage := pc / pageSize
	offset := pc % pageSize
	currentFrame := pageTable[currentPage]

	physicalAddress := (currentFrame * FrameSize) + offset

	return physicalAddress
}

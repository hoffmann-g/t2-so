package main

import "fmt"

type MemoryManager struct {
	FrameTable        []bool
	ProcessPageTables map[int][]int
}

func (mm *MemoryManager) Init() {
	mm.ProcessPageTables = make(map[int][]int)
	mm.FrameTable = make([]bool, FrameTableSize)

	for i := range MemorySize / FrameSize {
		mm.FrameTable[i] = true
	}

	mm.FrameTable[0] = false
	mm.FrameTable[1] = false
	mm.FrameTable[2] = false
	mm.FrameTable[3] = false
	mm.FrameTable[4] = false
}

func (mm *MemoryManager) AllocateProcess(p ProcessControlBlock) bool {
	programSize := len(p.Program)
	pageCount := (programSize + FrameSize - 1) / FrameSize // round up

	pageTable := make([]int, pageCount)
	frameUsed := 0
	programIndex := 0

	for i := 0; i < len(mm.FrameTable) && frameUsed < pageCount; i++ {
		if mm.FrameTable[i] {
			// mark frame as used
			mm.FrameTable[i] = false

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
		if !frame {
			LogDebug(fmt.Sprintf("Frame: %d Is used", i))
		}
	}

	return true
}

func (mm *MemoryManager) DeallocateProcess(pid int) {
	pageTable := mm.ProcessPageTables[pid]
	for _, frame := range pageTable {
		mm.FrameTable[frame] = true
	}

	delete(mm.ProcessPageTables, pid)

	if _, exists := mm.ProcessPageTables[pid]; !exists {
		// LogDebug("Page table deleted")
	}

	// NO NEED FOR ITERATING THROUGH DATA AND SETTING IT TO NIL
}

func (mm *MemoryManager) GetPhysicalPcAddress(pc int) int {
	if processor.InterruptionBits[2] {
		return pc
	}

	if kernel.PMU.CurrentProcessPID == -1 {
		return IdleStatePc
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

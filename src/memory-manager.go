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

func (mm *MemoryManager) DeallocateProcess() {
	pid := cpu.Registers["$v0"].(int)

	pageTable := mm.ProcessPageTables[pid]
	for _, frame := range pageTable {
		mm.FrameTable[frame] = true
	}

	delete(mm.ProcessPageTables, pid)

	if _, exists := mm.ProcessPageTables[pid]; !exists {
		LogDebug("Page table deleted")
	}

	// NO NEED FOR ITERATING THROUGH DATA AND SETTING IT TO NIL

	LogTrace("Calling scheduler directly after deallocation")
	kernel.Scheduler.Schedule()

	// Após o scheduler configurar o próximo processo, pula para o PC do processo
	if kernel.PMU.CurrentProcessPID != -1 {
		process, _, exists := kernel.PMU.GetRunningProcess()
		if exists {
			cpu.Pc = process.Pc
			cpu.Registers = copyRegisters(process.Registers)
			LogTrace(fmt.Sprintf("Jumping to next scheduled process PID %d at PC %d after deallocation", process.PID, process.Pc))
		}
	}
}

func (mm *MemoryManager) GetPhysicalPcAddress(pc int) int {
	if kernel.PMU.CurrentProcessPID == -1 {
		return IdleStatePc
	}

	currentProcessPid := kernel.PMU.CurrentProcessPID
	pageTable, exists := kernel.MMU.ProcessPageTables[currentProcessPid]
	if !exists {
		LogDebug(fmt.Sprintf("Page table not found for PID %d", currentProcessPid))
		return IdleStatePc
	}

	pageSize := FrameSize
	currentPage := pc / pageSize
	offset := pc % pageSize

	// Verificar se a página está dentro do range da page table
	if currentPage >= len(pageTable) {
		LogDebug(fmt.Sprintf("Page %d out of range for PID %d (page table size: %d)", currentPage, currentProcessPid, len(pageTable)))
		return IdleStatePc
	}

	currentFrame := pageTable[currentPage]
	physicalAddress := (currentFrame * FrameSize) + offset

	return physicalAddress
}

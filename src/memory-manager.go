package main

var FrameSize = 8

type Frame struct {
	IsFree bool
}

type MemoryManager struct {
	PaginationTable map[int][]int
	FrameTable      []Frame
}

func (mm *MemoryManager) Init() {
	// Initialize the memory manager with a pagination table
	// This is a placeholder; actual implementation would depend on the system architecture
	// mm.PaginationTable = make([]map[int]any, 0)
	mm.PaginationTable = make(map[int][]int)
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
	mm.PaginationTable[p.PID] = pageTable

	// Print the pagination table for debugging purposes
	// println("Pagination Table for PID", p.PID, ":")
	// for page, frame := range pageTable {
	// 	println("Page", page, "-> Frame", frame)
	// }

	return true
}

func (mm *MemoryManager) DeallocateProcess(p ProcessControlBlock) {
	// Deallocate the frames used by the process
	pageTable := mm.PaginationTable[p.PID]
	for _, frame := range pageTable {
		mm.FrameTable[frame].IsFree = true
	}

	delete(mm.PaginationTable, p.PID)

	// NO NEED FOR ITERATING THROUGH DATA AND SETTING IT TO NIL
}

func (mm *MemoryManager) GetPhysicalPcAddress(pc int) int {
	// println("Current PC: ", pc, "\n")

	currentProcessPid := k.ProcessManager.CurrentProcessPID
	// println("Current Process PID: ", currentProcessPid, "\n")

	pageTable := k.MemoryManager.PaginationTable[currentProcessPid]

	pageSize := FrameSize

	currentPage := pc / pageSize
	// println("Current Page: ", currentPage, "\n")

	offset := pc % pageSize
	// println("Offset: ", offset, "\n")

	currentFrame := pageTable[currentPage]
	// println("Current Frame: ", currentFrame, "\n")

	physicalAddress := (currentFrame * FrameSize) + offset
	// println("Physical Address: ", physicalAddress, "\n")

	return physicalAddress
}

package main

import "fmt"

// Estrutura de entrada da tabela de páginas
// Cada entrada representa uma página do processo
// FrameNumber: número do quadro na RAM, -1 se não está em RAM
// InMemory: true se está em RAM
// OnDisk: true se está no disco (swap)
// DiskBlock: índice do bloco no disco, -1 se não está no disco
// Referenced: para política de substituição
// LoadedBefore: true se já foi carregada alguma vez

type PageTableEntry struct {
	FrameNumber  int
	InMemory     bool
	OnDisk       bool
	DiskBlock    int
	Referenced   bool
	LoadedBefore bool
}

type MemoryManager struct {
	FrameTable        []bool
	ProcessPageTables map[int][]PageTableEntry
}

// Gerenciador de Disco para swap de páginas
// Cada bloco representa uma página salva no disco

type DiskManager struct {
	Blocks     [][]Instruction // cada bloco é uma página (slice de instruções/dados)
	FreeBlocks []int           // índices dos blocos livres
	Size       int             // número total de blocos
}

func (mm *MemoryManager) Init() {
	mm.ProcessPageTables = make(map[int][]PageTableEntry)
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

	pageTable := make([]PageTableEntry, pageCount)
	programIndex := 0
	frameAllocated := false

	// Tenta alocar apenas a primeira página em um quadro livre
	for i := 0; i < len(mm.FrameTable); i++ {
		if mm.FrameTable[i] {
			// marca quadro como usado
			mm.FrameTable[i] = false

			// Preenche a primeira entrada da tabela de páginas
			pageTable[0] = PageTableEntry{
				FrameNumber:  i,
				InMemory:     true,
				OnDisk:       false,
				DiskBlock:    -1,
				Referenced:   false,
				LoadedBefore: false,
			}

			// Copia as instruções da primeira página para a memória
			start := i * FrameSize
			end := start + FrameSize
			for j := start; j < end && programIndex < programSize; j++ {
				Data[j] = p.Program[programIndex]
				programIndex++
			}

			frameAllocated = true
			break
		}
	}

	if !frameAllocated {
		return false // não há quadro livre para a primeira página
	}

	// As demais páginas são marcadas como não carregadas
	for pg := 1; pg < pageCount; pg++ {
		pageTable[pg] = PageTableEntry{
			FrameNumber:  -1,
			InMemory:     false,
			OnDisk:       false,
			DiskBlock:    -1,
			Referenced:   false,
			LoadedBefore: false,
		}
	}

	// Salva a tabela de páginas do processo
	mm.ProcessPageTables[p.PID] = pageTable

	LogDebug(fmt.Sprintf("Process PID: %d", p.PID))
	LogDebug("Page Table:")
	for page, entry := range mm.ProcessPageTables[p.PID] {
		LogDebug(fmt.Sprintf("Page: %d, Frame: %d, InMemory: %v", page, entry.FrameNumber, entry.InMemory))
	}
	LogDebug("Frame Table:")
	for i, frame := range mm.FrameTable {
		if !frame {
			LogDebug(fmt.Sprintf("Frame: %d Is used", i))
		}
	}

	return true
}

func (mm *MemoryManager) DeallocateProcess() {
	pid := cpu.Registers["$v0"].(int)

	pageTable := mm.ProcessPageTables[pid]
	for _, entry := range pageTable {
		mm.FrameTable[entry.FrameNumber] = true
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

	entry := pageTable[currentPage]
	if !entry.InMemory {
		LogDebug(fmt.Sprintf("PAGE FAULT: PID %d, page %d não está em memória!", currentProcessPid, currentPage))
		// Aqui sinalizaremos o page-fault (tratamento será implementado depois)
		return -1 // ou IdleStatePc
	}

	physicalAddress := (entry.FrameNumber * FrameSize) + offset

	return physicalAddress
}

// Inicializa o disco com um número fixo de blocos
func (dm *DiskManager) Init(size int) {
	dm.Size = size
	dm.Blocks = make([][]Instruction, size)
	dm.FreeBlocks = make([]int, size)
	for i := 0; i < size; i++ {
		dm.FreeBlocks[i] = i
	}
}

// Aloca um bloco livre e retorna seu índice, ou -1 se não houver
func (dm *DiskManager) AllocateBlock() int {
	if len(dm.FreeBlocks) == 0 {
		return -1
	}
	idx := dm.FreeBlocks[0]
	dm.FreeBlocks = dm.FreeBlocks[1:]
	return idx
}

// Libera um bloco, tornando-o disponível novamente
func (dm *DiskManager) FreeBlock(idx int) {
	if idx >= 0 && idx < dm.Size {
		dm.Blocks[idx] = nil
		dm.FreeBlocks = append(dm.FreeBlocks, idx)
	}
}

// Salva uma página no bloco especificado
func (dm *DiskManager) SavePage(idx int, page []Instruction) {
	if idx >= 0 && idx < dm.Size {
		dm.Blocks[idx] = make([]Instruction, len(page))
		copy(dm.Blocks[idx], page)
	}
}

// Carrega uma página do bloco especificado
func (dm *DiskManager) LoadPage(idx int) []Instruction {
	if idx >= 0 && idx < dm.Size && dm.Blocks[idx] != nil {
		page := make([]Instruction, len(dm.Blocks[idx]))
		copy(page, dm.Blocks[idx])
		return page
	}
	return nil
}

// Trata page-fault para um processo e página específica
func (mm *MemoryManager) HandlePageFault(pid int, pageNumber int) {
	LogDebug(fmt.Sprintf("[PF] Tratando page-fault: PID %d, página %d", pid, pageNumber))
	LogDebug(fmt.Sprintf("[PF] PC atual do processo %d: %d", pid, kernel.PMU.GetProcessByPID(pid).Pc))
	for i, entry := range mm.ProcessPageTables[pid] {
		LogDebug(fmt.Sprintf("[PF] Página %d: InMemory=%v", i, entry.InMemory))
	}

	pageTable := mm.ProcessPageTables[pid]
	entry := &pageTable[pageNumber]

	// 0. Garante que o processo está bloqueado
	for i, process := range kernel.PMU.ReadyProcesses {
		if process.PID == pid {
			kernel.PMU.ReadyProcesses[i].Status = "BLOCKED"
			break
		}
	}

	// 1. Tenta alocar um quadro livre
	frame := -1
	for i, free := range mm.FrameTable {
		if free {
			frame = i
			mm.FrameTable[i] = false
			break
		}
	}

	// 2. Se não houver quadro livre, precisa vitimar uma página (FIFO: primeira encontrada)
	if frame == -1 {
		for victimPid, pt := range mm.ProcessPageTables {
			for victimPage, victimEntry := range pt {
				if victimEntry.InMemory && !(victimPid == pid && victimPage == pageNumber) {
					LogDebug(fmt.Sprintf("[PF] Vítima: PID %d, página %d, quadro %d", victimPid, victimPage, victimEntry.FrameNumber))
					// Salva página vítima no disco
					diskBlock := kernel.Disk.AllocateBlock()
					if diskBlock == -1 {
						LogDebug("[PF] ERRO: Disco cheio! Não é possível swapar página.")
						return
					}
					// Copia conteúdo do quadro para o disco
					start := victimEntry.FrameNumber * FrameSize
					end := start + FrameSize
					pageData := make([]Instruction, FrameSize)
					copy(pageData, Data[start:end])
					kernel.Disk.SavePage(diskBlock, pageData)
					// Atualiza tabela de páginas da vítima
					mm.ProcessPageTables[victimPid][victimPage].InMemory = false
					mm.ProcessPageTables[victimPid][victimPage].OnDisk = true
					mm.ProcessPageTables[victimPid][victimPage].DiskBlock = diskBlock
					mm.ProcessPageTables[victimPid][victimPage].FrameNumber = -1
					// Libera quadro
					mm.FrameTable[victimEntry.FrameNumber] = true
					frame = victimEntry.FrameNumber
					break
				}
			}
			if frame != -1 {
				break
			}
		}
	}

	if frame == -1 {
		LogDebug("[PF] ERRO: Não foi possível alocar quadro para page-fault!")
		return
	}

	// 3. Carrega a página demandada
	start := frame * FrameSize
	end := start + FrameSize
	var pageData []Instruction
	if entry.OnDisk {
		// Carrega do disco
		pageData = kernel.Disk.LoadPage(entry.DiskBlock)
		kernel.Disk.FreeBlock(entry.DiskBlock)
		entry.OnDisk = false
		entry.DiskBlock = -1
		LogDebug("[PF] Página carregada do disco.")
	} else {
		// Carrega do programa original
		proc := kernel.PMU.GetProcessByPID(pid)
		if proc.PID == -1 {
			LogDebug("[PF] ERRO: Processo não encontrado!")
			return
		}
		pageData = make([]Instruction, FrameSize)
		for i := 0; i < FrameSize; i++ {
			progIdx := pageNumber*FrameSize + i
			if progIdx < len(proc.Program) {
				pageData[i] = proc.Program[progIdx]
			} else {
				pageData[i] = Instruction{"", nil}
			}
		}
		LogDebug("[PF] Página carregada do programa original.")
	}
	copy(Data[start:end], pageData)

	// 4. Atualiza tabela de páginas
	entry.InMemory = true
	entry.FrameNumber = frame
	entry.LoadedBefore = true
	entry.Referenced = false

	// 5. Desbloqueia o processo (sempre coloca em READY)
	for i, process := range kernel.PMU.ReadyProcesses {
		if process.PID == pid {
			kernel.PMU.ReadyProcesses[i].Status = "READY"
			break
		}
	}
	LogDebug("[PF] Processo desbloqueado após page-fault.")
	kernel.Scheduler.Schedule()
}

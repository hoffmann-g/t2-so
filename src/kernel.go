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
	Disk      *DiskManager // Gerenciador de disco para swap
}

func (k *Kernel) Init() {
	k.MMU = &MemoryManager{}
	k.PMU = &ProcessManager{}
	k.Scheduler = &Scheduler{}
	k.Disk = &DiskManager{} // Inicializa o disco

	k.MMU.Init()
	k.PMU.Init(k.MMU)
	k.Scheduler.Init()
	k.Disk.Init(64) // Exemplo: 64 blocos de swap
}

// Estrutura para pedido de IO
type IORequest struct {
	PID     int
	Message string
	Input   string
	Done    bool
}

// Lista global de IO pendentes
var IORequests []IORequest

func (k *Kernel) AddIORequest(pid int, message string) {
	IORequests = append(IORequests, IORequest{PID: pid, Message: message, Done: false})
}

// Função para liberar processo após input
func (k *Kernel) CompleteIORequest(pid int, input string) {
	for i, req := range IORequests {
		if req.PID == pid && !req.Done {
			IORequests[i].Input = input
			IORequests[i].Done = true
			// Salva a resposta no registrador $t0 do processo
			for j, process := range k.PMU.ReadyProcesses {
				if process.PID == pid {
					if process.Registers == nil {
						process.Registers = make(map[string]any)
					}
					k.PMU.ReadyProcesses[j].Registers["$t0"] = input
					break
				}
			}
			k.PMU.UnblockProcess(pid)
			// Sinaliza interrupção de resposta de IO
			cpu.InterruptionBits[3] = true
			break
		}
	}
}

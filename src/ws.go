package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/gorilla/websocket"
)

// Estado que será enviado para o front-end
type StateMessage struct {
	InterruptionBits map[int]bool             `json:"InterruptionBits"`
	Pc               int                      `json:"Pc"`
	VirtualPc        int                      `json:"VirtualPc"`
	Registers        map[string]any           `json:"Registers"`
	Data             []any                    `json:"Data"`
	PID              int                      `json:"PID"`
	PageTable        []int                    `json:"PageTable"`
	CurrentFrame     int                      `json:"CurrentFrame"`
	ReadyProcess     []ProcessControlBlockDTO `json:"ReadyProcess"`
}

// PCB - estrutura do processo
type ProcessControlBlockDTO struct {
	PID           int            `json:"PID"`
	Status        string         `json:"Status"`
	Pc            int            `json:"Pc"`
	Registers     map[string]any `json:"Registers"`
	QuantumUsed   int            `json:"QuantumUsed"`
	ProgramLength int            `json:"ProgramLength"`
}

// Função que envia o estado atual para o WebSocket
func SendStatusToWS() {
	if cpu.Conn != nil {
		modifiedData := make([]any, len(Data))
		copy(modifiedData, Data[:])

		for i, v := range modifiedData {
			val := reflect.ValueOf(v)
			if val.Kind() == reflect.Func {
				funcName := runtime.FuncForPC(val.Pointer()).Name()
				funcName = strings.TrimPrefix(funcName, "main.")
				modifiedData[i] = funcName
			}
		}

		pageTable, exists := kernel.MMU.ProcessPageTables[kernel.PMU.CurrentProcessPID]
		if !exists {
			pageTable = []int{}
		}

		currentFrame := 0
		if exists && cpu.Pc/FrameSize < len(pageTable) {
			currentFrame = pageTable[cpu.Pc/FrameSize]
		}

		// Converte os PCBs reais em DTOs
		readyDTOs := make([]ProcessControlBlockDTO, len(kernel.PMU.ReadyProcesses))
		for i, pcb := range kernel.PMU.ReadyProcesses {
			readyDTOs[i] = ProcessControlBlockDTO{
				PID:           pcb.PID,
				Status:        pcb.Status,
				Pc:            pcb.Pc,
				Registers:     pcb.Registers,
				QuantumUsed:   pcb.QuantumUsed,
				ProgramLength: pcb.ProgramLength,
			}
		}

		stateMsg := StateMessage{
			InterruptionBits: cpu.InterruptionBits,
			Pc:               kernel.MMU.GetPhysicalPcAddress(cpu.Pc),
			VirtualPc:        cpu.Pc,
			Registers:        cpu.Registers,
			Data:             modifiedData,
			PID:              kernel.PMU.CurrentProcessPID,
			PageTable:        pageTable,
			CurrentFrame:     currentFrame,
			ReadyProcess:     readyDTOs,
		}

		jsonMsg, marshErr := json.MarshalIndent(stateMsg, "", "  ")
		if marshErr != nil {
			LogDebug(fmt.Sprintf("Failed to marshal message: %v", marshErr))
			return
		}

		if socketErr := cpu.Conn.WriteMessage(websocket.TextMessage, jsonMsg); socketErr != nil {
			LogDebug(fmt.Sprintf("Failed to send message: %v", socketErr))
			return
		}
	}
}

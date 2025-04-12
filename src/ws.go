package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/gorilla/websocket"
)

type StateMessage struct {
	InterruptionBits map[int]bool
	Pc               int
	VirtualPc        int
	Registers        map[string]any
	Data             []any
	PID              int
	PageTable        []int
	CurrentFrame     int
	// send ready processes
}

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

		currentFrame := -1
		if exists && cpu.Pc/FrameSize < len(pageTable) {
			currentFrame = pageTable[cpu.Pc/FrameSize]
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
		}

		jsonMsg, marshErr := json.MarshalIndent(stateMsg, "", "  ")
		if marshErr != nil {
			LogDebug(fmt.Sprintf("Failed to marshal message: %v", marshErr))
		}

		if socketErr := cpu.Conn.WriteMessage(websocket.TextMessage, jsonMsg); socketErr != nil {
			LogDebug(fmt.Sprintf("Failed to send message: %v", socketErr))
			return
		}
	}
}

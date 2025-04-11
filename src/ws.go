package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"

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
}

func SendStatusToWS() {
	if processor.Conn != nil {
		modifiedData := make([]any, len(Data))
		copy(modifiedData, Data[:])

		for i, v := range modifiedData {
			val := reflect.ValueOf(v)
			if val.Kind() == reflect.Func {
				funcName := runtime.FuncForPC(val.Pointer()).Name()
				modifiedData[i] = funcName
			}
		}

		pageTable, exists := kernel.MMU.ProcessPageTables[kernel.PMU.CurrentProcessPID]
		if !exists {
			pageTable = []int{}
		}

		currentFrame := -1
		if exists && processor.Pc/FrameSize < len(pageTable) {
			currentFrame = pageTable[processor.Pc/FrameSize]
		}

		stateMsg := StateMessage{
			InterruptionBits: processor.InterruptionBits,
			Pc:               kernel.MMU.GetPhysicalPcAddress(processor.Pc),
			VirtualPc:        processor.Pc,
			Registers:        processor.Registers,
			Data:             modifiedData,
			PID:              kernel.PMU.CurrentProcessPID,
			PageTable:        pageTable,
			CurrentFrame:     currentFrame,
		}

		jsonMsg, marshErr := json.MarshalIndent(stateMsg, "", "  ")
		if marshErr != nil {
			LogDebug(fmt.Sprintf("Failed to marshal message: %v", marshErr))
		}

		if socketErr := processor.Conn.WriteMessage(websocket.TextMessage, jsonMsg); socketErr != nil {
			LogDebug(fmt.Sprintf("Failed to send message: %v", socketErr))
			return
		}
	}
}

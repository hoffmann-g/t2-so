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
	Registers        map[string]any
	Data             []any
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

		stateMsg := StateMessage{
			InterruptionBits: processor.InterruptionBits,
			Pc:               kernel.MMU.GetPhysicalPcAddress(processor.Pc),
			Registers:        processor.Registers,
			Data:             modifiedData,
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

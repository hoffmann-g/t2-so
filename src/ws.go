package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type StateMessage struct {
	Interruption_bits map[int]bool
	Pc                int
	Registers         map[string]any
	Data              []any
}

func SendMessage() {
	if proc.Conn != nil {
		modifiedData := make([]any, len(Data))
		copy(modifiedData, Data[:])

		for i, v := range modifiedData {
			if fn, ok := v.(func(*Processor)); ok {
				modifiedData[i] = fmt.Sprintf("%T", fn)
			}
		}

		stateMsg := StateMessage{
			Interruption_bits: proc.InterruptionBits,
			Pc:                kernel.MMU.GetPhysicalPcAddress(proc.Pc),
			Registers:         proc.Registers,
			Data:              modifiedData,
		}

		jsonMsg, marshErr := json.MarshalIndent(stateMsg, "", "  ")
		if marshErr != nil {
			LogDebug(fmt.Sprintf("Failed to marshal message: %v", marshErr))
		}

		if socketErr := proc.Conn.WriteMessage(websocket.TextMessage, jsonMsg); socketErr != nil {
			LogDebug(fmt.Sprintf("Failed to send message: %v", socketErr))
			return
		}
	}
}

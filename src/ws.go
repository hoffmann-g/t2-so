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
	if proc.conn != nil {
		// Create a copy of the data slice for modification
		modifiedData := make([]any, len(Data))
		copy(modifiedData, Data[:])

		for i, v := range modifiedData {
			if fn, ok := v.(func(*Processor)); ok {
				modifiedData[i] = fmt.Sprintf("%T", fn)
			}
		}

		stateMsg := StateMessage{
			Interruption_bits: proc.Interruption_bits,
			Pc:                k.MemoryManager.GetPhysicalPcAddress(proc.Pc),
			Registers:         proc.Registers,
			Data:              modifiedData,
		}

		jsonMsg, error := json.MarshalIndent(stateMsg, "", "  ")
		if error != nil {
			fmt.Println("Failed to marshal message:", error)
		}

		if err := proc.conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
			fmt.Println("Failed to send message:", err)
			return
		}
	}
}

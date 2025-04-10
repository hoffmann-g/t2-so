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

func SendMessage(data [MemorySize]any, p *Processor) {
	if p.conn != nil {
		// Create a copy of the data slice for modification
		modifiedData := make([]any, len(data))
		copy(modifiedData, data[:])

		for i, v := range modifiedData {
			if fn, ok := v.(func(*Processor)); ok {
				modifiedData[i] = fmt.Sprintf("%T", fn)
			}
		}

		stateMsg := StateMessage{
			Interruption_bits: p.Interruption_bits,
			Pc:                k.MemoryManager.GetPhysicalPcAddress(p.Pc),
			Registers:         p.Registers,
			Data:              modifiedData,
		}

		jsonMsg, error := json.MarshalIndent(stateMsg, "", "  ")
		if error != nil {
			fmt.Println("Failed to marshal message:", error)
		}

		if err := p.conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
			fmt.Println("Failed to send message:", err)
			return
		}
	}
}

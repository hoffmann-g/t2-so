package main

import (
	"fmt"
	"net/http"
	"t1-so/src/memory"
	"t1-so/src/processor"
	"t1-so/src/so"

	"github.com/gorilla/websocket"
)

var proc *processor.Processor

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Falied to upgrade:", err)
		return
	}

	fmt.Println("Client connected.")

	proc.SetConn(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	fmt.Println("Closing client connection...")
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	conn.Close()
}

func runServer() {
	http.HandleFunc("/ws", handleWS)
	http.Handle("/", http.FileServer(http.Dir("static")))
	fmt.Println("Running observability server on http://localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Initializing kernel...")
	kernel := &so.Kernel{}
	kernel.Init()

	fmt.Println("Loading kernel into memory...")
	memory.InitMemory(kernel)
	data := memory.Data

	fmt.Println("Initializing processor...")
	proc = &processor.Processor{}
	proc.Init()

	fmt.Println("Running processor in background thread...\n")
	go proc.Run(data[:])

	runServer()
}

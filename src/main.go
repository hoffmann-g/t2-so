package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var proc *Processor

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// fmt.Println("Falied to upgrade:", err)
		return
	}

	// fmt.Println("Client connected.")

	proc.SetConn(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	// fmt.Println("Closing client connection...")
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	conn.Close()
}

func runServer() {
	http.HandleFunc("/ws", handleWS)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":8080", nil)
}

func main() {
	data := &Data

	fmt.Println("Initializing kernel...")
	kernel := &Kernel{}
	kernel.Init()

	fmt.Println("Initializing processor...")
	proc = &Processor{}
	proc.Init(kernel)

	fmt.Println("Running processor in background thread...\n")
	go proc.Run(data)

	fmt.Println("Starting observability server on http://localhost:8080\n")
	go runServer()

	HandleShell(kernel)
}

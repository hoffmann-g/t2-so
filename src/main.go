package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:embed static
var staticFiles embed.FS

var LogLevel string = "INFO"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		LogInfo(fmt.Sprintf("Failed to upgrade: %v", err))
		return
	}

	LogTrace("Client connected.")

	cpu.Conn = conn

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	LogTrace("Closing client connection...")
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	conn.Close()
}

func runServer() {
	http.HandleFunc("/ws", handleWS)
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.FS(staticFS)))
	http.ListenAndServe(":8080", nil)
}

func main() {
	LogInfo("Initializing processor...")
	cpu.Init()

	LogInfo("Initializing kernel...")
	kernel.Init()

	LogInfo("Loading kernel into memory...")
	LoadKernelIntoMemory()

	LogInfo("Running processor in background thread...")
	go cpu.Run()

	LogInfo("Starting observability server on http://localhost:8080")
	go runServer()

	HandleShell()
}

func LogTrace(message string) {
	if LogLevel == "TRACE" {
		fmt.Println("[TRACE]:", message)
	}
}

func LogDebug(message string) {
	if LogLevel == "DEBUG" || LogLevel == "TRACE" {
		fmt.Println("[DEBUG]:", message)
	}
}

func LogInfo(message string) {
	if LogLevel == "INFO" || LogLevel == "DEBUG" || LogLevel == "TRACE" {
		fmt.Println("[INFO]:", message)
	}
}

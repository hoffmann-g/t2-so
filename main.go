package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao fazer upgrade:", err)
		return
	}
	defer conn.Close()

	a := 1

	for {
		var list []string

		for i := 0; i < 10; i++ {
			list = append(list, fmt.Sprintf("Item %d %d", a, i))
		}
		a++

		time.Sleep(1 * time.Second)

		msg, _ := json.Marshal(list)
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func runServer() {
	http.HandleFunc("/ws", handleWS)
	http.Handle("/", http.FileServer(http.Dir("static")))
	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

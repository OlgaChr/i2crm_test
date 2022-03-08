package main

import (
	"fmt"

	"github.com/OlgaChr/i2crm_test/framesocket"
	"github.com/gorilla/websocket"
)

func main() {
	WebSocketUrl := "ws://localhost:8080/socket"

	ws, _, err := websocket.DefaultDialer.Dial(WebSocketUrl, nil)
	if err != nil {
		panic(err)
	}

	frameSocket := framesocket.NewFrameSocket(ws)
	defer frameSocket.Close()

	frameSocket.Write([]byte("Hello World!"))
	frame, _ := frameSocket.Read()
	fmt.Println("read in main ", frame)
}

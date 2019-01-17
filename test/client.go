package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
)

var origin = "http://127.0.0.1:8080/"

func main() {
	var url = "ws://127.0.0.1:8080/echo"

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	message := []byte("hello, world!你好")
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", message)

	var msg = make([]byte, 512)
	m, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg[:m])

	ws.Close() //关闭连接
}

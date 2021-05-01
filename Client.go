package main

import (
	"log"
	"time"
	"fmt"
	"os"
	"bufio"
	
	"github.com/gorilla/websocket"
)

func main() {
	var nickname string
	fmt.Print("Insert your nickname to join chat: ")
	fmt.Scan(&nickname)
	
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/socket", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	
	err = c.WriteMessage(websocket.TextMessage, []byte(nickname))
	if err != nil {
		log.Println("read:", err)
		return
	}
	
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()//read '/n' that remain from nickname scan
	go Getter(c)
	for {
		msg, _, err := reader.ReadLine()
		if err != nil {
			log.Println("scan:", err)
		}
		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("read:", err)
			return
		}
	}
}

func Getter(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("%s", message)
		time.Sleep(500*time.Millisecond)
	}
}

package main

import (
	"net/http"
	"io"
	"log"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var UsersConns [](*websocket.Conn)

var Mutex chan int = make(chan int, 1)

func Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	
	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {return }
		
		log.Printf("%s\n", data)
		io.WriteString(w, "successful post")
	} else if req.Method == "OPTIONS" {
		w.WriteHeader(204)
	} else {
		w.WriteHeader(405)
	}
	
}

func Socket(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	UsersConns = append(UsersConns,conn)
	
	messageType, nickname, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	
	SendMessageToOthers(conn, messageType, append(nickname, []byte(" joined")...))
	
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			DeleteConn(conn)
			conn.Close()
			SendMessageToOthers(nil, websocket.TextMessage, append(nickname, []byte(" left")...))
			return
		}
		
		p = append([]byte(": "), p...)
		
		go SendMessageToOthers(conn, messageType, append(nickname, p...))
	}
}

func DeleteConn(ToDel *websocket.Conn) {
	<- Mutex
	for i, conn := range UsersConns {
		if conn == ToDel {
			UsersConns = append(UsersConns[:i], UsersConns[i+1:]...)
			break
		}
	}
	Mutex <- 1
}

func SendMessageToOthers(except *websocket.Conn, messageType int, p []byte) {
	for _, conn := range UsersConns {
		if conn == except {
			continue
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	Mutex <- 1
	http.HandleFunc("/", Handler)
	http.HandleFunc("/socket", Socket)
	
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}



package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

var broadcast = make(chan string)

var clients = make(map[*Client]bool)

func main(){
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.Handle("/ws", websocket.Handler(handleConnections))

	go handleMessages()

	fmt.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(ws *websocket.Conn){
	client := &Client{Conn:ws}
	clients[client] = true

	defer func(){
		ws.Close()
		delete(clients, client)
	}()

	for {
		var msg string 
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			log.Printf("websocket message receive error: %v", err)
			break
		}
		broadcast <-msg
	}
}

func handleMessages(){
	for {
		msg := <-broadcast
		for client := range clients{
			err := websocket.Message.Send(client.Conn, msg)
			if err !=nil {
				log.Printf("websocket message send error: %v", err)
				client.Conn.Close()
				delete(clients,client)
			}
		}
	}
}
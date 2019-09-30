package main

import (
	"fmt"
	"log"
	"net/http"
	"io"
	"bytes"

	"github.com/gorilla/websocket"
)

//TODO: For now store users and chats in global variables but 
//		should move them to a db later. 
var users = make([]string, 0)
var chats = make(map[string][]*websocket.Conn)


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func getChats(w http.ResponseWriter, r *http.Request) {
	if r.Method !=  "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(w, "Bad request.")
		return 
	} 

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	b := new(bytes.Buffer)
	for key, _ := range chats {
		fmt.Fprintf(b, "%s\n", key)
	}
	
	io.WriteString(w, b.String())
	
}

func connect(w http.ResponseWriter, r *http.Request) {
	//Firt get the username from the request
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
	}
	var user string
	if val, ok := r.Form["user"]; ok {
		user = val[0]
		addUser(user)
	} else {
		log.Println("Error getting username from request")
		return
	}

	//Now that we have the user we can setup the websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Server has a new connection to user named:", user)

	handleChat(conn)

}

func addUser(user string) {
	users = append(users, user)
	log.Println("Current users:", users)
}

func addChat(chat_name string, conn *websocket.Conn) {
	if _, exists := chats[chat_name]; exists {
		log.Println("Chat name alreay exists...")
		return
	}

	conns := make([]*websocket.Conn, 1)
	conns[0] = conn
	chats[chat_name] = conns
}

func handleChat(conn *websocket.Conn) {

}

func main() {
	fmt.Println("Starting LiveChatApp...")

	http.HandleFunc("/connect", connect)
	http.HandleFunc("/chats", getChats)
	http.ListenAndServe(":8080", nil)
}

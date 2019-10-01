package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//TODO: For now store users and chats in global variables but
//		should move them to a db later.
var users = make(map[string]bool)
var chats = make(map[string][]*websocket.Conn)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//Handler function for the route: /chats
func getChats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
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

//Handler function for setting up the websocket connection for
//the route: "/connect"
func connect(w http.ResponseWriter, r *http.Request) {
	//Firt get the username from the request
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
	}
	var user string
	if val, ok := r.Form["user"]; ok {
		user = val[0]
		if added := addUser(user); !added {
			//username already exists so exit method
			return
		}
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

	handleChat(user, conn)

}

// Returns true if the user was successfully added to the user list
// and returns false otherwise.
func addUser(user string) bool {
	if _, exists := users[user]; exists {
		//that user name is already being used
		log.Println("Failed to add user: Username already exists.")
		return false
	}

	//otherwise username can be added to the map of users
	users[user] = true
	log.Println("Current users:", users)
	return true
}

//Attempts to create a new chat group
func addChatGroup(chat_name string, conn *websocket.Conn) {
	if _, exists := chats[chat_name]; exists {
		log.Println("Chat name alreay exists...")
		return
	}

	conns := make([]*websocket.Conn, 1)
	conns[0] = conn
	chats[chat_name] = conns
}

//Handles incoming and outgoing messages for a particular user.
func handleChat(user string, conn *websocket.Conn) {
	read_chan := make(chan string)
	go readMessage(user, read_chan, conn)

	for {
		select {
		case incomingMsg, ok := <-read_chan:
			if !ok {
				return
			}
			fmt.Println(incomingMsg)
		}
	}
}

//Waits for incoming messages from a user and then forwards
//them through a channel to the chat handler method of that
//particular user
func readMessage(user string, read_chan chan string, conn *websocket.Conn) {
	defer close(read_chan)

	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Connection to client is over...")
			log.Println(err)
			return
		} else {
			if msgType == 1 {
				read_chan <- user + ": " + string(p)
			}
		}
	}
}

func main() {
	fmt.Println("Starting LiveChatApp...")

	//setup handlers
	http.HandleFunc("/connect", connect)
	http.HandleFunc("/chats", getChats)

	//add a test chat
	chats["test"] = nil

	//start the server
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

//TODO: For now store users and chats in global variables but
//should move them to a db later.
var users = make(map[string]bool)

type connSlice []*websocket.Conn
var chats = make(map[string]*connSlice)

//used to upgrade the http server connection to the Websocket protocol 
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//Handler function for the route: "/chats"
//Only works for GET method
func getChats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest) 
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(w, "Bad request.")
		return
	}

	//otherwise send the keys from the global chats map
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
		//that username is already being used
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

	/*conns := make([]*websocket.Conn, 1)
	conns[0] = conn
	chats[chat_name] = conns*/
}

//Handles incoming and outgoing messages for a particular user.
func handleChat(user string, conn *websocket.Conn) {
	//TODO: Delete this code that adds each conn to the 
	//test chat later. It is just for testing purposes at 
	//the moment 
	var conns connSlice
	//var cs []*websocket.Conn
	conns = *chats["test"]
	//cs = conns
	fmt.Println("BEFORE", conns)
	conns = append(conns, conn)
	chats["test"] = &conns
	fmt.Println("AFTER", conns)

	read_chan := make(chan string)
	go readMessage(user, read_chan, conn)

	for {
		select {
		case incomingMsg, ok := <-read_chan:
			if !ok {
				return
			}

			broadcastMessage(user, incomingMsg)
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
				read_chan <- string(p)
			}
		}
	}
}

//Send the message to all connections that are in the slice 
//mapped to by the chat which is parsed from the msg itself 
func broadcastMessage(user string, msg string) {
	//parse which chat 
	var chat string
	if idx := strings.IndexByte(msg, ':'); idx >= 0 {
		chat = msg[:idx]
		fmt.Println("chat: ", chat)
	} else {
		log.Println("Error: Could not get chat name from message")
		return 
	}

	//iterate through all the connections and send all the messages 
	conns := *chats[chat]
	fmt.Println("I am here")
	fmt.Println(conns)
	for _, conn := range conns {
		if err := conn.WriteMessage(1, []byte(user + ": " + msg)); err != nil {
			fmt.Println(err)
			return
		}
	}

}

func main() {
	fmt.Println("Starting LiveChatApp...")

	//setup handlers
	http.HandleFunc("/connect", connect)
	http.HandleFunc("/chats", getChats)

	//add a test chat
	var conns connSlice
	cs := make([]*websocket.Conn, 0)
	conns = cs
	chats["test"] = &conns

	//start the server
	http.ListenAndServe(":8080", nil)
}

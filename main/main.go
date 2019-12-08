package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/seanmacdonald/LiveChatApp/data"
	"github.com/seanmacdonald/LiveChatApp/handler"
)

//TODO: For now store users and chats in global variable but
//should move them to a db later.
var chat_info data.ChatData

//used to upgrade the http server connection to the Websocket protocol
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Chats struct {
	Chats []string
}

//Handler function for the route: "/chats"
//Only works for GET method
func getChats(w http.ResponseWriter, r *http.Request) {
	chats := chat_info.Chats

	enableCors(&w)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(w, "Bad request.")
		return
	}

	//otherwise send the keys from the chats map
	chats_keys := make([]string, 0, len(chats))
	for k := range chats {
		chats_keys = append(chats_keys, k)
	}

	all_chats := Chats{chats_keys}

	js, err := json.Marshal(all_chats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//Function to enables CORS for http get requests to the chats route
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
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
		if valid := validUserName(user); !valid {
			log.Println("Username cannont have colons", user)
			return
		}

		if added := data.AddUser(user, &chat_info); !added {
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

func validUserName(user string) bool {
	if hasCol := strings.Contains(user, ":"); hasCol {
		return false
	}

	return true
}

//Handles incoming and outgoing messages for a particular user.
func handleChat(user string, conn *websocket.Conn) {
	//TODO: Delete this code that adds each conn to the
	//test chat later. It is just for testing purposes at
	//the moment
	conns := chat_info.Chats["test"]
	//fmt.Println("BEF:", conns)
	conns = append(conns, conn)
	chat_info.Chats["test"] = conns
	//fmt.Println("AFT:", conns)

	read_chan := make(chan string)
	go readMessage(user, read_chan, conn)

	for {
		select {
		case incomingMsg, ok := <-read_chan:
			if !ok {
				return
			}

			msgh.HandleMessage(user, incomingMsg, conn, &chat_info)
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
			log.Println("Connection to", user, "is over...")
			log.Println(err)

			//remove the username so it can be used by someone else later on
			data.RemoveUser(user, conn, &chat_info)
			return
		} else {
			if msgType == 1 {
				read_chan <- string(p)
			}
		}
	}
}

func main() {
	fmt.Println("Starting LiveChatApp...")

	//setup handlers
	http.HandleFunc("/connect", connect)
	http.HandleFunc("/chats", getChats)

	//setup chat info
	chat_info.Users = make(map[string][]string)
	chat_info.Chats = make(map[string][]*websocket.Conn)

	//add a test chat
	//TODO: delete this once a chat can be added another way
	cs := make([]*websocket.Conn, 0)
	chat_info.Chats["test"] = cs

	//start the server
	http.ListenAndServe(":8080", nil)
}

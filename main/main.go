package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/seanmacdonald/LiveChatApp/data"
	"github.com/gorilla/websocket"
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

//Handler function for the route: "/chats"
//Only works for GET method
func getChats(w http.ResponseWriter, r *http.Request) {
	chats := chat_info.Chats

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

//Attempts to create a new chat group
func addChatGroup(chat_name string, conn *websocket.Conn) {
	if _, exists := chat_info.Chats[chat_name]; exists {
		log.Println("Chat name alreay exists...")
		return
	}

	conns := make([]*websocket.Conn, 1)
	conns[0] = conn
	chat_info.Chats[chat_name] = conns
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

			broadcastMessage(user, incomingMsg)
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

			//remove the username so it can be used by someone else later on
			removeUser(user, conn)
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
	var parsedMsg string
	if idx := strings.IndexByte(msg, ':'); idx >= 0 {
		chat = strings.TrimSpace(msg[:idx])
		parsedMsg = strings.TrimSpace(msg[(idx+1):])
		log.Println("Broadcasting to", chat + ":", parsedMsg )
	} else {
		log.Println("Error: Could not get chat name from message")
		return 
	}

	//iterate through all the connections and send all the messages 
	conns := chat_info.Chats[chat]
	for _, conn := range conns {
		if err := conn.WriteMessage(1, []byte(user + ": " + parsedMsg)); err != nil {
			fmt.Println(err)
			return
		}
	}

}

func removeUser(user string, conn *websocket.Conn) {
	//first remove user from users map 
	delete(chat_info.Users, user)

	//next remove all the corresponding conn object from all chats 
	for chat_name, _ := range chat_info.Chats {
		conn_slice := chat_info.Chats[chat_name]

		i := getPos(conn_slice, conn)

		if i >= 0 {
			//then remove it from this slice 
			conn_slice[i] = conn_slice[len(conn_slice)-1] 
			conn_slice[len(conn_slice)-1] = nil  
			conn_slice = conn_slice[:len(conn_slice)-1]
			chat_info.Chats[chat_name] = conn_slice
		}
	}	
}

func getPos(s []*websocket.Conn, conn *websocket.Conn) int {
    for i, val := range s {
        if val == conn {
            return i
        }
    }
    return -1
}

func main() {
	fmt.Println("Starting LiveChatApp...")

	//setup handlers
	http.HandleFunc("/connect", connect)
	http.HandleFunc("/chats", getChats)

	//setup chat info 
	chat_info.Users = make(map[string]bool)
	chat_info.Chats = make(map[string][]*websocket.Conn)

	//add a test chat
	//TODO: delete this once a chat can be added another way 
	cs := make([]*websocket.Conn, 0)
	chat_info.Chats["test"] = cs

	//start the server
	http.ListenAndServe(":8080", nil)
}

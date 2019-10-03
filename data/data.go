package data

import (
	"log"

	"github.com/gorilla/websocket"
)

type ChatData struct {
	Users map[string]bool
	Chats map[string][]*websocket.Conn
}

// Returns true if the user was successfully added to the user list
// and returns false otherwise.
func AddUser(user string, chat_info *ChatData) bool {
	if _, exists := chat_info.Users[user]; exists {
		//that username is already being used
		log.Println("Failed to add user: Username already exists.")
		return false
	}

	//otherwise username can be added to the map of users
	chat_info.Users[user] = true
	//log.Println("Current users:", chat_info.Users)
	return true
}

//Removes the user from the Users map and its corresponding
//websocket connection objects in the Chats map which are both
//part of the ChatData struct
func RemoveUser(user string, conn *websocket.Conn, chat_info *ChatData) {
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

//Helper method for RemoveUser to find the position of the connection
//object in the slice it is contained in
func getPos(s []*websocket.Conn, conn *websocket.Conn) int {
	for i, val := range s {
		if val == conn {
			return i
		}
	}
	return -1
}

//Attempts to create a new chat group. First checks to see if the chat
//name already exits. Returns true if successful and false otherwise.
func AddChatGroup(chat_name string, conn *websocket.Conn, chat_info *ChatData) bool {
	if _, exists := chat_info.Chats[chat_name]; exists {
		log.Println("Chat name alreay exists...")
		return false
	}

	conns := make([]*websocket.Conn, 1)
	conns[0] = conn
	chat_info.Chats[chat_name] = conns
	return true
}

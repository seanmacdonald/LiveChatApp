package msgh

import (
	"log"
	"strings"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/seanmacdonald/LiveChatApp/data"
)

//Figures out what the incoming message is for. There are 4 cases: 
//	1: Broadcast message in form <chat_name>:<user_name>:<msg>
//	2: Make new chat in form :<chat_name>
//	3: Join existing chat in form ::<chat_name>
//	4: Delete existing chat in form :::<chat_name>  
//NOTE: a chat name CANNONT be an empty string. 
func HandleMessage(user string, msg string, conn *websocket.Conn, chat_info *data.ChatData) {
	//first figure out how many semicolons are in a row 
	count := 0 
	for _, char := range msg {
		if char != ':' {
			break
		}
		count++
	}

	//note that count corresponds to the index where chat_name starts in the message 
	switch count {
	case 0: 
		log.Println("Broadcast case")
		broadcastMessage(user, msg, chat_info)
	case 1: 
		log.Println("Make new chat case")
		createChat(count, user, msg, conn, chat_info)
	case 2: 
		log.Println("Join existing chat case")
	case 3: 
		log.Println("Leave existing chat case")
		leaveChat(count, user, msg, conn, chat_info)
	default: 
		log.Println("Error parsing message")
	}
}

//Send the message to all connections that are in the slice
//mapped to by the chat which is parsed from the msg itself
func broadcastMessage(user string, msg string, chat_info *data.ChatData) {
	//parse which chat
	var chat string
	var parsedMsg string
	if idx := strings.IndexByte(msg, ':'); idx >= 0 {
		chat = strings.TrimSpace(msg[:idx])
		parsedMsg = strings.TrimSpace(msg[(idx + 1):])
		log.Println("Broadcasting to", chat+":", parsedMsg)
	} else {
		log.Println("Error: Could not get chat name from message")
		return
	}

	//iterate through all the connections and send all the messages
	conns := chat_info.Chats[chat]
	for _, conn := range conns {
		send_string := chat + ": " + user + ": " + parsedMsg
		if err := conn.WriteMessage(1, []byte(send_string)); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func createChat(pos int, user string, msg string, conn *websocket.Conn, chat_info *data.ChatData) {
	chat_name := msg[pos:]
	data.AddChatGroup(chat_name, conn, chat_info)
}

func leaveChat(pos int, user string, msg string, conn *websocket.Conn, chat_info *data.ChatData) {
	chat_name := msg[pos:]
	data.LeaveChatGroup(chat_name, conn, chat_info)
}
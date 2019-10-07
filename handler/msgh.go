package msgh

import (
	"log"
	"strings"
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/seanmacdonald/LiveChatApp/data"
)

//Figures out what the incoming message is for. There are 4 cases: 
//	1: Broadcast message in form 1<chat_name>:<user_name>:<msg>
//	2: Make new chat in form 2<chat_name>
//	3: Join existing chat in form 3<chat_name>
//	4: Delete existing chat in form 4<chat_name>  
//NOTE: a chat name CANNONT be an empty string nor have semi colons  
func HandleMessage(user string, msg string, conn *websocket.Conn, chat_info *data.ChatData) {
	//get first character of msg which is the case identifier 
	var count int
	if val, err := strconv.ParseInt(string(msg[0]), 10, 0); err == nil {
		count = int(val)
	}

	parsedMsg := msg[1:]
	parsedMsg = strings.TrimSpace(parsedMsg)

	//note that count corresponds to the index where chat_name starts in the message 
	switch count {
	case 1: 
		log.Println("C1 Broadcast message")
		broadcastMessage(user, parsedMsg, chat_info)
	case 2: 
		log.Println("C2 Make new chat:", parsedMsg)
		createChat(user, parsedMsg, conn, chat_info)
	case 3: 
		log.Println("C3 Join existing chat:", parsedMsg)
		joinChat(user, parsedMsg, conn, chat_info)
	case 4: 
		log.Println("C4 Leave existing chat:", parsedMsg)
		leaveChat(user, parsedMsg, conn, chat_info)
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

func createChat(user string, msg string, conn *websocket.Conn, chat_info *data.ChatData) {
	//msg is the chat name
	if hasCol := strings.Contains(msg, ":"); !hasCol {
		data.AddChatGroup(user, msg, conn, chat_info)
	} else {
		log.Println("Invalid chat name because of colon", msg)
	}
}

func joinChat(user string, msg string, conn *websocket.Conn, chat_info *data.ChatData) {
	//msg us the chat name 
	if hasCol := strings.Contains(msg, ":"); !hasCol {
		data.JoinChatGroup(user, msg, conn, chat_info)
	} else {
		log.Println("Invalid chat name because of colon", msg)
	}
}

func leaveChat(user string, msg string, conn *websocket.Conn, chat_info *data.ChatData) {
	//msg is the chat name
	if hasCol := strings.Contains(msg, ":"); !hasCol {
		data.LeaveChatGroup(user, msg, conn, chat_info)
	} else {
		log.Println("Invalid chat name because of colon", msg)
	}
}
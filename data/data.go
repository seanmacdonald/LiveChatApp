package data

import (
	"log"

	"github.com/gorilla/websocket"
)

type ChatData struct {
	Users map[string][]string
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
	chat_info.Users[user] = make([]string, 0, 5) 
	//log.Println("Current users:", chat_info.Users)
	return true
}

//Removes the user from the Users map and its corresponding
//websocket connection objects in the Chats map which are both
//part of the ChatData struct
func RemoveUser(user string, conn *websocket.Conn, chat_info *ChatData) {
	//delete the user's conn from any chat groups they are in 
	deleteUserFromChats(user, conn, chat_info)

	//first remove user from users map
	delete(chat_info.Users, user)

	//next remove all the corresponding conn object from all chats
	for chat_name, _ := range chat_info.Chats {
		conn_slice := chat_info.Chats[chat_name]

		i := getConnPos(conn_slice, conn)

		if i >= 0 {
			//then remove it from this slice
			conn_slice[i] = conn_slice[len(conn_slice)-1]
			conn_slice[len(conn_slice)-1] = nil
			conn_slice = conn_slice[:len(conn_slice)-1]
			chat_info.Chats[chat_name] = conn_slice
		}
	}
}

//Helper method for finding the position of a connection
//object in the slice it is contained in. Returns the position 
//of the element if found and returns -1 if not found.
func getConnPos(s []*websocket.Conn, conn *websocket.Conn) int {
	for i, val := range s {
		if val == conn {
			return i
		}
	}
	return -1
}

//Attempts to create a new chat group. First checks to see if the chat
//name already exits. Returns true if successful and false otherwise.
func AddChatGroup(user_name string, chat_name string, conn *websocket.Conn, chat_info *ChatData) bool {
	if _, exists := chat_info.Chats[chat_name]; exists {
		log.Println("Chat alreay exists:", chat_name)
		return false
	}

	//add conn to the slice of conns in the chat (first make the slice)
	conns := make([]*websocket.Conn, 1)
	conns[0] = conn
	chat_info.Chats[chat_name] = conns

	//add chat to user's chat list 
	addChat(user_name, chat_name, chat_info)
	return true
}

func JoinChatGroup(user string, chat_name string, conn *websocket.Conn, chat_info *ChatData) bool{
	if _, exists := chat_info.Chats[chat_name]; !exists {
		log.Println("Can't join because chat does not exist:", chat_name)
		return false
	}

	//otherwise we add the users conn to the chat 
	conn_slice := chat_info.Chats[chat_name]
	i := getConnPos(conn_slice, conn)

	if i >= 0 {
		log.Println("Conn is already in the chat for:", user)
		return true 
	}

	//add conn to slice of conns
	conn_slice = append(conn_slice, conn)
	chat_info.Chats[chat_name] = conn_slice

	//add chat name to list of user's chats 
	addChat(user, chat_name, chat_info)
	return true
}

//Removes the given conn from the chat group. If the chat group is now 
//empty then the chat is also deleted.
func LeaveChatGroup(user_name string, chat_name string, conn *websocket.Conn, chat_info *ChatData) bool {
	if _, exists := chat_info.Chats[chat_name]; !exists {
		log.Println("Cannot leave chat because it doesn't exist:", chat_name)
		return false
	}

	conn_slice := chat_info.Chats[chat_name]
	i := getConnPos(conn_slice, conn)

	if i >= 0 {
		//then remove it from this slice
		conn_slice[i] = conn_slice[len(conn_slice)-1]
		conn_slice[len(conn_slice)-1] = nil
		conn_slice = conn_slice[:len(conn_slice)-1]
		chat_info.Chats[chat_name] = conn_slice

		removeChat(user_name, chat_name, chat_info)
	} else {
		log.Println("Nothing to remove: Conn not in this slice.")
	}

	//if there is no user left then delete the chat entirely 
	if len(conn_slice) <= 0 {
		deleteChatGroup(chat_name, chat_info)
	}

	return true; 
}

//Deletes the specified chat group from the given Chats map 
func deleteChatGroup(chat_name string, chat_info *ChatData) {
	delete(chat_info.Chats, chat_name)
	log.Println("Deleted chat group:", chat_name)
}

//Deletes the user all user conn objects from each of the chats
//that the user is in. 
func deleteUserFromChats(user string, conn *websocket.Conn, chat_info *ChatData) {
	for _, chat_name := range chat_info.Users[user] {
		LeaveChatGroup(user, chat_name, conn, chat_info)
	} 
}


//Adds the chat string to the user's corresponding slice in ChatData struct 
func addChat(user string, chat string, chat_info *ChatData) bool {
	chats := chat_info.Users[user]
	i := getChatPos(chats, chat)

	if i >= 0 {
		//then the chat already exists for the user 
		log.Println("Did not add chat because it is already there")
		return false 
	}

	//otherwise add the chat 
	chats = append(chats, chat)
	chat_info.Users[user] = chats
	return true 
}

//Removes the chat string from the user's corresponding slice in ChatData structure
func removeChat(user string, chat string, chat_info *ChatData) bool {
	chats := chat_info.Users[user]
	i := getChatPos(chats, chat)
	if i < 0 {
		//then the chat is not in the slice already 
		log.Println("Did not delete chat because it is not in slice")
		return false 
	}

	//otherwise remove the chat
	chats[i] = chats[len(chats)-1]
	chats[len(chats)-1] = ""
	chats = chats[:len(chats)-1]
	chat_info.Users[user] = chats

	return true
}

//Helper method for finding the position of a chat
//string in the slice it is contained in. Returns the position 
//of the element if found and returns -1 if not found.
func getChatPos(s []string, chat string) int {
	for i, val := range s {
		if val == chat {
			return i
		}
	}
	return -1
}
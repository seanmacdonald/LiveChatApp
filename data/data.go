package data 

import (

	"github.com/gorilla/websocket"
)

type ChatData struct {
	Users map[string]bool
	Chats map[string][]*websocket.Conn 

}

// Returns true if the user was successfully added to the user list
// and returns false otherwise.
/*func AddUser(user string) bool {
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

func RemoveUser(user string, conn *websocket.Conn) {
	//first remove user from users map 
	delete(users, user)

	//next remove all the corresponding conn object from all chats 
	for chat_name, _ := range chats {
		conn_slice := chats[chat_name]

		i := getPos(conn_slice, conn)

		if i >= 0 {
			//then remove it from this slice 
			conn_slice[i] = conn_slice[len(conn_slice)-1] 
			conn_slice[len(conn_slice)-1] = nil  
			conn_slice = conn_slice[:len(conn_slice)-1]
			chats[chat_name] = conn_slice
		}
	}	
}*/

func getPos(s []*websocket.Conn, conn *websocket.Conn) int {
    for i, val := range s {
        if val == conn {
            return i
        }
    }
    return -1
}
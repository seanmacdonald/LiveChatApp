package data

import (
	"testing" 

	"github.com/gorilla/websocket"
)

//test addChat function 
func TestAddChat(t *testing.T) {
	var chat_info ChatData
	chat_info.Users = make(map[string][]string)
	chat_info.Chats = make(map[string][]*websocket.Conn)

	//setup some example data 
	chats1 := []string{"first", "second", "third"}
	chat_info.Users["user1"] = chats1

	//add new chat to a list of chats that already exists  
	addChat("user1", "fourth", &chat_info)
	checkWith := []string{"first", "second", "third", "fourth"}
	result := chat_info.Users["user1"]

	if ok := equalSlice(result, checkWith); !ok {
		t.Errorf("Fourth chat was not added. \nGot: %v \nShould be: %v", result, checkWith)	
	} 

	//add new chat to a list that does not exist yet 
	addChat("user2", "first", &chat_info)
	checkWith = []string{"first"}
	result = chat_info.Users["user2"]

	if ok := equalSlice(result, checkWith); !ok {
		t.Errorf("First chat was not added. \nGot: %v \nShould be: %v", result, checkWith)	
	} 


}


//Helper method for testing if two string slices are equals 
func equalSlice(a, b []string) bool {
	if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}
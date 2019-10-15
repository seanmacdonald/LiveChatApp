package data

import (
	"testing"

	"github.com/gorilla/websocket"
)

/*
	Chat related tests *******************************************
*/

//TEST addChat function
func TestAddChat(t *testing.T) {
	var chat_info ChatData
	setUp(&chat_info)

	//add new chat to a list of chats that already exists
	success := addChat("user1", "fourth", &chat_info)
	checkWith := []string{"first", "second", "third", "fourth"}
	result := chat_info.Users["user1"]
	if ok := equalSlice(result, checkWith); !ok {
		t.Errorf("Fourth chat was not added. \nGot: %v \nShould be: %v", result, checkWith)
	} else if success != true {
		t.Errorf("Added chat but method returned fail value.")
	}

	//add new chat to a list that does not exist yet
	success = addChat("user2", "first", &chat_info)
	checkWith = []string{"first"}
	result = chat_info.Users["user2"]
	if ok := equalSlice(result, checkWith); !ok {
		t.Errorf("First chat was not added. \nGot: %v \nShould be: %v", result, checkWith)
	} else if success != true {
		t.Errorf("Added chat but method returned fail value.")
	}

	//try to add a chat that was already added
	success = addChat("user1", "third", &chat_info)
	checkWith = []string{"first", "second", "third", "fourth"}
	result = chat_info.Users["user1"]
	if ok := equalSlice(result, checkWith); !ok {
		t.Errorf("Some chat was added. \nGot: %v \nShould be: %v", result, checkWith)
	} else if success != false {
		t.Errorf("Did not add chat but method returned true value.")
	}

}

//TEST removeChat function
func TestRemoveChat(t *testing.T) {
	var chat_info ChatData
	setUp(&chat_info)

	//remove a middle chat - the "second" chat from the user1 chats
	success := removeChat("user1", "second", &chat_info)
	checkWith := []string{"first", "third"}
	result := chat_info.Users["user1"]
	if ok := equalSlice(result, checkWith); !ok {
		t.Errorf("Second chat was not removed. \nGot: %v \nShould be: %v", result, checkWith)
	} else if success != true {
		t.Errorf("Removed chat but method returned fail value.")
	}

	//remove the last chat in list -  the "third" chat
	success = removeChat("user1", "third", &chat_info)
	checkWith = []string{"first"}
	result = chat_info.Users["user1"]
	if ok := equalSlice(result, checkWith); !ok {
		t.Errorf("Third chat was not removed. \nGot: %v \nShould be: %v", result, checkWith)
	} else if success != true {
		t.Errorf("Removed chat but method returned fail value.")
	}

	//remove chat that doesnt exist in list
	success = removeChat("user1", "DNE", &chat_info)
	checkWith = []string{"first"}
	result = chat_info.Users["user1"]
	if ok := equalSlice(result, checkWith); !ok {
		t.Errorf("Removed something when it should not. \nGot: %v \nShould be: %v", result, checkWith)
	} else if success != false {
		t.Errorf("Didn't remove a chat but method returned true value.")
	}

	//remove the only chat in list -  the "first" chat
	success = removeChat("user1", "first", &chat_info)
	checkWith = []string{}
	result = chat_info.Users["user1"]
	if ok := equalSlice(result, checkWith); !ok {
		t.Errorf("The only chat was not removed. \nGot: %v \nShould be: %v", result, checkWith)
	} else if success != true {
		t.Errorf("Removed chat but method returned fail value.")
	}

}

//Test getChatPos function 
func TestGetChatPos(t *testing.T) {
	chats := []string{"schools", "kitesurfing", "hockey", "another chat"}

	//case 1: element is first in the slice 
	var checkWith int 
	checkWith = 0
	if val := getChatPos(chats, "schools"); val != checkWith {
		t.Errorf("The first chat was not found. \nGot: %v \nShould be: %v", val, checkWith)
	}

	//case 2: element is last in slice 
	checkWith = 3
	if val := getChatPos(chats, "another chat"); val != checkWith {
		t.Errorf("The first chat was not found. \nGot: %v \nShould be: %v", val, checkWith)
	}

	//case 3: element is in middle of slice
	checkWith = 1
	if val := getChatPos(chats, "kitesurfing"); val != checkWith {
		t.Errorf("The first chat was not found. \nGot: %v \nShould be: %v", val, checkWith)
	}

	//case 4: element does not exist in slice 
	checkWith = -1
	if val := getChatPos(chats, "this is not in list"); val != checkWith {
		t.Errorf("The first chat was not found. \nGot: %v \nShould be: %v", val, checkWith)
	}
}



/*
	User related tests *******************************************
*/

//TEST AddUser function 
func TestAddUser(t *testing.T) {
	var chat_info ChatData
	setUp(&chat_info)

	//try to add already existing user 
	success := AddUser("user1", &chat_info)
	if _, ok := chat_info.Users["user1"]; !ok {
		t.Errorf("user1 does not exist in users map for some reason")
	} else if success != false {
		t.Errorf("Method should not return true when the user was already in the map")
	}

	//try adding a new user 
	success = AddUser("newuser", &chat_info)
	if _, ok := chat_info.Users["newuser"]; !ok {
		t.Errorf("newuser was not added to the map")
	} else if success != true {
		t.Errorf("Method should return true if the addition was successful")
	}
}

func TestRemoveUser(t *testing.T) {
	var chat_info ChatData
	setUp(&chat_info)

	//try removing a user that does not exist 
	var conn *websocket.Conn
	conn = nil  
	RemoveUser("userDNE", conn, &chat_info)
	if _, ok := chat_info.Users["userDNE"]; ok {
		t.Errorf("userDNE is in the map")
	}

	//try removing user that does exist in map 
	RemoveUser("user1", conn, &chat_info)
	if _, ok := chat_info.Users["user1"]; ok {
		t.Errorf("user1 is still in the map")
	}
	
}


/*
	Helper methods for test cases ******************************
*/

//Helper method for setting up test data
func setUp(chat_info *ChatData) {
	chat_info.Users = make(map[string][]string)
	chat_info.Chats = make(map[string][]*websocket.Conn)

	//setup some example data
	chats1 := []string{"first", "second", "third"}
	chat_info.Users["user1"] = chats1
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

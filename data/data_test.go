package data

import (
	"testing"

	"github.com/gorilla/websocket"
)

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

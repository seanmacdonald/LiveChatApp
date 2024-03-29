# LiveChatApp
A backend system written in Golang for a live messaging application that handles multiple concurrent users. 

## API Reference 

### List Chats 
Returns text data containing all the current live chats. 
**URL** - /chats  
**Method** - GET  
**URL Params** - none  
**Success Response**  
    Code: 200  
    Content: SomeChat  
             AnotherChat  
**Error Response**  
    Code: 400  
    Content: Bad request.   

### Make Websocket Connection 
Connect a new user to the server using a websocket.  
**URL** - /connect  
**URL Params (required)** - user=[username]  
**Success** - Websocket connection is granted. 
**Fail** - Webscoket connection is not granted most likely because the user name is already in use.  

### Messaging Protocol
The protocol used between the client and server through the websocket connection.


**Client to Server** - Send message to the server with the format:  
There are 4 types of messages the client can send the server.  
- Case 1: Broadcast message in form: 1*chatname*:*username*:*msg* 
- Case 2: Make new chat in form: 2*chatname*  
- Case 3: Join existing chat in form: 3*chatname*  
- Case 4: Leave existing chat in form: 4*chatname*    

Where chatname is the name of the chat that the user wants to send a message to and msg is the message they want to send. Note that the client should ensure that there are not colons used in the chatname nor the username. However, there may be colons in the msg. Also note that once all the users in a chat have left a chat it is then deleted and not longer exists.   

**Server to Client** - The server will send a message to all users part of the chat in the form:   
*chatname*: *username*: *message*   
The client is responsible for parsing the chatname and username for displaying the information however it would like.    

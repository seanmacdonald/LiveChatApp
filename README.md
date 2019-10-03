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

### Messaging Protocol
The protocol used between the client and server through the websocket connection.  
**Client to Server** - Send message to the server with the format:  
*chat_name:message*   
Where chat_name is the name of the chat that the user wants to send a message to and message is the message they want to send. Note that the client should ensure that there are not colons used in the chat_name. However, there may be colons in the message.   
**Server to Client** - The server will send a message to all users part of the chat in the form:   
*chat_name: user_name: message*   
The client is responsible for parsing the chat_name and possibly user_name for displaying the information however it would like.    

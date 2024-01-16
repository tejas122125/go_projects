// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"github.com/gorilla/websocket"
// )

//   type connection struct {
// connections *websocket.Conn
// }

// var connectUser = make(map[*connection]bool)

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }
// func broadcast (p []byte){

// 	for client := range connectUser {
// 		err:= client.connections.WriteMessage(websocket.TextMessage,p)
// 		if err != nil{
			
//          log.Println("errorasasd",err)
// 			return
// 		}
// 	}
// }

// func handleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	upgrader.CheckOrigin = func(r *http.Request) bool {return true}
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
		
// 		log.Println(err)
// 		return
// 	}
// 	// var c = connection{connections:conn}
// 	var client =  &connection{connections: conn}
// 	connectUser[client] =true
// 	defer conn.Close()

// 	for {
// 		messageType, p, err := conn.ReadMessage()
// 		fmt.Println(string(p))
		
// 		if err != nil {

// 			log.Println(err)
// 			return
// 		}
		
// 		broadcast(p)
// 		fmt.Println("printing recieved message");


// 		conn.WriteMessage(messageType,p)
// 		if err := conn.WriteMessage(messageType, p); err != nil {
			
// 			log.Println(err)
// 			return
// 		}
// 	}
// }
// func router (){
// 	http.HandleFunc("/ws",handleWebSocket)
// }
// func main() {
// 	// http.HandleFunc("/ws", handleWebSocket)
// 	router()
// 	fmt.Println("Server is running on :8080")
// 	http.ListenAndServe(":8080", nil)
// }

package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client represents a connected WebSocket client.
type Client struct {
	conn *websocket.Conn
}

// clients stores all connected clients.
var clients = make(map[*Client]bool)

// broadcast sends a message to all connected clients.
func broadcast(message []byte) {
	for client := range clients {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error broadcasting message:", err)
		}
	}
}

// handleWebSocket handles WebSocket connections.
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {return true}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	client := &Client{conn: conn}
	clients[client] = true

	log.Println("Client connected")

	// Listen for messages from the client
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			delete(clients, client)
			break
		}
		// You can handle different message types (TextMessage, BinaryMessage, etc.)
		// Here, we're broadcasting all incoming messages as text messages.
		broadcast(p)
		log.Printf("Received message: %s\n", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Error echoing message:", err)
			delete(clients, client)
			break
		}
	}

	log.Println("Client disconnected")
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	addr := ":8080"
	fmt.Printf("Server is running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

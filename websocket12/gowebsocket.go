package websocket12

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type CurrentClientBroadcast struct {
	message       Message
	currentclient Client
}
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

type Client struct {
	User *websocket.Conn
	Send chan Message
	Room string
	hub *Hub
}

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan CurrentClientBroadcast
	Rooms      map[string]map[*Client]bool
}

func NewPool() *Hub {
	newpool := &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),

		Broadcast: make(chan CurrentClientBroadcast),
		Rooms:     make(map[string]map[*Client]bool),
	}
	return newpool
}

func (Pool *Hub) start() {

	for {
		select {
		case client := <-Pool.Register:
			Pool.Clients[client] = true
			fmt.Println("new usser connected lenght of clients is", len(Pool.Clients))
			_, err := Pool.Rooms[client.Room]
			if err == false {
				Pool.Rooms[client.Room] = make(map[*Client]bool)
				Pool.Rooms[client.Room][client] = true
			}

		case client := <-Pool.Unregister:

			if _, ok := Pool.Clients[client]; ok {
				delete(Pool.Clients, client)
				close(client.Send)

			}
			if _, ok := Pool.Rooms[client.Room]; ok {
				delete(Pool.Rooms[client.Room], client)
				if len(Pool.Rooms[client.Room]) == 0 {
					delete(Pool.Rooms, client.Room)
				}
			}

		case current := <-Pool.Broadcast:

			for client := range Pool.Rooms[current.currentclient.Room] {
				select {

				case client.Send <- current.message:
				default:
					fmt.Println("no client with given room present")
				}
			}

		}
	}
}
func(client *Client) read (){

	defer func ()  {
		client.hub.Unregister <- client
		client.User.Close()
	}()

for{
	newclient := client.User
	messageType,message,err := newclient.ReadMessage()
	if err!= nil {
	msgobj := Message{Type: messageType,Body: string(message)}
	
	broadobj:= CurrentClientBroadcast{message:  msgobj,currentclient: *client}
	client.hub.Broadcast<- broadobj
	fmt.Println("message recieved from frontend")
	}

}
	

}
func (client *Client) write(){
	defer func ()  {
	client.User.Close()	
	}()


	for{
		select
		{

			// to write message from server create a message object and sen it to user connection channel
		case message:= <- client.Send:
			client.User.WriteMessage(message.Type,[]byte(message.Body))
			fmt.Println("written message propertly")

		}
	}
}

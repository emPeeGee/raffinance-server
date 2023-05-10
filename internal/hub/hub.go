package hub

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   uint
	Conn *websocket.Conn
}

func NewClient(id uint, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
	}
}

func (c *Client) Send(message []byte) {
	err := c.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Println("Error writing message to client:", err)
	}
}

func (c *Client) Listen(hub *Hub) {
	defer func() {
		hub.RemoveClient(c.ID)
		c.Conn.Close()
	}()

	for {
		messageType, message, err := c.Conn.ReadMessage()
		fmt.Println("MSG TYPE", messageType)
		if err != nil {
			fmt.Println("Error reading message from client:", err)
			return
		}

		// Send message only to the sender
		c.Send(message)
	}

	// for {
	// 	_, message, err := c.Conn.ReadMessage()
	// 	if err != nil {
	// 		fmt.Println("Error reading message:", err)
	// 		break
	// 	}

	// 	fmt.Printf("Received message from client %s: %s\n", c.ID, message)

	// 	// Broadcast the message to all clients except the sender
	// 	hub.Broadcast(message, c)
	// }
}

type Hub struct {
	clients map[uint]*Client
	mutex   sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint]*Client),
	}
}

func (r *Hub) AddClient(client *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.clients[client.ID] = client
}

func (r *Hub) RemoveClient(userID uint) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.clients, userID)
}

func (r *Hub) GetClient(userID uint) *Client {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.clients[userID]
}

func (r *Hub) Broadcast(message []byte, sender *Client) {
	for userID, client := range r.clients {
		if userID != sender.ID {
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error writing message to client:", err)
			}
		}
	}
}

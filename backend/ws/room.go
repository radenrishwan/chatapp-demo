package ws

import (
	"context"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Room struct {
	Join      chan *Client
	Left      chan *Client
	Broadcast chan Message
	Clients   map[*Client]string
	Pool      *Pool
}

func NewRoom(pool *Pool) *Room {
	return &Room{
		Join:      make(chan *Client),
		Left:      make(chan *Client),
		Broadcast: make(chan Message),
		Clients:   map[*Client]string{},
		Pool:      pool,
	}
}

func (r Room) Run() {
	for {
		select {
		case client := <-r.Join:
			log.Printf("User with id : %s joined to room chat on room %s", client.Id, client.RoomId)

			// add user to connection pool
			r.Clients[client] = client.Id

			// broadcast to all user when new user joined to chat
			for b := range r.Clients {
				err := wsjson.Write(context.Background(), b.Conn, Message{
					For:         JOIN,
					UserId:      client.Id,
					MessageType: websocket.MessageText,
					Body:        "An user joined with id : " + client.Id,
				})

				if err != nil {
					log.Println("Error on user joined : ", err.Error())
				}
			}

			break
		case client := <-r.Left:
			delete(r.Clients, client) // remove user from connection pool
			r.Pool.DeleteEmptyRoom(client.RoomId)

			// broadcast to all user when user has left from chat
			for b := range r.Clients {
				err := wsjson.Write(context.Background(), b.Conn, Message{
					For:         LEFT,
					UserId:      client.Id,
					MessageType: websocket.MessageText,
					Body:        "An user left with id : " + client.Id,
				})

				if err != nil {
					log.Println("Error on user left : ", err.Error())
				}
			}

			break
		case message := <-r.Broadcast:
			for b := range r.Clients {
				err := wsjson.Write(context.Background(), b.Conn, Message{
					For:         MSG,
					UserId:      message.UserId,
					MessageType: websocket.MessageText,
					Body:        message.Body,
				})

				log.Println("an user send message")

				if err != nil {
					log.Println("Error on user broadcast : ", err.Error())
				}
			}
		}
	}
}

package ws

import (
	"context"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Client struct {
	Id              string
	RoomId          string
	*websocket.Conn // add pool connection later
	*Pool
}

func (c *Client) Read(ctx context.Context) {
	// TODO: unregister/delete user from connection pool
	defer func() {
		c.Pool.Left <- c
	}()
	defer c.Conn.Close(websocket.StatusNormalClosure, "Connection has closed normally")

	// read message from all user
	var message Message
	for {
		err := wsjson.Read(ctx, c.Conn, &message)
		if err != nil {
			log.Println("Error on read message : ", err.Error())

			break
		}

		//  broadcast message
		c.Pool.Broadcast <- message
	}
}

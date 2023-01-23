package ws

import "nhooyr.io/websocket"

type Message struct {
	UserId                string `json:"user_id"`
	For                   string `json:"for"`
	websocket.MessageType `json:"message_type"`
	Body                  any `json:"body"`
}

const (
	MSG  = "message"
	JOIN = "join"
	LEFT = "left"
)

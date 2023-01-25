package ws

type Pool struct {
	Rooms     map[string]*Room
	Join      chan *Client
	Left      chan *Client
	Broadcast chan Message
}

func NewPool() *Pool {
	return &Pool{
		Rooms:     map[string]*Room{},
		Join:      make(chan *Client),
		Left:      make(chan *Client),
		Broadcast: make(chan Message),
	}
}

func (p Pool) Run() {
	for {
		select {
		case client := <-p.Join:
			// check if room is exist
			if _, ok := p.Rooms[client.RoomId]; !ok {
				p.Rooms[client.RoomId] = NewRoom()
				go p.Rooms[client.RoomId].Run()
			}

			// add user to room
			p.Rooms[client.RoomId].Join <- client
			break
		case client := <-p.Left:
			// remove user from room
			p.Rooms[client.RoomId].Left <- client
			// check if client is empty, remove room from pool
			if len(p.Rooms[client.RoomId].Clients) == 0 {
				delete(p.Rooms, client.RoomId)
			}

			break
		case message := <-p.Broadcast:
			// check if room is exist
			if _, ok := p.Rooms[message.RoomId]; !ok {
				return // TODO; send error message to client
			}

			// broadcast message to all user in room
			p.Rooms[message.RoomId].Broadcast <- message
		}
	}
}

package app

// TODO Insert to Redis
var AllRoom map[uint32]*Room

// Room control the set of active clients and broadcasts messages to the clients.
type Room struct {
	clients map[*Client]bool

	broadcast chan []byte

	// Register requests from the clients.Like enter the room.
	register chan *Client

	// Unregister requests from clients.Like leave the room.
	unregister chan *Client
}

// Init Room
func init() {
	AllRoom = make(map[uint32]*Room)
}

// TODO use Redis/DB to get roomInfo
func GetRoom(id uint32) *Room {
	if room, ok := AllRoom[id]; ok {
		return room
	} else {
		room = newRoom()
		go room.run()
		AllRoom[id] = room
		return room
	}
}

func newRoom() *Room {
	return &Room{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// TODO Start live commit io pool.
func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.msg)
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.msg <- message:
				default:
					close(client.msg)
					delete(r.clients, client)
				}
			}
		}
	}
}

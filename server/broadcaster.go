package server

import (
	"encoding/json"
	"fmt"
	"rtdocs/entities"
)

// util to broadcast functions to different go connections
type Broadcaster struct {
	// currently active Clients
	Clients map[*WebsocketConn]bool

	// the Broadcast channel that'll receive messages
	Broadcast chan *entities.Payload

	Add    chan *WebsocketConn
	Remove chan *WebsocketConn
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		Clients:   make(map[*WebsocketConn]bool),
		Broadcast: make(chan *entities.Payload),
		Add:       make(chan *WebsocketConn),
		Remove:    make(chan *WebsocketConn),
	}
}

func (b *Broadcaster) Run() {
	for {
		select {
		case msg := <-b.Broadcast:
			fmt.Println(msg.X, msg.Y, msg.Uuid)
			jsonPayload, err := json.Marshal(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			for client := range b.Clients {
				if client.uuid != msg.Uuid {
					client.sendMsg(jsonPayload)
				}
			}
			fmt.Println("incoming message sent to ", len(b.Clients)-1)
		case newConn := <-b.Add:
			b.Clients[newConn] = true
			fmt.Println(newConn.tcpConn.RemoteAddr().String() + " with uuid: " + newConn.uuid + " connected.")
		case removeConn := <-b.Remove:
			if _, alive := b.Clients[removeConn]; alive {
				delete(b.Clients, removeConn)
				fmt.Println(removeConn.tcpConn.RemoteAddr().String() + " left.")
			}
		}
	}
}

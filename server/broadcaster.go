package server

import (
	"encoding/json"
	"fmt"
	"lostinchaos/entities"
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
			fmt.Println(msg.Type)
			jsonPayload, err := json.Marshal(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			for client := range b.Clients {
				if client.uuid != "asdf" {
					client.sendMsg(jsonPayload)
				}
			}
			fmt.Println("incoming message sent to ", len(b.Clients)-1)
		case newConn := <-b.Add:
			b.Clients[newConn] = true
			fmt.Println(newConn.tcpConn.RemoteAddr().String() + " with uuid: " + newConn.uuid + " connected.")
			joinNotif := entities.JoinNotif{Uuid: newConn.uuid, Type: "new_join"}
			jsonPayload, err := json.Marshal(joinNotif)
			if err != nil {
				return
			}
			// notify all clients someone joined
			for client := range b.Clients {
				if client.uuid != newConn.uuid {
					client.sendMsg(jsonPayload)
				}
			}
			// also send the connected client all existing clients and their x y values
		case removeConn := <-b.Remove:
			if _, alive := b.Clients[removeConn]; alive {
				leaveNotif := entities.JoinNotif{Uuid: removeConn.uuid, Type: "left"}

				delete(b.Clients, removeConn)
				fmt.Println(removeConn.tcpConn.RemoteAddr().String() + " left.")

				jsonPayload, err := json.Marshal(leaveNotif)
				if err != nil {
					return
				}
				fmt.Println(string(jsonPayload))
				for client := range b.Clients {
					client.sendMsg(jsonPayload)
				}

			}
			// notify all clients someone left
		}
	}
}

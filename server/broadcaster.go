package server

import "fmt"

// util to broadcast functions to different go connections
type Broadcaster struct {
	// currently active Clients
	Clients map[*WebsocketConn]bool

	// the Broadcast channel that'll receive messages
	Broadcast chan *Msg

	Add    chan *WebsocketConn
	Remove chan *WebsocketConn
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		Clients:   make(map[*WebsocketConn]bool),
		Broadcast: make(chan *Msg),
		Add:       make(chan *WebsocketConn),
		Remove:    make(chan *WebsocketConn),
	}
}

func (b *Broadcaster) Run() {
	for {
		select {
		case msg := <-b.Broadcast:
			fmt.Println(string(msg.Data))
			for client := range b.Clients {
				if client != msg.WConn {
					client.sendMsg(msg.Data)
				}
			}
			fmt.Println("incoming message sent to ", len(b.Clients))
		case newConn := <-b.Add:
			b.Clients[newConn] = true
			fmt.Println(newConn.tcpConn.RemoteAddr().String() + " connected.")
		case removeConn := <-b.Remove:
			if _, alive := b.Clients[removeConn]; alive {
				delete(b.Clients, removeConn)
				fmt.Println(removeConn.tcpConn.RemoteAddr().String() + " left.")
			}
		}
	}
}

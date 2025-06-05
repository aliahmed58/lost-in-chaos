package main

import (
	"context"
	"lostinchaos/server"
	"net"
	"net/http"
)

func SaveConnInContext(ctx context.Context, c net.Conn) context.Context {

	return context.WithValue(ctx, server.ConnContextKey, c)
}

func main() {

	server.Init()
	broadcast := server.NewBroadcaster()
	go broadcast.Run()

	server.R.HandleFunc("/", server.HandleHomePage)
	server.R.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		server.HandleWebsocketConn(broadcast, w, r)
	})

	s := http.Server{
		Addr:        ":80",
		ConnContext: SaveConnInContext,
		Handler:     server.R,
	}

	panic(s.ListenAndServe())

}

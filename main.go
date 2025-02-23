package main

import (
	"context"
	"net"
	"net/http"
	"rtdocs/server"
)

func SaveConnInContext(ctx context.Context, c net.Conn) context.Context {

	return context.WithValue(ctx, server.ConnContextKey, c)
}

func main() {

	server.Init()

	server.R.HandleFunc("/", server.HandleHomePage)
	server.R.HandleFunc("/websocket", server.HandleWebsocketConn)
	server.R.HandleFunc("/test", server.HandleTest)

	s := http.Server{
		Addr:        ":80",
		ConnContext: SaveConnInContext,
		Handler:     server.R,
	}

	panic(s.ListenAndServe())
}

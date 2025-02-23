package server

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"html/template"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

var R *mux.Router
var Templ *template.Template

type contextKey struct {
	key string
}

var ConnContextKey = &contextKey{key: "http-conn"}

func Init() {
	R = mux.NewRouter()
	Templ = template.Must(Templ.ParseFiles(
		"templates/index.html",
	))
}

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	Templ.Execute(w, nil)
}

func HandleWebsocketConn(w http.ResponseWriter, r *http.Request) {
	cHeaders, err := InitClientHeaders(r)
	if err != nil {
		//TODO: send a bad response
	}
	// fmt.Println(cHeaders.Headers)
	acceptVal := cHeaders.Headers["Sec-WebSocket-Key"]
	acceptVal += "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hasher := sha1.New()
	hasher.Write([]byte(acceptVal))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// tcpConn := r.Context().Value(http.LocalAddrContextKey).(*net.TCPAddr)
	conn, _, err := http.NewResponseController(w).Hijack()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(conn.RemoteAddr().String())
	// conn.Write([]byte("HTTP/1.1 101 Switching Protocols\r\n"))
	// go handleConn(conn)

	header := "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: " + sha + "\r\n\r\n"

	_, err1 := conn.Write([]byte(header))
	if err1 != nil {
		fmt.Println(err1)
	} // w.Header().Set("Upgrade", "websocket")
	// w.Header().Set("Connection", "Upgrade")
	// w.Header().Set("Sec-WebSocket-Accept", sha)
	// w.WriteHeader(http.StatusSwitchingProtocols)

	go handleConn(conn)
}

func HandleTest(w http.ResponseWriter, r *http.Request) {
	conn := r.Context().Value(ConnContextKey).(net.Conn)
	fmt.Println(conn.RemoteAddr().String())
	fmt.Fprintf(w, "ok fine damn")
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println(conn.RemoteAddr().String())
	// Read recvData from the connection
	reader := bufio.NewReader(conn)
	for {
		fmt.Println("worked till here")
		header, err := reader.ReadByte()
		fmt.Println(header)
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				fmt.Println("Error reading from connection:", err)
				break
			}
		}
		// fmt.Println(header)
	}
}

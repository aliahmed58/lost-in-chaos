package server

import (
	"bufio"
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
	conn, err := UpgradeConn(w, r)
	if err != nil {
		fmt.Println(err)
	}
	go conn.ReadMsg()
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
		header, err := reader.ReadByte()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				fmt.Println("Error reading from connection:", err)
				break
			}
		}
		fin := header & 1
		opcode := header & 0b00001111
		// fmt.Println(fin)
		fmt.Println(fin)
		fmt.Println(opcode)

		next, err := reader.ReadByte()
		if err != nil {
			fmt.Println(err)
			break
		}
		mask := next & 1
		length := next & 0b01111111
		fmt.Println(mask, length)

		maskData := make([]byte, 4)
		_, err1 := reader.Read(maskData)
		if err1 != nil {
			fmt.Println(err)
		}

		fmt.Printf("%08b", maskData)

		payloadData := make([]byte, length)
		_, err2 := reader.Read(payloadData)
		if err2 != nil {
			fmt.Println(err2)
		}
		decodedData := make([]byte, length)
		for i, d := range payloadData {
			decodedData[i] = d ^ maskData[i%4]
		}
		fmt.Println(string(decodedData))
		// fmt.Println(header)
	}
}

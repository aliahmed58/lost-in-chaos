package server

import (
	"fmt"
	"html/template"
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

func HandleWebsocketConn(b *Broadcaster, w http.ResponseWriter, r *http.Request) {
	conn, err := UpgradeConn(b, w, r)
	if err != nil {
		fmt.Println(err)
	}
	go conn.ReadMsg()
}

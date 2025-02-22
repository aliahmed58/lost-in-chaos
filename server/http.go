package server

import (
	"fmt"
	"net/http"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my websockets website!")
}

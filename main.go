package main

import (
	"fmt"
	"net/http"
	"os"
	"rtdocs/server"
)

func main() {

	http.HandleFunc("/", server.HandleHomePage)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

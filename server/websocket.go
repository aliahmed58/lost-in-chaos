package server

import "net/http"

var clientHeaders = []string{
	"Host",
	"Upgrade",
	"Connection",
	"Sec-WebSocket-Key",
	"Sec-WebSocket-Version",
}

type ClientHeaders struct {
	Headers map[string]string
}

func InitClientHeaders(r *http.Request) (ClientHeaders, error) {
	c := ClientHeaders{Headers: make(map[string]string)}
	for _, header := range clientHeaders {
		c.Headers[header] = r.Header.Get(header)
		if c.Headers[header] == "" {
			// handle and throw error TODO: learn how to throw custom errors
		}
		// TODO: if any header is wrong or misunderstood
	}
	return c, nil
}

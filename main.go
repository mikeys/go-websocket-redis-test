package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	log.Printf("%v", *addr)
	flag.Parse()
	go h.run()
	http.Handle("/", websocket.Handler(wsHandler))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

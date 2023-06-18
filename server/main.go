package main

import (
	"flag"
	"log"
	"net/http"
	"server/lib"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {

	var (
		entering = make(chan lib.Client)
		leaving  = make(chan lib.Client)
		messages = make(chan []byte) // all incoming client messages
	)

	// broadcaster is a goroutine that broadcasts messages to all clients
	broadcaster := lib.NewBroadcaster()
	go broadcaster.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// handle websocket
		broadcaster.HandleWebsocket(w, r, entering, leaving, messages)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

package lib

import (
	"encoding/json"
	"log"
	"net/http"
	"server/models"
	"time"
)

type Broadcaster interface {
	Run()
	HandleWebsocket(w http.ResponseWriter, r *http.Request, entering chan Client, leaving chan Client, messages chan []byte)
}

type broadcaster struct {
	entering chan Client
	leaving  chan Client
	messages chan []byte
}

func NewBroadcaster() Broadcaster {
	return &broadcaster{
		entering: make(chan Client),
		leaving:  make(chan Client),
		messages: make(chan []byte),
	}
}

func (b *broadcaster) Run() {
	clients := make(map[Client]bool) // all connected clients
	for {
		select {
		case msg := <-b.messages: // incoming message
			// broadcast incoming message to all clients' outgoing message channels
			for cli := range clients {
				cli.send <- msg
			}
		case cli := <-b.entering: // incoming client
			clients[cli] = true
		case cli := <-b.leaving: // leaving client
			delete(clients, cli)
			close(cli.send)
		}
	}
}

func (b *broadcaster) HandleWebsocket(w http.ResponseWriter, r *http.Request, entering chan Client, leaving chan Client, messages chan []byte) {
	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// create client
	client := &Client{conn: conn, send: make(chan []byte, 256), broadcaster: b}
	// register client
	b.entering <- *client

	enterMessage := &models.Message{
		User:      conn.RemoteAddr().String(),
		Body:      conn.RemoteAddr().String() + " entered",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Type:      "enter",
	}

	jsonBytes, _ := json.Marshal(enterMessage)
	b.messages <- jsonBytes

	// handle all reads
	go client.readPump()

	// handle all writes
	go client.writePump()

}

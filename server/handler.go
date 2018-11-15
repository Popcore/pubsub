package server

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
)

type pubSubHandler struct{
	subscribers map[chan []byte]bool
	openConn chan chan []byte
	closeConn chan chan []byte
	message chan string
}

func (p *pubSubHandler) Listen() {
	for {
		select {
		case s := <- p.openConn:
			log.Println("subscriber added")
			p.subscribers[s] = true
		case s := <- p.closeConn:
			log.Println("subscriber removed")
			delete(p.subscribers, s)
		case s := <- p.message:
			for subChannel, _ := range p.subscribers {
				subChannel <- []byte(s)
			}
		}
	}
}

func (p pubSubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		PostMessageHandler(w, r, &p)
	} else {
		SubscribeHandler(w, r, &p)
	}
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request, p *pubSubHandler) {

	// each connection has its own channel
	msgChannel := make(chan []byte)
	p.openConn <- msgChannel

	flusher, _ := w.(http.Flusher)

	defer func() {
		p.closeConn <- msgChannel
	}()

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<- notify
		p.closeConn <- msgChannel
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// start listening to incoming messages
	for {
		fmt.Fprintf(w, "%s \n", <-msgChannel)
		flusher.Flush()
	}
}

// message is the representation of a JSON message that can
// be broadcasted to subscribers.
type message struct {
	Body string `json:"message"`
}

func PostMessageHandler(w http.ResponseWriter, r *http.Request, p *pubSubHandler) {
	decoder := json.NewDecoder(r.Body)
	msg := message{}
	err := decoder.Decode(&msg)

	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		p.message <- msg.Body
	}()

	// print connected clients
	fmt.Fprintf(w, "Connected clients: %d", len(p.subscribers))
}
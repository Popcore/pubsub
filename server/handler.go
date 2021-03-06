package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// pubSubManager manages subscribers and their connections.
type PubSubManager struct {
	subscribers map[chan []byte]bool
	openConn    chan chan []byte
	closeConn   chan chan []byte
	mux         sync.Mutex
}

func (p PubSubManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		postMessageHandler(w, r, &p)
	} else {
		subscribeHandler(w, r, &p)
	}
}

func subscribeHandler(w http.ResponseWriter, r *http.Request, p *PubSubManager) {
	msgChannel := make(chan []byte)

	p.mux.Lock()
	p.subscribers[msgChannel] = true
	p.mux.Unlock()

	// use the http flusher interface for sse
	flusher, _ := w.(http.Flusher)

	// set headers
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/event-stream")

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		p.mux.Lock()
		delete(p.subscribers, msgChannel)
		p.mux.Unlock()
	}()

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

func postMessageHandler(w http.ResponseWriter, r *http.Request, p *PubSubManager) {
	decoder := json.NewDecoder(r.Body)
	msg := message{}
	err := decoder.Decode(&msg)

	if err != nil {
		log.Println(err)
		return
	}

	// publish message
	go func() {
		for subChannel, _ := range p.subscribers {
			subChannel <- []byte(msg.Body)
		}
	}()

	// print connected clients
	fmt.Fprintf(w, "Connected clients: %d", len(p.subscribers))
}

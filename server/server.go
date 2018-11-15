package server

import (
	"log"
	"net/http"
)

// System is a custom type used to hold information about
// the application
type System struct {
	Subscribers []int
	PubSubChannel chan string
}

// Server is the custom type used to handle http connections
type Server struct {
	Address string
	Mux *http.ServeMux
	System *System
}

// New server returns a concrete Server instance
func NewServer(addr string) *Server {
	system := System{
		Subscribers: []int{},
		PubSubChannel: make(chan string),
	}

	p := pubSubHandler{
		subscribers: make(map[chan []byte]bool),
		closeConn: make(chan chan []byte),
		openConn: make(chan chan []byte),
		message: make(chan string),
	}
	mux := http.NewServeMux()
	mux.Handle("/", p)

	go p.Listen()

	return &Server{
		Address: addr,
		Mux: mux,
		System: &system,
	}
}

// Start runs the application http server
func (s Server) Start() {
	server := &http.Server{
		Addr: s.Address,
		Handler: s.Mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
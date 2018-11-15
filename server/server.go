package server

import (
	"log"
	"net/http"
)

// Server is the custom type used to handle http connections
type Server struct {
	Address string
	Mux *http.ServeMux
}

// New server returns a concrete Server instance
func NewServer(addr string) *Server {

	psManager := PubSubManager{
		subscribers: make(map[chan []byte]bool),
		closeConn: make(chan chan []byte),
		openConn: make(chan chan []byte),
	}

	mux := http.NewServeMux()
	mux.Handle("/", psManager)

	return &Server{
		Address: addr,
		Mux: mux,
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
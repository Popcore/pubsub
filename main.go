package main

import (
	"github.com/popcore/pubsub/server"
)

func main() {
	s := server.NewServer(":9090")
	s.Start()
}

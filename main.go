package main

import (
	"github.com/popcore/jaakpubsub/server"
)

func main() {
	s := server.NewServer(":9090")
	s.Start()
}
package main

import (
	"flag"

	"raft-sister/src/core"
)

func main() {
	var port int
	var hostfile string

	flag.IntVar(&port, "port", 50000, "The port on which the server will listen")
	flag.StringVar(&hostfile, "hostfile", "./config/Hostfile", "The port on which the server will listen")
	flag.Parse()

	var node = core.NewNode(port, nil)

	node.InitServer(hostfile)
	node.SanityCheck()

	for {
	}
}

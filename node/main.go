package main

import (
	node "Node/lib"
	"flag"
)

func main() {
	var addr string
	var timeoutAvg int
	var hostfile string

	flag.StringVar(&addr, "addr", "0.0.0.0:60000", "The port on which the server will listen")
	flag.IntVar(&timeoutAvg, "timeout", 200, "The average time before a node declares a timeout")
	flag.StringVar(&hostfile, "hostfile", "/config/Hostfile", "The list of all servers in the cluster")
	flag.Parse()

	node := node.NewNode(addr)
	node.Init(".."+hostfile, timeoutAvg)

	for node.Running {
	}
}

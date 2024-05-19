package main

import (
	"flag"

	"raft-sister/src/core"
)

func main() {
	var addr string
	var timeoutAvg int
	var hostfile string

	flag.StringVar(&addr, "addr", "0.0.0.0:60000", "The port on which the server will listen")
	flag.IntVar(&timeoutAvg, "timeout", 15, "The average time before a node declares a timeout")
	flag.StringVar(&hostfile, "hostfile", "./config/Hostfile", "The list of all servers in the cluster")
	flag.Parse()

	var node = core.NewNode(addr, timeoutAvg, nil)

	node.Init(hostfile)
	node.CheckSanity()

	for node.Running {
	}
}

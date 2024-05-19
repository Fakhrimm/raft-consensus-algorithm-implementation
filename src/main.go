package main

import (
	"flag"

	"raft-sister/src/core"
)

func main() {
	var port int
	var timeoutAvg int
	var hostfile string

	flag.IntVar(&port, "port", 50000, "The port on which the server will listen")
	flag.IntVar(&timeoutAvg, "timeout", 15, "The average time before a node declares a timeout")
	flag.StringVar(&hostfile, "hostfile", "./config/Hostfile", "The list of all servers in the cluster")
	flag.Parse()

	var node = core.NewNode(port, timeoutAvg, nil)

	node.InitServer(hostfile)
	node.ReadServerList(hostfile)
	node.CheckSanity()

	for node.Running {
	}
}

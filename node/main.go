package main

import (
	node "Node/lib"
	"flag"
	"log"
)

func main() {
	var addr string
	var timeoutAvg int
	var hostfile string
	var hostfileNew string
	var isJointConsensus bool

	flag.StringVar(&addr, "addr", "0.0.0.0:60000", "The port on which the server will listen")
	flag.IntVar(&timeoutAvg, "timeout", 200, "The average time before a node declares a timeout")
	flag.StringVar(&hostfile, "hostfile", "/config/Hostfile", "The list of all servers in the cluster")
	flag.StringVar(&hostfileNew, "hostfilenew", "", "The list of new servers in the cluster")
	flag.BoolVar(&isJointConsensus, "isjointconsensus", false, "condition to transition to new configuration")
	flag.Parse()

	newNode := node.NewNode(addr)
	log.Printf("hostFile: %v", ".."+hostfile)
	log.Printf("hostFileNew: %v", ".."+hostfileNew)
	log.Printf("isJointConsensus: %v", isJointConsensus)
	newNode.Init(".."+hostfile, timeoutAvg, ".."+hostfileNew, isJointConsensus)

	for newNode.Running {
	}

	newNode.Status()
	newNode.SaveLogs()

	log.Println("Program is stopped")
}

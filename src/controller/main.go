package main

import controller "raft-sister/src/controller/Controller"

func main() {
	controller := controller.NewController()
	controller.Run()
}

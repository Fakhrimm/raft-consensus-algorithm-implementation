package main

import controller "Controller/lib"

func main() {
	controller := controller.NewController()
	controller.Run()
}

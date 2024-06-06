package controller

import (
	"Controller/grpc/comm"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Controller struct {
	// General Attribute
	Running bool

	// Grpc purposes
	client comm.CommServiceClient

	// Application purposes
	reader bufio.Reader
}

func NewController() *Controller {
	controller := &Controller{
		reader: *bufio.NewReader(os.Stdin),
	}
	return controller
}

func (c *Controller) Run() {
	c.Running = true

	fmt.Println("Node Controller Terminal, Enter command below")

	for c.Running {
		fmt.Print("> ")
		input, err := c.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
			continue
		}

		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")
		length := len(parts)

		if length <= 0 {
			continue
		}

		switch strings.ToLower(parts[0]) {
		case "exit":
			log.Println("Exiting...")
			c.Running = false
			return
		case "ping":
			if length < 2 {
				log.Println("Usage: ping <address>")
				continue
			}
			c.Ping(parts[1])
		case "status":
			if length < 2 {
				log.Println("Usage: status <address>")
				continue
			}
			c.Status(parts[1])
		case "get":
			if length < 3 {
				log.Println("Usage: get <address> <key>")
				continue
			}
			c.Get(parts[1], parts[2])
		case "set":
			if length < 4 {
				log.Println("Usage: get <address> <key> <value>")
				continue
			}
			c.Set(parts[1], parts[2], parts[3])
		case "strlen":
			if length < 3 {
				log.Println("Usage: strlen <address> <key>")
				continue
			}
			c.Strlen(parts[1], parts[2])
		case "delete":
			if length < 3 {
				log.Println("Usage: strlen <address> <key>")
				continue
			}
			c.Delete(parts[1], parts[2])
		case "append":
			if length < 4 {
				log.Println("Usage: append <address> <key> <value>")
				continue
			}
			c.Append(parts[1], parts[2], parts[3])
		case "stop":
			if length < 2 {
				log.Println("Usage: stop <address>")
				continue
			}
			c.Stop(parts[1])
		case "config":
			if length < 3 {
				log.Println("Usage: config <address> <config file>")
				continue
			}
			c.Config(parts[1], parts[2])
		default:
			log.Println("Unknown Command: ", parts[0])
		}
	}
}

func (c *Controller) Call(address string, callable func()) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	c.client = comm.NewCommServiceClient(conn)
	callable()

	conn.Close()
}

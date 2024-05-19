package controller

import (
	"Controller/grpc/comm"
	"bufio"
	"context"
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

		switch parts[0] {
		case "exit":
			fmt.Println("Exiting...")
			c.Running = false
			return
		case "ping":
			if length < 2 {
				fmt.Println("Usage: ping <address>")
				continue
			}
			c.Call(parts[1], func() {
				response_ping, err_ping := c.client.Ping(context.Background(), &comm.PingRequest{})
				if err_ping != nil {
					log.Printf("ping error: %v", err_ping)
				}
				fmt.Println(response_ping.Message)
			})
		case "get":
			if length < 3 {
				fmt.Println("Usage: get <address> <key>")
				continue
			}
			c.Call(parts[1], func() {
				response_ping, err_ping := c.client.RequestValue(context.Background(), &comm.RequestValueRequest{Key: parts[2]})
				if err_ping != nil {
					log.Printf("get error: %v", err_ping)
				}
				fmt.Println(response_ping.Value)
			})
		case "set":
			if length < 4 {
				fmt.Println("Usage: get <address> <key> <value>")
				continue
			}
			c.Call(parts[1], func() {
				response_ping, err_ping := c.client.SetValue(context.Background(), &comm.SetValueRequest{Key: parts[2], Value: parts[3]})
				if err_ping != nil {
					log.Printf("set error: %v", err_ping)
				}
				fmt.Println(response_ping.Message)
			})
		case "stop":
			if length < 2 {
				fmt.Println("Usage: stop <address>")
				continue
			}
			c.Call(parts[1], func() {
				_, err_ping := c.client.Stop(context.Background(), &comm.PingRequest{})

				fmt.Println("Server response:", err_ping)
			})

		default:
			fmt.Println("Unknown Command: ", parts[0])
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

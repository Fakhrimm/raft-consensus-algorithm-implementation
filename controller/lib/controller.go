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
			log.Println("Exiting...")
			c.Running = false
			return
		case "ping":
			if length < 2 {
				log.Println("Usage: ping <address>")
				continue
			}
			c.Call(parts[1], func() {
				response, err := c.client.Ping(context.Background(), &comm.BasicRequest{})
				if err != nil {
					log.Printf("ping error: %v", err)
				}
				log.Printf("Response: {%v}", response)
			})
		case "get":
			if length < 3 {
				log.Println("Usage: get <address> <key>")
				continue
			}
			c.Call(parts[1], func() {
				response, err := c.client.GetValue(context.Background(), &comm.GetValueRequest{Key: parts[2]})
				if err != nil {
					log.Printf("get error: %v", err)
				}
				log.Printf("Response: {%v}", response)
			})
		case "set":
			if length < 4 {
				log.Println("Usage: get <address> <key> <value>")
				continue
			}
			c.Call(parts[1], func() {
				response, err := c.client.SetValue(context.Background(), &comm.SetValueRequest{Key: parts[2], Value: parts[3]})
				if err != nil {
					log.Printf("set error: %v", err)
				}
				log.Printf("Response: {%v}", response)
			})
		case "strlen":
			if length < 3 {
				log.Println("Usage: strlen <address> <key>")
				continue
			}
			c.Call(parts[1], func() {
				response, err := c.client.StrlnValue(context.Background(), &comm.StrlnValueRequest{Key: parts[2]})
				if err != nil {
					log.Printf("strlen error: %v", err)
				}
				log.Printf("Response: {%v}", response)
			})
		case "delete":
			if length < 3 {
				log.Println("Usage: strlen <address> <key>")
				continue
			}
			c.Call(parts[1], func() {
				response, err := c.client.DeleteValue(context.Background(), &comm.DeleteValueRequest{Key: parts[2]})
				if err != nil {
					log.Printf("delete error: %v", err)
				}
				log.Printf("Response: {%v}", response)
			})
		case "append":
			if length < 4 {
				log.Println("Usage: append <address> <key> <value>")
				continue
			}
			c.Call(parts[1], func() {
				response, err := c.client.AppendValue(context.Background(), &comm.AppendValueRequest{Key: parts[2], Value: parts[3]})
				if err != nil {
					log.Printf("append error: %v", err)
				}
				log.Printf("Response: {%v}", response)
			})
		case "stop":
			if length < 2 {
				log.Println("Usage: stop <address>")
				continue
			}
			c.Call(parts[1], func() {
				_, err := c.client.Stop(context.Background(), &comm.BasicRequest{})

				log.Println("Server response:", err)
			})

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

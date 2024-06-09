package controller

import (
	"Controller/grpc/comm"
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strings"
)

func (c *Controller) Ping(address string) {
	c.Call(address, func() {
		response, err := c.client.Ping(context.Background(), &comm.BasicRequest{})
		if err != nil {
			log.Printf("ping error: %v", err)
		}
		log.Printf("Response: {%v}", response)
	})
}

func (c *Controller) Status(address string) {
	c.Call(address, func() {
		response, err := c.client.Status(context.Background(), &comm.BasicRequest{})
		if err != nil {
			log.Printf("ping error: %v", err)
		}
		log.Printf("Response: {%v}", response)
	})
}

func (c *Controller) Stop(address string) {
	c.Call(address, func() {
		_, err := c.client.Stop(context.Background(), &comm.BasicRequest{})

		log.Println("Server response:", err)
	})
}

func (c *Controller) Get(address string, key string) {
	c.Call(address, func() {
		response, err := c.client.GetValue(context.Background(), &comm.GetValueRequest{Key: key})
		if err != nil {
			log.Printf("get error: %v", err)
		}
		log.Printf("Response: {%v}", response)
	})
}

func (c *Controller) GetLogs(address string) {
	c.Call(address, func() {
		response, err := c.client.GetLogs(context.Background(), &comm.BasicRequest{})
		if err != nil {
			log.Printf("getlogs error: %v", err)
		}
		log.Printf("Response: {%v}", response)

		log.Printf("/\n\nlog info:")
		for index, entry := range response.Value {
			log.Printf("[%v]: command: %v, key:%v, value:%v, term:%v", index, entry.Command, entry.Key, entry.Value, entry.Term)
		}

	})
}

func (c *Controller) Strlen(address string, key string) {
	c.Call(address, func() {
		response, err := c.client.StrlnValue(context.Background(), &comm.StrlnValueRequest{Key: key})
		if err != nil {
			log.Printf("strlen error: %v", err)
		}
		log.Printf("Response: {%v}", response)
	})
}

func (c *Controller) Delete(address string, key string) {
	c.Call(address, func() {
		response, err := c.client.DeleteValue(context.Background(), &comm.DeleteValueRequest{Key: key})
		if err != nil {
			log.Printf("delete error: %v", err)
		}
		log.Printf("Response: {%v}", response)
	})
}

func (c *Controller) Set(address string, key string, value string) {
	c.Call(address, func() {
		response, err := c.client.SetValue(context.Background(), &comm.SetValueRequest{Key: key, Value: value})
		if err != nil {
			log.Printf("set error: %v", err)
		}
		log.Printf("Response: {%v}", response)
	})
}

func (c *Controller) Append(address string, key string, value string) {
	c.Call(address, func() {
		response, err := c.client.AppendValue(context.Background(), &comm.AppendValueRequest{Key: key, Value: value})
		if err != nil {
			log.Printf("append error: %v", err)
		}
		log.Printf("Response: {%v}", response)
	})
}

func (c *Controller) Config(address string, filename string) {
	file, err := os.Open(filename)

	if err != nil {
		log.Printf("Failed to load hostfile: %v", err)

		currentDir, err := os.Getwd()
		if err != nil {
			log.Printf("Error getting current directory: %v\n", err)
			return
		}
		log.Printf("Current directory is: %v", currentDir)
		log.Printf("File attempted to be open is: %v", filename)

		return
	}
	defer file.Close()

	serverList := ""
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		_, err := net.ResolveTCPAddr("tcp", line)
		if err != nil {
			continue
		}
		serverList += line + ","

	}

	if len(serverList) > 0 {
		serverList = strings.TrimSuffix(serverList, ",")
	}

	log.Printf("Sending new server configuration: %v", serverList)
	c.Call(address, func() {
		response, err := c.client.ChangeMembership(context.Background(), &comm.ChangeMembershipRequest{
			NewClusterAddresses: serverList,
		})

		if err != nil {
			log.Printf("Config error: %v", err)
		}
		log.Println("Server response:", response)
	})
}

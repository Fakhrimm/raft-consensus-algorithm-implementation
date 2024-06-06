package controller

import (
	"Controller/grpc/comm"
	"context"
	"log"
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

func (c *Controller) Config(address string, file string) {
	c.Call(address, func() {
		_, err := c.client.Stop(context.Background(), &comm.BasicRequest{})

		log.Println("Server response:", err)
	})
}

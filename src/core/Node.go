package core

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"raft-sister/src/proto/comm"
)

type Node struct {
	// Core components
	Server  *grpc.Server
	Client  comm.CommServiceClient
	Port    int
	Running bool

	// Cluster components
	ServerList  []net.TCPAddr
	ServerCount int

	// Raft protocol
	State      ServerState
	timeoutAvg int

	// Functional components
	Map map[string]string
}

func NewNode(port int, timeoutAvg int, serverList []net.TCPAddr) *Node {
	if serverList == nil {
		serverList = []net.TCPAddr{}
	}
	node := &Node{
		Port:       port,
		ServerList: serverList,
		Map:        make(map[string]string),
		timeoutAvg: timeoutAvg,
	}
	return node
}

func (n *Node) InitServer(hostfile string) {
	log.Printf("Initializing node")

	n.Running = true
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(n.Port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := server{Node: n}
	n.Server = grpc.NewServer()
	comm.RegisterCommServiceServer(n.Server, &s)

	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err := n.Server.Serve(lis); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()
}

func (n *Node) ReadServerList(filename string) []net.TCPAddr {
	log.Printf("Initializing server list")

	var serverList []net.TCPAddr

	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Failed to load hostfile")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		log.Print(line)

		addr, err := net.ResolveTCPAddr("tcp", line)
		if err != nil {
			log.Fatalf("Invalid address: %v", line)
		}
		serverList = append(serverList, *addr)
	}

	n.ServerList = serverList
	n.ServerCount = len(serverList)

	return n.ServerList
}

func (n *Node) Call(address string, callable func()) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	n.Client = comm.NewCommServiceClient(conn)
	callable()

	conn.Close()
}

func (n *Node) CheckSanity() {
	key := "__CheckSanity__"

	oldVal := n.Map[key]
	n.Call("0.0.0.0:"+strconv.Itoa(n.Port), func() {
		response_test, err_test := n.Client.Ping(context.Background(), &comm.PingRequest{})
		if err_test != nil {
			log.Printf("Sanity check error: %v", err_test)
		}
		log.Printf("Test Response: %v", response_test.Message)

		response_set, err_set := n.Client.SetValue(context.Background(), &comm.SetValueRequest{Key: key, Value: "OK"})
		if err_set != nil {
			log.Printf("Sanity check error: %v", err_set)
		}
		log.Printf("Set Response: %v", response_set.Message)

		response_get, err_get := n.Client.RequestValue(context.Background(), &comm.RequestValueRequest{Key: key})
		if err_get != nil {
			log.Printf("Sanity check error: %v", err_get)
		}
		log.Printf("Get Response: %v", response_get.Value)
	})

	n.Map[key] = oldVal
}

func (n *Node) Stop() {
	n.Running = false
}

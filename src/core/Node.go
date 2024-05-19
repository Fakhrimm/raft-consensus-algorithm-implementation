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

	pb "raft-sister/src/proto/example"
)

type Node struct {
	// Core components
	Server *grpc.Server
	Client pb.ExampleServiceClient
	Port   int

	// Cluster components
	ServerList  []net.TCPAddr
	ServerCount int
}

func NewNode(port int, serverList []net.TCPAddr) *Node {
	if serverList == nil {
		serverList = []net.TCPAddr{}
	}
	node := &Node{
		Port:       port,
		ServerList: serverList,
	}
	return node
}

func (n *Node) InitServer(hostfile string) {
	log.Printf("Initializing node")
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(n.Port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := server{Node: n}
	n.Server = grpc.NewServer()
	pb.RegisterExampleServiceServer(n.Server, &s)

	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err := n.Server.Serve(lis); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()
}

func (n *Node) ReadServerList(filename string) {
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

		addr, err := net.ResolveTCPAddr("tcp", line)
		if err != nil {
			log.Fatalf("Invalid address: %v", line)
		}
		serverList = append(serverList, *addr)
	}

	n.ServerList = serverList
	n.ServerCount = len(serverList)
}

func (n *Node) Call(address string, callable func()) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	n.Client = pb.NewExampleServiceClient(conn)
	callable()

	conn.Close()
}

func (n *Node) SanityCheck() {
	n.Call("0.0.0.0:"+strconv.Itoa(n.Port), func() {
		response, err := n.Client.SayHello(context.Background(), &pb.HelloRequest{Name: strconv.Itoa(n.Port)})
		if err != nil {
			log.Printf("Sanity check error: %v", err)
		}
		log.Printf("Response: %v", response.GetMessage())
	})
}

package Node

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"raft-sister/src/proto/comm"
)

type Node struct {
	// Core components
	server  *grpc.Server
	client  comm.CommServiceClient
	Running bool
	addr    net.TCPAddr

	// Cluster components
	serverList  []net.TCPAddr
	serverCount int

	// Raft protocol
	state      ServerState
	id         int
	timeoutAvg int

	// Functional components
	Map map[string]string
}

func NewNode(addr string, timeoutAvg int, serverList []net.TCPAddr) *Node {
	if serverList == nil {
		serverList = []net.TCPAddr{}
	}
	resolved, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		log.Fatalf("Invalid address of %v set to node", addr)
	}

	node := &Node{
		addr:       *resolved,
		serverList: serverList,
		Map:        make(map[string]string),
		timeoutAvg: timeoutAvg,
	}
	return node
}

func (n *Node) Init(hostfile string) {
	log.Printf("Initializing node")
	n.InitServer()
	n.ReadServerList(hostfile)
	log.Printf("Node initialization complete with address %v and id %v", n.addr.String(), n.id)
}

func (n *Node) InitServer() {
	log.Printf("Initializing grpc server")

	n.Running = true
	lis, err := net.Listen("tcp4", n.addr.String())

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := server{Node: n}
	n.server = grpc.NewServer()
	comm.RegisterCommServiceServer(n.server, &s)

	log.Printf("server set address is %v, listening at %v", n.addr.String(), lis.Addr())

	go func() {
		if err := n.server.Serve(lis); err != nil {
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
	index := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		addr, err := net.ResolveTCPAddr("tcp", line)
		if err != nil {
			log.Fatalf("Invalid address: %v", line)
		}
		if addr.String() == n.addr.String() {
			n.id = index
		}

		serverList = append(serverList, *addr)
		log.Print(addr)
		index++
	}

	n.serverList = serverList
	n.serverCount = len(serverList)

	return n.serverList
}

func (n *Node) Call(address string, callable func()) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	n.client = comm.NewCommServiceClient(conn)
	callable()

	conn.Close()
}

func (n *Node) CheckSanity() {
	key := "__CheckSanity__"

	oldVal := n.Map[key]
	log.Printf("Checking own address of %v", n.addr.String())
	n.Call(n.addr.String(), func() {
		response_test, err_test := n.client.Ping(context.Background(), &comm.PingRequest{})
		if err_test != nil {
			log.Printf("Sanity check error: %v", err_test)
		}
		log.Printf("Test Response: %v", response_test.Message)

		response_set, err_set := n.client.SetValue(context.Background(), &comm.SetValueRequest{Key: key, Value: "OK"})
		if err_set != nil {
			log.Printf("Sanity check error: %v", err_set)
		}
		log.Printf("Set Response: %v", response_set.Message)

		response_get, err_get := n.client.RequestValue(context.Background(), &comm.RequestValueRequest{Key: key})
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

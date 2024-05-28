package node

import (
	"Node/grpc/comm"
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type State int

const (
	Leader State = iota // https://go.dev/wiki/Iota
	Candidate
	Follower
)

type LogEntry struct {
	term  int
	key   string
	value string
}

type Info struct {
	// Persistent
	id               int
	leaderId         int
	currentTerm      int
	votedFor         int
	log              []LogEntry
	clusterAddresses []net.TCPAddr
	clusterCount     int
	timeoutAvgTime   int

	// Volatile
	commitIndex int
	lastApplied int

	// Leaders
	nextIndex  []int
	matchIndex []int
}

type Node struct {
	// General attribute
	Running bool
	mutex   sync.Mutex

	// Grpc purposes
	address    net.TCPAddr
	grpcServer *grpc.Server
	grpcClient comm.CommServiceClient

	// Raft purposes
	state               State
	info                Info
	electionTimer       *time.Timer
	electionResetSignal chan bool

	// Application purposes
	app Application
}

func NewNode(addr string) *Node {
	resolved, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		log.Fatalf("Invalid address of %v set to node", addr)
	}

	node := &Node{
		address:             *resolved,
		electionResetSignal: make(chan bool),
	}

	return node
}

func (node *Node) Init(hostfile string, timeoutAvgTime int) {
	log.Printf("Initializing node")
	node.InitServer()
	node.ReadServerList(hostfile)

	node.info.timeoutAvgTime = timeoutAvgTime

	node.resetElectionTimer()

	// node.CheckSanity()
	log.Printf("Node initialization complete with address %v and id %v", node.address.String(), node.info.id)
	node.Running = true
	node.state = Follower

	go node.ElectionTimerHandler()
}

func (node *Node) InitServer() {
	log.Printf("Initializing grpc server")

	lis, err := net.Listen("tcp4", node.address.String())

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := server{Node: node}
	node.grpcServer = grpc.NewServer()
	comm.RegisterCommServiceServer(node.grpcServer, &s)

	log.Printf("server set address is %v, listening at %v", node.address.String(), lis.Addr())

	go func() {
		if err := node.grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve %v", err)
			node.Running = false
		}
	}()
}

func (node *Node) ReadServerList(filename string) []net.TCPAddr {
	log.Printf("Initializing server list")

	var serverList []net.TCPAddr

	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Failed to load hostfile: %v", err)
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
		if addr.String() == node.address.String() {
			node.info.id = index
		}

		serverList = append(serverList, *addr)
		log.Print(addr)
		index++
	}

	node.info.clusterAddresses = serverList
	node.info.clusterCount = len(serverList)

	return node.info.clusterAddresses
}

func (node *Node) Call(address string, callable func()) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	node.grpcClient = comm.NewCommServiceClient(conn)
	callable()

	conn.Close()
}

func (node *Node) CheckSanity() {
	key := "__CheckSanity__"

	oldVal := node.app.Get(key)
	log.Printf("Checking own address of %v", node.address.String())
	node.Call(node.address.String(), func() {
		responseTest, errTest := node.grpcClient.Ping(context.Background(), &comm.BasicRequest{})
		if errTest != nil {
			log.Printf("Sanity check error: %v", errTest)
		}
		log.Printf("Test Response: %v", responseTest.Message)

		responseSet, errSet := node.grpcClient.SetValue(context.Background(), &comm.SetValueRequest{Key: key, Value: "OK"})
		if errSet != nil {
			log.Printf("Sanity check error: %v", errSet)
		}
		log.Printf("Set Response: %v", responseSet.Message)

		responseGet, errGet := node.grpcClient.GetValue(context.Background(), &comm.GetValueRequest{Key: key})
		if errGet != nil {
			log.Printf("Sanity check error: %v", errGet)
		}
		log.Printf("Get Response: %v", responseGet.Value)
	})

	node.app.Set(key, oldVal)
}

func (node *Node) Stop() {
	node.Running = false
}

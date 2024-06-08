package node

import (
	"Node/grpc/comm"
	"bufio"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type State int

const (
	Null State = iota
	Set
	Append
	Delete
	NewOldConfig
	NewConfig
)

const (
	Leader State = iota // https://go.dev/wiki/Iota
	Candidate
	Follower
)

type Info struct {
	// Persistent
	id                  int
	leaderId            int
	currentTerm         int
	votedFor            int
	log                 []comm.Entry
	clusterAddresses    []net.TCPAddr
	newClusterAddresses []net.TCPAddr
	clusterCount        int
	newClusterCount     int
	isJointConsensus    bool
	timeoutAvgTime      int

	// Volatile
	commitIndex int
	// lastApplied int

	// Leaders
	nextIndex     []int
	matchIndex    []int
	nextIndexNew  []int
	matchIndexNew []int
	serverUp      bool
}

type Node struct {
	// General attribute
	Running bool
	mutex   Mutex

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

func (node *Node) Init(hostfile string, timeoutAvgTime int, newHostfile string, isJointConsensus bool) {
	log.Printf("Initializing node")
	node.LoadLogs()
	node.InitServer()

	node.info.clusterAddresses, node.info.clusterCount = node.ReadServerList(hostfile)
	log.Printf("Check: %v", node.info.clusterAddresses)

	if isJointConsensus {
		node.info.newClusterAddresses, node.info.newClusterCount = node.ReadServerList(newHostfile)
	} else {
		node.info.newClusterAddresses = []net.TCPAddr{}
		node.info.newClusterCount = 0
	}
	node.info.isJointConsensus = isJointConsensus

	node.info.timeoutAvgTime = timeoutAvgTime

	node.resetElectionTimer()

	// node.CheckSanity()
	log.Printf("Node initialization complete with address %v and id %v", node.address.String(), node.info.id)
	node.Running = true
	node.state = Follower

	go node.ElectionTimerHandler()

	// Start the grpc-web server
	grpcWebServer := grpcweb.WrapServer(
		node.grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	// 0.0.0.0:60000 (grpc) -> 0.0.0.0:50000 (grpc-web) (added with 10000)

	addr := strings.Split(node.address.String(), ":")
	port, _ := strconv.Atoi(addr[1])
	port -= 10000
	addr[1] = strconv.Itoa(port)

	srv := &http.Server{
		Handler: grpcWebServer,
		Addr:    strings.Join(addr, ":"),
	}
	log.Printf("grpc-web server listening at %v", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && node.Running {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
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
		if err := node.grpcServer.Serve(lis); err != nil && node.Running {
			node.Running = false
			log.Fatalf("Failed to serve %v", err)
		}
	}()
}

func (node *Node) ReadServerList(filename string) ([]net.TCPAddr, int) {
	log.Printf("Initializing server list")

	var serverList []net.TCPAddr

	file, err := os.Open(filename)

	if err != nil {
		log.Printf("Failed to load hostfile: %v", err)
		return serverList, 0
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

	return serverList, len(serverList)
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

func (node *Node) CommitLogEntries(newCommitIndex int) {

	for i := node.info.commitIndex + 1; i <= newCommitIndex; i++ {
		entry := node.info.log[i]

		switch entry.Command {
		case int32(Set):
			node.app.Set(entry.Key, entry.Value)
		case int32(Append):
			node.app.Append(entry.Key, entry.Value)
		case int32(Delete):
			node.app.Del(entry.Key)
		case int32(NewConfig):
			if node.state == Leader {
				contained := false
				for _, addr := range node.info.newClusterAddresses {
					if addr.String() == node.address.String() {
						contained = true
						break
					}
				}
				if contained {
					node.info.clusterAddresses = node.info.newClusterAddresses
					node.info.clusterCount = node.info.newClusterCount
					node.info.isJointConsensus = false
				} else {
					log.Printf("[CommitLogEntries] Node %v is not in new cluster but its a leader, precedes to kill itsel", node.address.String())
					node.Stop()
				}
			}
		case int32(NewOldConfig):
			if node.state == Leader {
				newEntry := comm.Entry{
					Term:    int32(node.info.currentTerm),
					Command: int32(NewConfig),
				}
				node.appendToLog(newEntry)
			}
		}
	}
	node.info.commitIndex = newCommitIndex
}

func (node *Node) appendToLog(newEntry comm.Entry) {
	node.mutex.Lock()
	node.info.log = append(node.info.log, newEntry)

	node.info.matchIndex[node.info.id] = len(node.info.log) - 1
	// node.info.nextIndex[node.info.id] = len(node.info.log)

	for index, nextIndex := range node.info.nextIndex {
		if nextIndex == node.info.matchIndex[index] {
			node.info.nextIndex[index]++
		}
	}

	log.Printf("[Persistence] writing log via leader request")
	node.SaveLog(newEntry)
	node.mutex.Unlock()
}

func (node *Node) Stop() {
	node.Running = false
}

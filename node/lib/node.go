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
	newId               int
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
	Running    bool
	nodeMutex  Mutex
	timerMutex Mutex

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
	node.InitServer()

	node.info.clusterAddresses, node.info.clusterCount, node.info.id = node.ReadServerList(hostfile)
	log.Printf("Check: %v", node.info.clusterAddresses)

	if isJointConsensus {
		node.info.newClusterAddresses, node.info.newClusterCount, node.info.newId = node.ReadServerList(newHostfile)
		node.info.isJointConsensus = isJointConsensus
	}

	node.LoadLogs()
	node.info.timeoutAvgTime = timeoutAvgTime
	node.info.commitIndex = -1

	node.resetElectionTimer()

	// node.CheckSanity()
	log.Printf("Node initialization complete with address %v and id %v", node.address.String(), node.info.id)
	node.Running = true

	// log.Printf("[DEBUG] State is now follower")
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
		if err := node.grpcServer.Serve(lis); err != nil && node.Running {
			node.Running = false
			log.Fatalf("Failed to serve %v", err)
		}
	}()

	// Start the grpc-web server
	log.Printf("Initializing grpc web server")

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

func (node *Node) ReadServerList(filename string) ([]net.TCPAddr, int, int) {
	log.Printf("Initializing server list")

	var serverList []net.TCPAddr
	id := -1

	file, err := os.Open(filename)

	if err != nil {
		log.Printf("Failed to load hostfile: %v", err)
		return serverList, 0, id
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
			id = index
		}

		serverList = append(serverList, *addr)
		log.Print(addr)
		index++
	}

	return serverList, len(serverList), id
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
	if newCommitIndex <= node.info.commitIndex {
		return
	}
	// log.Printf("[DEBUG] locking nodemutex")
	// node.nodeMutex.Lock()

	// log.Printf("\n\n[DEBUG] committing index %v with current commit idx is %v", newCommitIndex, node.info.commitIndex)
	// log.Printf("nextIndex: %v", node.info.nextIndex)
	// log.Printf("matchIndex: %v", node.info.matchIndex)
	// log.Printf("nextIndexNew: %v", node.info.nextIndexNew)
	// log.Printf("matchIndexNew: %v", node.info.matchIndexNew)
	// log.Printf("\n\n")
	for i := node.info.commitIndex + 1; i <= newCommitIndex; i++ {
		entry := node.info.log[i]
		log.Printf("[CommitLogEntries] Committing entry index %v with command %v", i, entry.Command)

		switch entry.Command {
		case int32(Set):
			node.app.Set(entry.Key, entry.Value)
		case int32(Append):
			node.app.Append(entry.Key, entry.Value)
		case int32(Delete):
			node.app.Del(entry.Key)
		case int32(NewConfig):
			contained := false
			for _, addr := range node.info.newClusterAddresses {
				if addr.String() == node.address.String() {
					contained = true
					break
				}
			}
			if contained {
				if node.state == Leader {
					log.Printf("[CommitLogEntries] Node Leader %v change to NewConfig", node.address.String())
					node.info.clusterAddresses = node.info.newClusterAddresses
					node.info.clusterCount = node.info.newClusterCount
					node.info.id = node.info.newId
					node.info.isJointConsensus = false

					// TODO: Stay leader or reset?
					node.info.leaderId = node.info.newId

					lastLogIdx := int32(len(node.info.log) - 1)
					node.info.matchIndex = make([]int, node.info.clusterCount)
					node.info.nextIndex = make([]int, node.info.clusterCount)

					for i := range node.info.nextIndex {
						node.info.nextIndex[i] = int(lastLogIdx) + 1
						node.info.matchIndex[i] = -1
					}
					node.info.matchIndex[node.info.id] = int(lastLogIdx)
					node.info.nextIndex[node.info.id] = int(lastLogIdx) + 1
				} else {
					log.Printf("[CommitLogEntries] Node Follower %v change to NewConfig", node.address.String())
					node.info.clusterAddresses = node.info.newClusterAddresses
					node.info.clusterCount = node.info.newClusterCount
					node.info.id = node.info.newId
					node.info.isJointConsensus = false
				}
			} else {
				if node.state == Leader {
					if entry.Term == int32(node.info.currentTerm) {
						log.Printf("[CommitLogEntries] Leader %v is not in new cluster, procedes to kill itself", node.address.String())
						node.Stop()
					} else {
						log.Printf("[CommitLogEntries] Leader %v is not in new cluster but its old term, ignore", node.address.String())
					}
				}
			}
		case int32(NewOldConfig):
			if node.state == Leader {
				newEntry := comm.Entry{
					Term:    int32(node.info.currentTerm),
					Command: int32(NewConfig),
					Key:     "-",
					Value:   "-",
				}
				node.appendToLog(newEntry)
			}
		}
	}
	node.info.commitIndex = newCommitIndex

	// log.Printf("[DEBUG] unlocking nodemutex")
	// node.nodeMutex.Unlock()
}

func (node *Node) appendToLog(newEntry comm.Entry) {
	// log.Printf("[DEBUG] locking nodemutex")
	// node.nodeMutex.Lock()
	node.info.log = append(node.info.log, newEntry)

	if node.info.id != -1 {
		node.info.matchIndex[node.info.id] = len(node.info.log) - 1
		node.info.nextIndex[node.info.id] = len(node.info.log)
	}

	// for index, nextIndex := range node.info.nextIndex {
	// 	if nextIndex == node.info.matchIndex[index] {
	// 		node.info.nextIndex[index]++
	// 	}
	// }

	if node.info.isJointConsensus {
		if node.info.newId != -1 {
			node.info.matchIndexNew[node.info.newId] = len(node.info.log) - 1
			node.info.nextIndexNew[node.info.newId] = len(node.info.log)
		}

		// 	for index, nextIndex := range node.info.nextIndexNew {
		// 		if nextIndex == node.info.matchIndexNew[index] {
		// 			node.info.nextIndexNew[index]++
		// 		}
		// 	}
	}

	log.Printf("[Persistence] writing log via request")
	node.SaveLog(newEntry)

	// log.Printf("[DEBUG] locking nodemutex")
	// node.nodeMutex.Unlock()
}

func (node *Node) Stop() {
	node.Running = false
}

func (node *Node) Status() {
	log.Printf("\n\n---------------Node Info:")
	log.Printf("id: %v", node.info.id)
	log.Printf("state: %v", node.state)
	log.Printf("leaderId: %v", node.info.leaderId)
	log.Printf("currentTerm: %v", node.info.currentTerm)
	log.Printf("commitIndex %v", node.info.commitIndex)
	// log.Printf("lastApplied %v", node.info.lastApplied)
	log.Printf("clusterAddresses: %v", node.info.clusterAddresses)
	log.Printf("clusterCount: %v", node.info.clusterCount)
	log.Printf("timeoutAvgTime: %v", node.info.timeoutAvgTime)
	log.Printf("isJointConsensus: %v", node.info.isJointConsensus)

	log.Printf("/\n\nlog info:")
	for index, entry := range node.info.log {
		log.Printf("[%v]: command: %v, key:%v, value:%v, term:%v", index, entry.Command, entry.Key, entry.Value, entry.Term)
	}

	if node.info.isJointConsensus {
		log.Printf("/\n\nJoint Consensus info:")

		log.Printf("newId: %v", node.info.newId)
		log.Printf("newClusterAddresses: %v", node.info.newClusterAddresses)
		log.Printf("newClusterCount: %v", node.info.newClusterCount)
	}

	if node.state == Leader {
		log.Printf("/\n\nLeader info:")

		log.Printf("serverUp: %v", node.info.serverUp)
		log.Printf("nextIndex: %v", node.info.nextIndex)
		log.Printf("matchIndex: %v", node.info.matchIndex)

		if node.info.isJointConsensus {
			log.Printf("nextIndexNew: %v", node.info.nextIndexNew)
			log.Printf("matchIndexNew: %v", node.info.matchIndexNew)
		}
	}
}

// func (node *Node) PrepNewConfig(info string) {
// 	node.info.isJointConsensus = true
// 	node.info.newClusterAddresses = parseAddresses(info)
// 	node.info.newClusterCount = len(node.info.newClusterAddresses)
// 	node.info.idNew = -1
// 	for index, address := range node.info.newClusterAddresses {
// 		if node.address.String() == address.String() {
// 			node.info.idNew = index
// 			break
// 		}
// 	}
// }

// func (node *Node) ApplyNewConfig() {
// 	if node.info.idNew == -1 {
// 		log.Printf("[Config] Node %v is not in new cluster, precedes to kill itself", node.address.String())
// 		node.Stop()
// 	}

// 	node.info.clusterAddresses = node.info.newClusterAddresses
// 	node.info.clusterCount = node.info.newClusterCount
// 	node.info.id = node.info.idNew
// 	node.info.isJointConsensus = false

// 	// TODO: Reset leadership?
// 	// node.info.leaderId = node.info.idNew
// }

package node

import (
	"Node/grpc/comm"
	"log"
	"math/rand"
	"time"
)

func (node *Node) GenerateTimeoutDelay() time.Duration {
	minTimeout := node.info.timeoutAvgTime / 2
	maxTimeout := node.info.timeoutAvgTime + minTimeout

	// log.Printf("Node avg timeout duration is: %v", node.info.timeoutAvgTime)
	// log.Printf("Generated min duration is: %v", minTimeout)
	// log.Printf("Generated max duration is: %v", maxTimeout)

	timeoutTime := rand.Int63n(int64(maxTimeout-minTimeout)) + int64(minTimeout)
	timeout := time.Duration(timeoutTime) * time.Millisecond
	log.Printf("Generated timeout duration is: %v", timeout)
	return timeout
}

func (node *Node) resetElectionTimer() {
	if node.electionTimer != nil {
		node.electionTimer.Stop()
	}

	node.electionTimer = time.NewTimer(node.GenerateTimeoutDelay())
}

func (node *Node) ElectionTimerHandler() {
	for node.Running {
		select {
		case <-node.electionTimer.C:
			node.mutex.Lock()
			log.Println("Election timeout occured")
			node.resetElectionTimer()
			node.mutex.Unlock()
		case <-node.electionResetSignal:
			node.mutex.Lock()
			node.resetElectionTimer()
			node.mutex.Unlock()
		}
	}
}

func (node *Node) onHeartBeat(info *comm.HeartbeatRequest) {
	log.Printf("Received heartbeat from: %v", info.LeaderId)
	node.electionResetSignal <- true
}

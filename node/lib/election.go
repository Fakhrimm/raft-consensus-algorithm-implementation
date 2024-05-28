package node

import (
	"Node/grpc/comm"
	"context"
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
	// timeout := time.Duration(timeoutTime) * time.Millisecond
	timeout := time.Duration(timeoutTime) * time.Second / 100 // Second to visualize it better, switch to millis on actual

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
	for node.Running && node.state == Follower {
		select {
		case <-node.electionTimer.C:
			node.mutex.Lock()
			log.Println("Election timeout occured")
			node.startElection()
			node.resetElectionTimer()
			node.mutex.Unlock()
		case <-node.electionResetSignal:
			node.mutex.Lock()
			node.resetElectionTimer()
			node.mutex.Unlock()
		}
	}
}

func (node *Node) startElection() {
	node.info.currentTerm++
	node.state = Candidate
	vote := 1

	log.Printf("[Election] Starting election for cluster size of %v", node.info.clusterCount)
	for _, peer := range node.info.clusterAddresses {
		if peer.String() == node.address.String() {
			continue
		}

		log.Printf("[Election] node is : %v", node.address.String())
		log.Printf("[Election] Asking for vote from node: %v", peer)

		// TODO: Implement
		node.Call(peer.String(), func() {

			// TODO: review, last log index committed or uncommited? not sure
			lastLogTerm := 0
			if len(node.info.log) > 0 {
				lastLogTerm = node.info.log[node.info.commitIndex].term
			}
			response, err := node.grpcClient.RequestVote(context.Background(), &comm.RequestVoteRequest{
				Term:         int32(node.info.currentTerm),
				CandidateId:  int32(node.info.id),
				LastLogIndex: int32(node.info.commitIndex),
				LastLogTerm:  int32(lastLogTerm),
			})

			if err != nil {
				log.Printf("[Election] Error reaching node %v: %v", peer.String(), err)
			} else if response.VoteGranted {
				log.Printf("[Election] Received vote from node %v, current vote is %v/%v", peer.String(), vote, node.info.clusterCount)
				vote += 1
			}
		})

		if vote > node.info.clusterCount/2 {
			log.Printf("[Election] Enough vote received")
			break
		}
	}

	if vote > node.info.clusterCount/2 {
		log.Printf("[Election] Node is now a leader for term %v", node.info.currentTerm)
		node.state = Leader
		node.initHeartbeat()
	} else {
		node.state = Follower
		log.Print("[Election] Not enough vote!")
	}
}

func (node *Node) initHeartbeat() {
	interval := time.Duration(node.info.timeoutAvgTime) * time.Second / 100 / 6
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for node.state == Leader {
		for range ticker.C {
			for _, peer := range node.info.clusterAddresses {
				if peer.String() == node.address.String() {
					continue
				}
				log.Printf("Sending heartbeat to node: %v", peer)

				// TODO: Implement
				node.Call(peer.String(), func() {
					node.grpcClient.Heartbeat(context.Background(), &comm.HeartbeatRequest{})
				})
			}
		}
	}
}

func (node *Node) onHeartBeat(info *comm.HeartbeatRequest) {
	// TODO: Handle if a candidate node receives a heartbeat
	// TODO: Implement
	log.Printf("Received heartbeat from: %v", info.LeaderId)
	node.electionResetSignal <- true
}

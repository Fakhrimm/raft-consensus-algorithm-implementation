package node

import (
	"Node/grpc/comm"
	"context"
	"log"
	"math/rand"
	"sync"
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
	interval := time.Duration(node.info.timeoutAvgTime) * time.Second / 100 / 3

	node.info.currentTerm++
	node.state = Candidate
	vote := 1

	log.Printf("[Election] Starting election for cluster size of %v", node.info.clusterCount)

	// Create channels
	var waitGroup sync.WaitGroup
	voteCh := make(chan bool, node.info.clusterCount)

	for _, peer := range node.info.clusterAddresses {
		if peer.String() == node.address.String() {
			continue
		}

		log.Printf("[Election] node is : %v", node.address.String())
		log.Printf("[Election] Asking for vote from node: %v", peer)

		waitGroup.Add(1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*interval)

		// TODO: Implement vote request here
		go func(peerAddr string) {
			defer cancel()
			defer waitGroup.Done()

			node.Call(peer.String(), func() {

				// TODO: review, last log index committed or uncommited? not sure
				lastLogTerm := 0
				if len(node.info.log) > 0 {
					lastLogTerm = node.info.log[node.info.commitIndex].term
				}
				response, err := node.grpcClient.RequestVote(ctx, &comm.RequestVoteRequest{
					Term:         int32(node.info.currentTerm),
					CandidateId:  int32(node.info.id),
					LastLogIndex: int32(node.info.commitIndex),
					LastLogTerm:  int32(lastLogTerm),
				})

				voteCh <- false
				if err != nil {
					log.Printf("[Election] Error reaching node %v: %v", peer.String(), err)
				} else if response.VoteGranted {
					log.Printf("[Election] Received vote from node %v", peer.String())
					voteCh <- true
				}
			})
		}(peer.String())
	}

	// Create barrier, reduce votes
	go func() {
		waitGroup.Wait()
		close(voteCh)
	}()

	for granted := range voteCh {
		if granted {
			vote += 1
			log.Printf("[Election] Counting votes %v/%v", vote, node.info.clusterCount)
		}
	}
	log.Printf("[Election] Counting completed with %v/%v", vote, node.info.clusterCount)

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
				log.Printf("[Heartbeat] Sending heartbeat to node: %v", peer)

				ctx, cancel := context.WithTimeout(context.Background(), 2*interval)

				go func(peerAddr string) {
					defer cancel()

					// TODO: Implement heartbeat functions
					node.Call(peer.String(), func() {
						_, err := node.grpcClient.Heartbeat(ctx, &comm.HeartbeatRequest{})
						if err != nil {
							log.Printf("[Heartbeat] failed to send heartbeat to %v: %v", peer.String(), err)
						}
					})
				}(peer.String())
			}
		}
	}
}

func (node *Node) onHeartBeat(info *comm.HeartbeatRequest) {
	// TODO: Handle if a candidate node receives a heartbeat
	// TODO: Implement
	log.Printf("[Heartbeat] Received heartbeat from: %v", info.LeaderId)
	node.electionResetSignal <- true
}

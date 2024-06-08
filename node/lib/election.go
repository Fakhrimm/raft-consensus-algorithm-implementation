package node

import (
	"Node/grpc/comm"
	"context"
	"log"
	"math/rand"
	"net"
	"sort"
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

	// log.Printf("[Election] Generated timeout duration is: %v", timeout)
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
			log.Println("[Election] Election timeout occured")
			node.startElection()
			node.resetElectionTimer()
			node.mutex.Unlock()
		case <-node.electionResetSignal:
			// log.Println("[Election] Received election reset signal")
			node.mutex.Lock()
			node.resetElectionTimer()
			node.mutex.Unlock()
		}
	}
}

func (node *Node) startElection() {
	interval := time.Duration(node.info.timeoutAvgTime) * time.Second / 100 / 3
	log.Printf("[Election Start] Current TERM BEFORE: %v", node.info.currentTerm)
	node.info.currentTerm++
	log.Printf("[Election Start] START ELECTION; TERM: %v; CLUSTER SIZE: %v", node.info.currentTerm, node.info.clusterCount)
	node.state = Candidate
	lastLogIdx := int32(len(node.info.log) - 1)

	// TODO: Test joint consensus & improve concurrency
	vote := node.sendElection(node.info.clusterAddresses, interval)
	if node.info.isJointConsensus {
		voteNew := node.sendElection(node.info.newClusterAddresses, interval)

		if (vote > node.info.clusterCount/2) && (voteNew > node.info.newClusterCount/2) {
			log.Printf("[Election] ELECTED; TERM: %v; isJointConsensus: TRUE", node.info.currentTerm)
			node.info.serverUp = true
			node.state = Leader

			// Old
			node.info.matchIndex = make([]int, node.info.clusterCount)
			node.info.matchIndex[node.info.id] = len(node.info.log) - 1

			node.info.nextIndex = make([]int, node.info.clusterCount)
			for i := range node.info.nextIndex {
				node.info.nextIndex[i] = int(lastLogIdx)
				node.info.matchIndex[i] = -1
			}

			// New
			node.info.matchIndexNew = make([]int, node.info.newClusterCount)
			node.info.matchIndexNew[node.info.id] = len(node.info.log) - 1

			node.info.nextIndexNew = make([]int, node.info.newClusterCount)
			for i := range node.info.nextIndexNew {
				node.info.nextIndexNew[i] = int(lastLogIdx)
				node.info.matchIndexNew[i] = -1
			}
			// log.Printf("[Election] matchIndex: %v", node.info.matchIndex)
			// log.Printf("[Election] nextIndex: %v", node.info.nextIndex)

			node.startAppendEntries()
		} else {
			node.state = Follower
			log.Print("[Election] Not enough vote!")
		}
	} else {
		if vote > node.info.clusterCount/2 {
			log.Printf("[Election] ELECTED; TERM: %v; isJointConsensus: FALSE", node.info.currentTerm)
			node.info.serverUp = true
			node.state = Leader

			node.info.matchIndex = make([]int, node.info.clusterCount)
			node.info.matchIndex[node.info.id] = len(node.info.log) - 1

			node.info.nextIndex = make([]int, node.info.clusterCount)
			for i := range node.info.nextIndex {
				node.info.nextIndex[i] = int(lastLogIdx)
				node.info.matchIndex[i] = -1
			}
			// log.Printf("[Election] matchIndex: %v", node.info.matchIndex)
			// log.Printf("[Election] nextIndex: %v", node.info.nextIndex)

			node.startAppendEntries()
		} else {
			node.state = Follower
			log.Print("[Election] Not enough vote!")
		}
	}

}

func (node *Node) sendElection(address []net.TCPAddr, interval time.Duration) int {
	vote := 1
	length := len(address)

	// Fetch data for voting
	lastLogIdx := int32(len(node.info.log) - 1)
	lastLogTerm := int32(0)
	if lastLogIdx > -1 {
		lastLogTerm = node.info.log[lastLogIdx].Term
	}

	// Create channels
	var waitGroup sync.WaitGroup
	voteCh := make(chan bool, length)

	for _, peer := range address {
		if peer.String() == node.address.String() {
			continue
		}

		//log.Printf("[Election] node is : %v", node.address.String())
		//log.Printf("[Election] Asking for vote from node: %v", peer)

		waitGroup.Add(1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*interval)

		go func(peerAddr string) {
			defer cancel()
			defer waitGroup.Done()

			node.Call(peer.String(), func() {
				response, err := node.grpcClient.RequestVote(ctx, &comm.RequestVoteRequest{
					Term:         int32(node.info.currentTerm),
					CandidateId:  int32(node.info.id),
					LastLogIndex: lastLogIdx,
					LastLogTerm:  lastLogTerm,
				})

				voteCh <- false
				if err != nil {
					log.Printf("[Election] Error reaching node %v: %v", peer.String(), err)
				} else if response.VoteGranted {
					log.Printf("[Election] Received vote from node %v", peer.String())
					voteCh <- true
				} else if response.Term > int32(node.info.currentTerm) {
					log.Printf("[Election] Received higher term %v from node, turns to follower %v", peer.String())
					log.Printf("[Election] [UPDATE TERM] Before TERM: %v, After TERM: %v", node.info.currentTerm, response.Term)
					node.info.currentTerm = int(response.Term)
					node.state = Follower
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
			log.Printf("[Election] Counting votes %v/%v", vote, length)
		}
	}
	log.Printf("[Election] Counting completed with %v/%v", vote, length)

	return vote
}

func (node *Node) startAppendEntries() {
	// Unlock mutex from election
	node.mutex.Unlock()

	interval := time.Duration(node.info.timeoutAvgTime) * time.Second / 100 / 6
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for node.state == Leader {
		for range ticker.C {
			aliveNode := 0
			aliveNode = node.sendAppendEntries(node.info.clusterAddresses, interval)
			if node.state != Leader {
				return
			}

			// TODO: Test joint consensus & improve by concurrency
			if node.info.isJointConsensus {
				newAliveNode := 0
				newAliveNode = node.sendAppendEntries(node.info.newClusterAddresses, interval)
				log.Printf("[Heartbeat]")

				if (2*aliveNode <= node.info.clusterCount) && (2*newAliveNode <= node.info.newClusterCount) {
					node.info.serverUp = false
				} else {
					node.updateMajority()
					node.info.serverUp = true
				}
			} else {
				// log.Printf("[Heartbeat] Total received heartbeat is %v/%v", aliveNodeCount, node.info.clusterCount)
				if 2*aliveNode <= node.info.clusterCount {
					node.info.serverUp = false
				} else {
					node.updateMajority()
					node.info.serverUp = true
				}
			}
		}
	}
}

func (node *Node) sendAppendEntries(address []net.TCPAddr, interval time.Duration) int {
	var waitGroup sync.WaitGroup
	heartbeatCh := make(chan bool, node.info.clusterCount)

	for index, peer := range node.info.clusterAddresses {
		if peer.String() == node.address.String() {
			continue
		}
		// log.Printf("[Heartbeat] Sending heartbeat to node %v: %v", index, peer)

		prevLogIndex := max(node.info.nextIndex[index]-1, -1)
		prevLogTerm := int32(0)
		if prevLogIndex > -1 {
			prevLogTerm = node.info.log[prevLogIndex].Term
		}

		// TODO: Review, I don't think this is efficient
		length := len(node.info.log) - prevLogIndex - 1

		sentEntries := make([]*comm.Entry, length)
		for i := 0; i < length; i++ {
			sentEntries[i] = &comm.Entry{
				Term:    node.info.log[prevLogIndex+1+i].Term,
				Key:     node.info.log[prevLogIndex+1+i].Key,
				Value:   node.info.log[prevLogIndex+1+i].Value,
				Command: node.info.log[prevLogIndex+1+i].Command,
			}
		}

		data := &comm.AppendEntriesRequest{
			Term:         int32(node.info.currentTerm),
			LeaderId:     int32(node.info.id),
			PrevLogTerm:  prevLogTerm,
			PrevLogIndex: int32(prevLogIndex),
			LeaderCommit: int32(node.info.commitIndex),
			Entries:      sentEntries,
		}

		// log.Printf("[Heartbeat] data to send: %v", data)
		// log.Printf("[Heartbeat] data term: %v", data.Term)
		// log.Printf("[Heartbeat] data leader id: %v", data.LeaderId)
		// log.Printf("[Heartbeat] data entries: %v", data.Entries)

		waitGroup.Add(1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*interval)

		go func(peerAddr string) {
			defer cancel()
			defer waitGroup.Done()

			node.Call(peer.String(), func() {
				response, err := node.grpcClient.AppendEntries(ctx, data)
				if err != nil {
					log.Printf("[Heartbeat] failed to send heartbeat to %v: %v", peer.String(), err)
					log.Printf("[Heartbeat] Node is a: %v", node.state)
					heartbeatCh <- false
				} else {
					if response.Term > int32(node.info.currentTerm) {
						log.Printf("[Heartbeat] Received higher term %v from node %v, turns to follower", response.Term, peer.String())
						log.Printf("[Heartbeat] [UPDATE TERM] Before TERM: %v, After TERM: %v", node.info.currentTerm, response.Term)
						node.info.currentTerm = int(response.Term)
						node.state = Follower
						return
					} else {
						if response.Success {
							node.info.matchIndex[index] = int(prevLogIndex) + len(sentEntries)
							node.info.nextIndex[index] = node.info.matchIndex[index] + 1
						} else {
							node.info.nextIndex[index]--
						}
					}
					heartbeatCh <- true
				}
			})
		}(peer.String())
	}
	go func() {
		waitGroup.Wait()
		close(heartbeatCh)
	}()

	aliveNodeCount := 1
	for getResponse := range heartbeatCh {
		if getResponse {
			aliveNodeCount += 1
		}
	}

	return aliveNodeCount
}

func (node *Node) updateMajority() {
	majority := node.getMajority(node.info.matchIndex, node.info.clusterCount)

	// TODO: Test joint consensus & improve concurrency
	if node.info.isJointConsensus {
		newMajority := node.getMajority(node.info.matchIndexNew, node.info.newClusterCount)
		node.CommitLogEntries(min(majority, newMajority))
	} else {
		node.CommitLogEntries(majority)
	}
}

func (node *Node) getMajority(matchArray []int, clusterCount int) int {
	matchIndex := make(map[int]int)
	for _, index := range matchArray {
		if matchIndex[index] == 0 {
			matchIndex[index] = 1
		} else {
			matchIndex[index]++
		}
	}

	// sort matchIndex descending by the key
	keys := make([]int, 0, len(matchIndex))
	for k := range matchIndex {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	majority := clusterCount / 2
	sum := 0

	for _, key := range keys {
		sum += matchIndex[key]
		if sum > majority {
			// log.Printf("[Transaction] updating majority keys: %v", keys)
			return key
		}
	}
	return 0
}

func (node *Node) onHeartBeat() {
	node.electionResetSignal <- true
}

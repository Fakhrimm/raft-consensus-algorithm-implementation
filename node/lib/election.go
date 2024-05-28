package node

import (
	"Node/grpc/comm"
	"context"
	"log"
	"math/rand"
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

	node.info.currentTerm++
	node.state = Candidate
	vote := 1

	log.Printf("[Election] Starting election for cluster size of %v", node.info.clusterCount)

	// Fetch data for voting
	lastLogIdx := int32(len(node.info.log) - 1)
	lastLogTerm := int32(0)
	if lastLogIdx > -1 {
		lastLogTerm = node.info.log[lastLogIdx].Term
	}

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
		node.info.serverUp = true
		node.state = Leader

		node.info.matchIndex = make([]int, node.info.clusterCount)
		node.info.nextIndex = make([]int, node.info.clusterCount)
		for i := range node.info.nextIndex {
			node.info.nextIndex[i] = int(lastLogIdx)
		}
		log.Printf("[Election] matchIndex: %v", node.info.matchIndex)
		log.Printf("[Election] nextIndex: %v", node.info.nextIndex)

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

			var waitGroup sync.WaitGroup
			heartbeatCh := make(chan bool, node.info.clusterCount)

			for index, peer := range node.info.clusterAddresses {
				if peer.String() == node.address.String() {
					continue
				}
				// log.Printf("[Heartbeat] Sending heartbeat to node %v: %v", index, peer)

				prevLogIndex := node.info.nextIndex[index] - 1
				prevLogTerm := int32(0)
				if prevLogIndex > -1 {
					prevLogTerm = node.info.log[prevLogIndex].Term
				}

				// TODO: Review, I don't think this is efficient
				length := len(node.info.log) - prevLogIndex - 2
				sentEntries := make([]*comm.Entry, length)
				for i := 0; i < length; i++ {
					sentEntries[i] = &comm.Entry{
						Term:          node.info.log[prevLogIndex+1+i].Term,
						Key:           node.info.log[prevLogIndex+1+i].Key,
						Value:         node.info.log[prevLogIndex+1+i].Value,
						IsConfigEntry: false,
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

				// TODO: Count failed heartbeats to indicate majority down
				go func(peerAddr string) {
					defer cancel()
					defer waitGroup.Done()

					// TODO: Implement append entries functions
					node.Call(peer.String(), func() {
						response, err := node.grpcClient.AppendEntries(ctx, data)
						if err != nil {
							log.Printf("[Heartbeat] failed to send heartbeat to %v: %v", peer.String(), err)
							heartbeatCh <- false
						} else {
							if response.Success {
								node.info.matchIndex[index] = int(prevLogIndex) + len(sentEntries)
								node.info.nextIndex[index] = node.info.matchIndex[index] + 1
							} else {
								node.info.nextIndex[index]--
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
			// log.Printf("[Heartbeat] Total received heartbeat is %v/%v", aliveNodeCount, node.info.clusterCount)

			if 2*aliveNodeCount <= node.info.clusterCount {
				// TODO: implement if mayority of server is not alive
				node.info.serverUp = false
			} else {
				node.updateMajority()
				node.info.serverUp = true
			}
		}
	}
}

func (node *Node) updateMajority() {
	matchIndex := make(map[int]int)
	for _, index := range node.info.matchIndex {
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

	majority := node.info.clusterCount / 2
	sum := 0
	for _, key := range keys {
		sum += matchIndex[key]
		if sum > majority {
			node.CommitLogEntries(key)
			break
		}
	}
}

func (node *Node) onHeartBeat() {
	node.electionResetSignal <- true
}

package node

import (
	"Node/grpc/comm"
	"context"
	"log"
	"net"
)

type server struct {
	comm.UnimplementedCommServiceServer
	Node        *Node
	serverMutex Mutex
}

// Application purposes
func (s *server) ValidateRequest() (code int32, message string) {
	if s.Node.state != Leader {
		log.Printf("[Transaction] Transaction refused, node is not leader")
		return 501, s.Node.info.clusterAddresses[s.Node.info.leaderId].String()
	} else if !s.Node.info.serverUp {
		log.Printf("[Transaction] Transaction refused, not enough active nodes")
		return 503, "Service Unavailable"
	} else {
		return 200, "OK"
	}
}

func (s *server) Ping(ctx context.Context, in *comm.BasicRequest) (*comm.BasicResponse, error) {
	log.Printf("[System] Received ping")
	return &comm.BasicResponse{Code: 200, Message: "PONG"}, nil
}

func (s *server) Status(ctx context.Context, in *comm.BasicRequest) (*comm.BasicResponse, error) {
	log.Printf("[System] Received status request")
	s.Node.Status()

	// s.Node.SaveLogs()

	return &comm.BasicResponse{Code: 200, Message: "STAT"}, nil
}

func (s *server) Stop(ctx context.Context, in *comm.BasicRequest) (*comm.BasicResponse, error) {
	defer s.Node.Stop()
	log.Printf("[Config] Received stop signal")
	return &comm.BasicResponse{Code: 200, Message: "STOPPING"}, nil
}

func (s *server) GetValue(ctx context.Context, in *comm.GetValueRequest) (*comm.GetValueResponse, error) {
	var code int32
	var message string
	var value string

	// TODO: Give feedback when value is not set
	code, message = s.ValidateRequest()
	if code == 200 {
		message = "Value fetched successfully"
		value = s.Node.app.Get(in.Key)
	}

	return &comm.GetValueResponse{Code: code, Message: message, Value: value}, nil
}

func (s *server) SetValue(ctx context.Context, in *comm.SetValueRequest) (*comm.SetValueResponse, error) {
	log.Printf("[Transaction] Recevied set request of %v:%v", in.Key, in.Value)
	var code int32
	var message string
	var value string

	code, message = s.ValidateRequest()
	if code == 200 {
		message = "Value set request completed"

		s.Node.appendToLog(comm.Entry{
			Term:    int32(s.Node.info.currentTerm),
			Key:     in.Key,
			Value:   in.Value,
			Command: int32(Set),
		})

		transactionIndex := len(s.Node.info.log) - 1
		for transactionIndex != s.Node.info.commitIndex {
			// log.Printf("[Transaction] Waiting to sync transactionIndex: %v, commitIndex: %v", transactionIndex, s.Node.info.commitIndex)
		}
	}

	log.Printf("[Transaction] Set request is completed")
	return &comm.SetValueResponse{Code: code, Message: message, Value: value}, nil
}

func (s *server) StrlnValue(ctx context.Context, in *comm.StrlnValueRequest) (*comm.StrlnValueResponse, error) {
	log.Printf("[Transaction] Recevied strlen request of %v", in.Key)
	var code int32
	var message string
	var value int32

	code, message = s.ValidateRequest()
	if code == 200 {
		message = "Strlen fetched successfully"
		value = int32(s.Node.app.Strln(in.Key))
	}

	return &comm.StrlnValueResponse{Code: code, Message: message, Value: value}, nil
}

func (s *server) DeleteValue(ctx context.Context, in *comm.DeleteValueRequest) (*comm.DeleteValueResponse, error) {
	log.Printf("[Transaction] Recevied delete request of %v", in.Key)
	var code int32
	var message string
	var value string

	code, message = s.ValidateRequest()
	if code == 200 {
		message = "Value delete request completed"

		s.Node.appendToLog(comm.Entry{
			Term:    int32(s.Node.info.currentTerm),
			Key:     in.Key,
			Value:   "-",
			Command: int32(Delete),
		})

		transactionIndex := len(s.Node.info.log) - 1
		for transactionIndex != s.Node.info.commitIndex {
			// log.Printf("[Transaction] Waiting to sync transactionIndex: %v, commitIndex: %v", transactionIndex, s.Node.info.commitIndex)
		}
	}

	log.Printf("[Transaction] Delete request is completed")
	return &comm.DeleteValueResponse{Code: code, Message: message, Value: value}, nil
}

func (s *server) AppendValue(ctx context.Context, in *comm.AppendValueRequest) (*comm.AppendValueResponse, error) {
	log.Printf("[Transaction] Recevied append request of %v:%v", in.Key, in.Value)
	var code int32
	var message string
	var value string

	code, message = s.ValidateRequest()
	if code == 200 {
		message = "Value append request completed"

		s.Node.appendToLog(comm.Entry{
			Term:    int32(s.Node.info.currentTerm),
			Key:     in.Key,
			Value:   in.Value,
			Command: int32(Append),
		})

		transactionIndex := len(s.Node.info.log) - 1
		for transactionIndex != s.Node.info.commitIndex {
			// log.Printf("[Transaction] Waiting to sync transactionIndex: %v, commitIndex: %v", transactionIndex, s.Node.info.commitIndex)
		}
	}

	log.Printf("[Transaction] Append request is completed")
	return &comm.AppendValueResponse{Code: code, Message: message, Value: value}, nil
}

func (s *server) ChangeMembership(ctx context.Context, in *comm.ChangeMembershipRequest) (*comm.ChangeMembershipResponse, error) {
	log.Printf("[Config] Received membership change request")

	// clusterAddresses & newClusterAddresses example:
	// "10.1.78.242:60000,10.1.78.242:60001,10.1.78.242:60002"
	// (separated by comma)
	var code int32
	var message string

	var newClusterAddresses []net.TCPAddr

	// Validate the request first
	code, message = s.ValidateRequest()
	if code == 200 {
		// Parse the (old) addresses and new addresses
		newClusterAddresses = parseAddresses(in.NewClusterAddresses)
		log.Printf("[Config] New addresses: %v", newClusterAddresses)

		// Set it in the node
		s.Node.info.newClusterAddresses = newClusterAddresses
		s.Node.info.newClusterCount = len(newClusterAddresses)
		s.Node.info.newId = -1
		for index, address := range newClusterAddresses {
			if s.Node.address.String() == address.String() {
				s.Node.info.newId = index
				break
			}
		}

		// Set that the node is in the joint consensus state
		s.Node.info.isJointConsensus = true

		// Initialize arrays
		lastLogIdx := int32(len(s.Node.info.log) - 1)
		s.Node.info.matchIndexNew = make([]int, s.Node.info.newClusterCount)
		s.Node.info.nextIndexNew = make([]int, s.Node.info.newClusterCount)

		for i := range s.Node.info.nextIndexNew {
			s.Node.info.nextIndexNew[i] = int(lastLogIdx) + 1
			s.Node.info.matchIndexNew[i] = -1
		}

		// Append entry to log
		s.Node.appendToLog(comm.Entry{
			Term:    int32(s.Node.info.currentTerm),
			Key:     "-",
			Value:   in.NewClusterAddresses,
			Command: int32(NewOldConfig),
		})
	}

	return &comm.ChangeMembershipResponse{Code: code, Message: message, Value: in.NewClusterAddresses}, nil
}

// Raft purposes
func (s *server) AppendEntries(ctx context.Context, in *comm.AppendEntriesRequest) (*comm.AppendEntriesResponse, error) {
	s.Node.nodeMutex.Lock()
	defer s.Node.nodeMutex.Unlock("[AppendEntries]")

	log.Printf("[AppendEntries] Received append entries request from %v", in.LeaderId)

	//Reply false if term from leader < currentTerm (§5.1)
	if in.Term < int32(s.Node.info.currentTerm) {
		return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: false}, nil
	} else if in.Term > int32(s.Node.info.currentTerm) {
		// Update currentTerm if leaderTerm > currentTerm
		log.Printf("[AppendEntries] [UPDATE TERM] Before: %v, After: %v", s.Node.info.currentTerm, in.Term)
		s.Node.info.currentTerm = int(in.Term)
	}

	// Reset election timer
	s.Node.onHeartBeat()
	s.Node.info.leaderId = int(in.LeaderId)

	// Reply false if log does not contain an entry at prevLogIndex
	if len(s.Node.info.log) <= int(in.PrevLogIndex) {
		log.Printf("[AppendEntries] No entry, len is %v against index of %v", len(s.Node.info.log), in.PrevLogIndex)
		return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: false}, nil
	}

	if in.PrevLogIndex >= 0 {
		// Reply false if log entry term at prevLogIndex does not match prevLogTerm (§5.3)
		if s.Node.info.log[in.PrevLogIndex].Term != in.PrevLogTerm {
			log.Printf("[AppendEntries] no match")
			return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: false}, nil
		}
	}

	overwrite := false
	if len(s.Node.info.log) > 0 {
		maxIdx := max(0, in.PrevLogIndex+1)
		overwrite = maxIdx < int32(len(s.Node.info.log))

		// log.Printf("[DEBUG] taking index of %v from %v", in.PrevLogIndex+1, len(s.Node.info.log))

		s.Node.info.log = s.Node.info.log[:maxIdx]
	}

	newConfigIdx := -1
	oldnewConfigIdx := -1
	for index, entry := range in.Entries {
		newLog := comm.Entry{
			Term:    entry.Term,
			Key:     entry.Key,
			Value:   entry.Value,
			Command: entry.Command,
		}
		log.Printf("[AppendEntries] Appending log with index: %v, term: %v, key: %v, value: %v, command: %v", len(s.Node.info.log)-1, entry.Term, entry.Key, entry.Value, entry.Command)
		s.Node.info.log = append(s.Node.info.log, newLog)

		if !overwrite {
			log.Printf("[Persistence] writing log via append entries")
			s.Node.SaveLog(newLog)
		}

		if entry.Command == int32(NewOldConfig) {
			oldnewConfigIdx = index
		} else if entry.Command == int32(NewConfig) {
			newConfigIdx = index
		}
	}

	log.Printf("[AppendEntries] oldnewConfigIdx: %v, newConfigIdx: %v", oldnewConfigIdx, newConfigIdx)
	if oldnewConfigIdx > newConfigIdx {
		s.Node.info.isJointConsensus = true
		s.Node.info.newClusterAddresses = parseAddresses(in.Entries[oldnewConfigIdx].Value)
		s.Node.info.newClusterCount = len(s.Node.info.newClusterAddresses)
		s.Node.info.newId = -1
		for index, address := range s.Node.info.newClusterAddresses {
			if s.Node.address.String() == address.String() {
				s.Node.info.newId = index
				break
			}
		}
	} else if newConfigIdx > oldnewConfigIdx {
		if oldnewConfigIdx != -1 {
			log.Printf("[DEBUG] oldnewindex: %v, newindex, %v, entries: %v", oldnewConfigIdx, newConfigIdx, in.Entries)
			s.Node.info.newClusterAddresses = parseAddresses(in.Entries[oldnewConfigIdx].Value)
			s.Node.info.newClusterCount = len(s.Node.info.newClusterAddresses)
			s.Node.info.newId = -1
			for index, address := range s.Node.info.newClusterAddresses {
				if s.Node.address.String() == address.String() {
					s.Node.info.newId = index
					break
				}
			}

			s.Node.info.isJointConsensus = false
			s.Node.info.clusterCount = s.Node.info.newClusterCount
			s.Node.info.clusterAddresses = s.Node.info.newClusterAddresses

			if s.Node.info.newId == -1 {
				defer s.Node.Stop()
			} else {
				s.Node.info.id = s.Node.info.newId
			}
		} else {
			log.Printf("[AppendEntries] change to newConfig state, but no oldnewConfigIdx")
			s.Node.info.isJointConsensus = false
			s.Node.info.clusterCount = s.Node.info.newClusterCount
			s.Node.info.clusterAddresses = s.Node.info.newClusterAddresses
			s.Node.info.id = -1
			for index, address := range s.Node.info.clusterAddresses {
				if s.Node.address.String() == address.String() {
					s.Node.info.id = index
					break
				}
			}

			if s.Node.info.id == -1 {
				defer s.Node.Stop()
			}
		}
	}

	if overwrite {
		log.Printf("[Persistence] rewriting logs because of conflict")
		s.Node.SaveLogs()
	}

	// If leaderCommit > commitIndex, set commitIndex =
	// min(leaderCommit, index of last new entry)
	if in.LeaderCommit > int32(s.Node.info.commitIndex) {
		newCommitIndex := min(int(in.LeaderCommit), len(s.Node.info.log)-1)
		s.Node.CommitLogEntries(newCommitIndex)
	}

	return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: true}, nil
}

// TODO: handle joint election
func (s *server) RequestVote(ctx context.Context, in *comm.RequestVoteRequest) (*comm.RequestVoteResponse, error) {
	s.serverMutex.Lock()

	vote := false
	term := s.Node.info.currentTerm

	lastLogIdx := int32(len(s.Node.info.log) - 1)
	lastLogTerm := int32(0)
	if lastLogIdx > -1 {
		lastLogTerm = int32(s.Node.info.log[lastLogIdx].Term)
	}

	if term < int(in.Term) && (in.LastLogTerm > lastLogTerm || (in.LastLogTerm == lastLogTerm && in.LastLogIndex >= lastLogIdx)) {
		log.Printf("[Election] Giving vote for ID: %v, TERM: %v", in.CandidateId, in.Term)
		log.Printf("[Election] [UPDATE TERM] Before: %v, After: %v", s.Node.info.currentTerm, in.Term)
		vote = true
		s.Node.info.votedFor = int(in.CandidateId)
		s.Node.electionResetSignal <- true
		s.Node.info.currentTerm = int(in.Term)
	} else {
		log.Printf("[Election] Declined vote for ID: %v, TERM: %v", in.CandidateId, in.Term)
	}
	s.serverMutex.Unlock("[RequestVote]")
	return &comm.RequestVoteResponse{Term: int32(term), VoteGranted: vote}, nil
}

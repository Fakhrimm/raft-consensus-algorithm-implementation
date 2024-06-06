package node

import (
	"Node/grpc/comm"
	"context"
	"log"
	"net"
	"strings"
)

type server struct {
	comm.UnimplementedCommServiceServer
	Node *Node
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
	log.Printf("[Transaction] Received ping")
	return &comm.BasicResponse{Code: 200, Message: "PONG"}, nil
}

func (s *server) Stop(ctx context.Context, in *comm.BasicRequest) (*comm.BasicResponse, error) {
	defer s.Node.Stop()
	log.Printf("[Transaction] Received stop signal")
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

		s.Node.info.log = append(s.Node.info.log, comm.Entry{
			Term:    int32(s.Node.info.currentTerm),
			Key:     in.Key,
			Value:   in.Value,
			Command: int32(Set),
		})

		transactionIndex := len(s.Node.info.log) - 1
		s.Node.info.matchIndex[s.Node.info.id] = len(s.Node.info.log) - 1

		for index, nextIndex := range s.Node.info.nextIndex {
			if nextIndex == s.Node.info.matchIndex[index] {
				s.Node.info.nextIndex[index]++
			}
		}
		log.Printf("[Transaction] matchIndex: %v", s.Node.info.matchIndex)
		log.Printf("[Transaction] nextIndex: %v", s.Node.info.nextIndex)

		for transactionIndex != s.Node.info.commitIndex {
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

		s.Node.info.log = append(s.Node.info.log, comm.Entry{
			Term:    int32(s.Node.info.currentTerm),
			Key:     in.Key,
			Value:   "",
			Command: int32(Delete),
		})

		transactionIndex := len(s.Node.info.log) - 1
		s.Node.info.matchIndex[s.Node.info.id] = len(s.Node.info.log) - 1

		for index, nextIndex := range s.Node.info.nextIndex {
			if nextIndex == s.Node.info.matchIndex[index] {
				s.Node.info.nextIndex[index]++
			}
		}
		log.Printf("[Transaction] matchIndex: %v", s.Node.info.matchIndex)
		log.Printf("[Transaction] nextIndex: %v", s.Node.info.nextIndex)

		for transactionIndex != s.Node.info.commitIndex {
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

		s.Node.info.log = append(s.Node.info.log, comm.Entry{
			Term:    int32(s.Node.info.currentTerm),
			Key:     in.Key,
			Value:   in.Value,
			Command: int32(Append),
		})

		transactionIndex := len(s.Node.info.log) - 1
		s.Node.info.matchIndex[s.Node.info.id] = len(s.Node.info.log) - 1

		for index, nextIndex := range s.Node.info.nextIndex {
			if nextIndex == s.Node.info.matchIndex[index] {
				s.Node.info.nextIndex[index]++
			}
		}
		log.Printf("[Transaction] matchIndex: %v", s.Node.info.matchIndex)
		log.Printf("[Transaction] nextIndex: %v", s.Node.info.nextIndex)

		for transactionIndex != s.Node.info.commitIndex {
		}
	}

	log.Printf("[Transaction] Append request is completed")
	return &comm.AppendValueResponse{Code: code, Message: message, Value: value}, nil
}

func (s *server) ChangeMembership(ctx context.Context, in *comm.ChangeMembershipRequest) (*comm.ChangeMembershipResponse, error) {
	log.Printf("[Transaction] Received membership change request")

	// clusterAddresses & newClusterAddresses example:
	// "10.1.78.242:60000,10.1.78.242:60001,10.1.78.242:60002"
	// (separated by comma)
	var code int32
	var message string

	var clusterAddresses []net.TCPAddr
	var newClusterAddresses []net.TCPAddr

	// Validate the request first
	code, message = s.ValidateRequest()
	if code == 200 {
		// Parse the (old) addresses and new addresses
		clusterAddresses, newClusterAddresses = parseAddresses(in.ClusterAddresses, in.NewClusterAddresses)

		// Set it in the node
		s.Node.info.clusterAddresses = clusterAddresses
		s.Node.info.newClusterAddresses = newClusterAddresses
		s.Node.info.clusterCount = len(clusterAddresses)
		s.Node.info.newClusterCount = len(clusterAddresses)

		// Set that the node is in the joint consensus state
		s.Node.info.isJointConsensus = true

		// Append entry to log
		s.Node.info.log = append(s.Node.info.log, comm.Entry{
			Term:    int32(s.Node.info.currentTerm),
			Value:   in.NewClusterAddresses,
			Command: int32(NewOldConfig),
		})
	}

	return &comm.ChangeMembershipResponse{Code: code, Message: message, Value: in.NewClusterAddresses}, nil
}

func parseAddresses(clusterAddresses string, newClusterAddresses string) ([]net.TCPAddr, []net.TCPAddr) {
	var clusterAddressesTcp []net.TCPAddr
	var newClusterAddressesTcp []net.TCPAddr

	if clusterAddresses != "" {

		addresses := strings.Split(clusterAddresses, ",")

		for _, address := range addresses {
			a, err := net.ResolveTCPAddr("tcp", address)
			if err != nil {
				log.Fatalf("Invalid address: %v", address)
			}

			clusterAddressesTcp = append(clusterAddressesTcp, *a)
		}
	}

	if newClusterAddresses != "" {

		newAddresses := strings.Split(newClusterAddresses, ",")

		for _, newAddress := range newAddresses {
			a, err := net.ResolveTCPAddr("tcp", newAddress)
			if err != nil {
				log.Fatalf("Invalid address: %v", newAddress)
			}

			clusterAddressesTcp = append(clusterAddressesTcp, *a)
		}
	}

	return clusterAddressesTcp, newClusterAddressesTcp
}

// Raft purposes
func (s *server) AppendEntries(ctx context.Context, in *comm.AppendEntriesRequest) (*comm.AppendEntriesResponse, error) {
	//Reply false if term from leader < currentTerm (§5.1)
	if in.Term < int32(s.Node.info.currentTerm) {
		return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: false}, nil
	}

	// Reset election timer
	s.Node.onHeartBeat()
	s.Node.info.leaderId = int(in.LeaderId)

	// Reply false if log does not contain an entry at prevLogIndex
	if len(s.Node.info.log) < int(in.PrevLogIndex) {
		return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: false}, nil
	}

	if in.PrevLogIndex >= 0 {
		// Reply false if log entry term at prevLogIndex does not match prevLogTerm (§5.3)
		if s.Node.info.log[in.PrevLogIndex].Term != in.PrevLogTerm {
			return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: false}, nil
		}
	}

	if len(s.Node.info.log) > 0 {
		maxIdx := max(0, in.PrevLogIndex+1)
		s.Node.info.log = s.Node.info.log[:maxIdx]
	}

	for _, entry := range in.Entries {
		newLog := comm.Entry{
			Term:    entry.Term,
			Key:     entry.Key,
			Value:   entry.Value,
			Command: entry.Command,
		}
		s.Node.info.log = append(s.Node.info.log, newLog)

		if entry.Command == int32(NewOldConfig) {
			s.Node.info.isJointConsensus = true
			_, s.Node.info.newClusterAddresses = parseAddresses("", entry.Value)
			s.Node.info.newClusterCount = len(s.Node.info.newClusterAddresses)

		} else if entry.Command == int32(NewConfig) {
			s.Node.info.isJointConsensus = false
			s.Node.info.clusterCount = s.Node.info.newClusterCount
			s.Node.info.clusterAddresses = s.Node.info.newClusterAddresses
		}
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
	log.Printf("")

	vote := false
	term := s.Node.info.currentTerm

	lastLogIdx := int32(len(s.Node.info.log) - 1)
	lastLogTerm := int32(0)
	if lastLogIdx > -1 {
		lastLogTerm = int32(s.Node.info.log[lastLogIdx].Term)
	}

	if term < int(in.Term) && (in.LastLogTerm > lastLogTerm || (in.LastLogTerm == lastLogTerm && in.LastLogIndex >= lastLogIdx)) {
		log.Printf("[Election] Giving vote for id %v", in.CandidateId)
		vote = true
		s.Node.info.votedFor = int(in.CandidateId)
		s.Node.electionResetSignal <- true
		s.Node.info.currentTerm = int(in.Term)
	} else {
		log.Printf("[Election] Node %v requests for vote, declined", in.CandidateId)
	}

	return &comm.RequestVoteResponse{Term: int32(term), VoteGranted: vote}, nil
}

package node

import (
	"Node/grpc/comm"
	"context"
	"log"
)

type server struct {
	comm.UnimplementedCommServiceServer
	Node *Node
}

// Application purposes
func (s *server) Ping(ctx context.Context, in *comm.BasicRequest) (*comm.BasicResponse, error) {
	return &comm.BasicResponse{Code: 0, Message: "PONG"}, nil
}

func (s *server) Stop(ctx context.Context, in *comm.BasicRequest) (*comm.BasicResponse, error) {
	defer s.Node.Stop()
	return &comm.BasicResponse{Code: 0, Message: "STOPPING"}, nil
}

func (s *server) GetValue(ctx context.Context, in *comm.GetValueRequest) (*comm.GetValueResponse, error) {
	return &comm.GetValueResponse{Code: 0, Value: s.Node.app.Get(in.Key)}, nil
}

func (s *server) SetValue(ctx context.Context, in *comm.SetValueRequest) (*comm.SetValueResponse, error) {
	s.Node.app.Set(in.Key, in.Value)
	return &comm.SetValueResponse{Code: 0, Message: "Value changed sucessfully"}, nil
}

func (s *server) StrlnValue(ctx context.Context, in *comm.StrlnValueRequest) (*comm.StrlnValueResponse, error) {
	length := s.Node.app.Strln(in.Key)
	return &comm.StrlnValueResponse{Code: 0, Message: "Length fetched", Value: int32(length)}, nil
}

func (s *server) DeleteValue(ctx context.Context, in *comm.DeleteValueRequest) (*comm.DeleteValueResponse, error) {
	value := s.Node.app.Del(in.Key)
	return &comm.DeleteValueResponse{Code: 0, Message: "Value Deleted", Value: value}, nil
}

func (s *server) AppendValue(ctx context.Context, in *comm.AppendValueRequest) (*comm.AppendValueResponse, error) {
	s.Node.app.Append(in.Key, in.Value)
	return &comm.AppendValueResponse{Code: 0, Message: "Value appended sucessfully"}, nil
}

// Raft purposes

func (s *server) AppendEntries(ctx context.Context, in *comm.AppendEntriesRequest) (*comm.AppendEntriesResponse, error) {
	//Reply false if term from leader < currentTerm (§5.1)
	if in.Term < int32(s.Node.info.currentTerm) {
		return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: false}, nil
	}

	// Reset election timer
	s.Node.onHeartBeat()

	// Reply false if log doesn’t contain an entry at prevLogIndex
	if len(s.Node.info.log) < int(in.PrevLogIndex) {
		return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: false}, nil
	}

	if in.PrevLogIndex >= 0 {
		// Reply false if log entry term at prevLogIndex does not match prevLogTerm (§5.3)
		if s.Node.info.log[in.PrevLogIndex].term != int(in.PrevLogTerm) {
			return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: false}, nil
		}
	}

	s.Node.info.log = s.Node.info.log[:in.PrevLogIndex+1]

	for _, entry := range in.Entries {
		newLog := LogEntry{
			term:  int(entry.Term),
			key:   entry.Key,
			value: entry.Value,
		}
		s.Node.info.log = append(s.Node.info.log, newLog)
	}

	// If leaderCommit > commitIndex, set commitIndex =
	// min(leaderCommit, index of last new entry)
	if in.LeaderCommit > int32(s.Node.info.commitIndex) {
		newCommitIndex := min(int(in.LeaderCommit), len(s.Node.info.log)-1)
		s.Node.CommitLogEntries(newCommitIndex)
	}

	return &comm.AppendEntriesResponse{Term: int32(s.Node.info.currentTerm), Success: true}, nil
}

func (s *server) RequestVote(ctx context.Context, in *comm.RequestVoteRequest) (*comm.RequestVoteResponse, error) {
	log.Printf("")

	vote := false
	term := s.Node.info.currentTerm

	if term < int(in.Term) {
		log.Printf("[Election] Giving vote for id %v", in.CandidateId)
		vote = true
		s.Node.info.votedFor = int(in.CandidateId)
		s.Node.resetElectionTimer()
	} else {
		log.Printf("[Election] Node %v requests for vote, declined", in.CandidateId)
	}

	return &comm.RequestVoteResponse{Term: int32(term), VoteGranted: vote}, nil
}

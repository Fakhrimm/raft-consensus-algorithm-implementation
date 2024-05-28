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
	len := s.Node.app.Strln(in.Key)
	return &comm.StrlnValueResponse{Code: 0, Message: "Length fetched", Value: int32(len)}, nil
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
	// TODO: Implement
	return &comm.AppendEntriesResponse{}, nil
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

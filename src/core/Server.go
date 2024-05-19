package core

import (
	"context"

	pb "raft-sister/src/proto/comm"
)

type server struct {
	pb.UnimplementedCommServiceServer
	Node *Node
}

func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Code: 0, Message: "PONG"}, nil
}

func (s *server) Stop(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	defer s.Node.Stop()
	return &pb.PingResponse{Code: 0, Message: "STOPPING"}, nil
}

func (s *server) RequestValue(ctx context.Context, in *pb.RequestValueRequest) (*pb.RequestValueResponse, error) {
	return &pb.RequestValueResponse{Code: 0, Value: s.Node.Map[in.Key]}, nil
}

func (s *server) SetValue(ctx context.Context, in *pb.SetValueRequest) (*pb.SetValueResponse, error) {
	s.Node.Map[in.Key] = in.Value
	return &pb.SetValueResponse{Code: 0, Message: "Value changed sucessfully"}, nil
}

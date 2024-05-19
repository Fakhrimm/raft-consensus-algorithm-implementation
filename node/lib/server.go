package node

import (
	"Node/grpc/comm"
	"context"
)

type server struct {
	comm.UnimplementedCommServiceServer
	Node *Node
}

func (s *server) Ping(ctx context.Context, in *comm.PingRequest) (*comm.PingResponse, error) {
	return &comm.PingResponse{Code: 0, Message: "PONG"}, nil
}

func (s *server) Stop(ctx context.Context, in *comm.PingRequest) (*comm.PingResponse, error) {
	defer s.Node.Stop()
	return &comm.PingResponse{Code: 0, Message: "STOPPING"}, nil
}

func (s *server) RequestValue(ctx context.Context, in *comm.RequestValueRequest) (*comm.RequestValueResponse, error) {
	return &comm.RequestValueResponse{Code: 0, Value: s.Node.app.Get(in.Key)}, nil
}

func (s *server) SetValue(ctx context.Context, in *comm.SetValueRequest) (*comm.SetValueResponse, error) {
	s.Node.app.Set(in.Key, in.Value)
	return &comm.SetValueResponse{Code: 0, Message: "Value changed sucessfully"}, nil
}

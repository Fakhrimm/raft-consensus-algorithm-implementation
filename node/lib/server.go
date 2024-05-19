package node

import (
	"Node/grpc/comm"
	"context"
)

type server struct {
	comm.UnimplementedCommServiceServer
	Node *Node
}

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

func (s *server) AppendValue(ctx context.Context, in *comm.AppendValueRequest) (*comm.AppendValueResponse, error) {
	s.Node.app.Append(in.Key, in.Value)
	return &comm.AppendValueResponse{Code: 0, Message: "Value appended sucessfully"}, nil
}

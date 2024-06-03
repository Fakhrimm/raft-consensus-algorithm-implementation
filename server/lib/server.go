package server

import (
	"bufio"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	pb "server/grpc/comm"
)

type Server struct {
	pb.UnimplementedCommServiceServer

	client pb.CommServiceClient

	reader *bufio.Reader
}

func NewServer() *Server {
	server := &Server{
		reader: bufio.NewReader(os.Stdin),
	}
	return server
}

func (s *Server) Ping(ctx context.Context, in *pb.BasicRequest) (*pb.BasicResponse, error) {
	// TODO : isi address manual
	address := ""
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	s.client = pb.NewCommServiceClient(conn)

	response, err := s.client.Ping(ctx, in)

	return response, err
}

func (s *Server) GetValue(ctx context.Context, in *pb.GetValueRequest) (*pb.GetValueResponse, error) {
	// TODO: isi address manual
	address := ""
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	s.client = pb.NewCommServiceClient(conn)

	response, err := s.client.GetValue(ctx, in)
	return response, err
}

func (s *Server) SetValue(ctx context.Context, in *pb.SetValueRequest) (*pb.SetValueResponse, error) {
	// TODO: isi address manual
	address := ""
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	s.client = pb.NewCommServiceClient(conn)

	response, err := s.client.SetValue(ctx, in)
	return response, err
}

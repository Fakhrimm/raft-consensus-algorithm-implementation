package core

import (
	"context"
	"log"
	"strconv"

	pb "raft-sister/src/proto/example"
)

type server struct {
	pb.UnimplementedExampleServiceServer
	Node *Node
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Responded to: %v", in.Name)
	return &pb.HelloResponse{Message: "Hello from " + strconv.Itoa(s.Node.Port)}, nil
}

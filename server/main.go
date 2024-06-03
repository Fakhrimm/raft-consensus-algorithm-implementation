package main

import (
	"flag"
	"fmt"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	pb "server/grpc/comm"
	"server/lib"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// TODO: Implement dynamic address
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCommServiceServer(s, &server.Server{})
	log.Printf("grpc server listening at %v", lis.Addr())
	go func() {
		log.Fatalf("failed to serve: %v", s.Serve(lis))
	}()

	grpcWebServer := grpcweb.WrapServer(
		s,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	srv := &http.Server{
		Handler: grpcWebServer,
		Addr:    fmt.Sprintf("localhost:%d", *port+1),
	}
	log.Printf("grpc-web server listening at %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

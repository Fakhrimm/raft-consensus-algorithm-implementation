package main

import (
	node "Node/lib"
	"flag"
	"fmt"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"log"
	"net/http"

	pb "Node/grpc/comm"
)

func main() {
	var addr string
	var timeoutAvg int
	var hostfile string
	var hostfileNew string
	var isJointConsensus bool

	flag.StringVar(&addr, "addr", "0.0.0.0:60000", "The port on which the server will listen")
	flag.IntVar(&timeoutAvg, "timeout", 200, "The average time before a node declares a timeout")
	flag.StringVar(&hostfile, "hostfile", "/config/Hostfile", "The list of all servers in the cluster")
	flag.StringVar(&hostfileNew, "hostfilenew", "", "The list of new servers in the cluster")
	flag.BoolVar(&isJointConsensus, "isjoinconsensus", false, "The list of new servers in the cluster")
	flag.Parse()

	var service pb.CommServiceServer
	grpcServer := grpc.NewServer()
	pb.RegisterCommServiceServer(
		grpcServer,
		service,
	)

	node := node.NewNode(addr)
	node.Init(".."+hostfile, timeoutAvg, ".."+hostfileNew, isJointConsensus)

	wrappedGrpc := grpcweb.WrapServer(grpcServer)
	httpServer := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 {
				wrappedGrpc.ServeHTTP(w, r)
			} else {
				http.Error(w, "Only gRPC requests are accepted", http.StatusNotImplemented)
			}
		}),
	}
	fmt.Println("Starting server on port 8080")
	log.Fatal(httpServer.ListenAndServe())

	for node.Running {
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "golang/user/usertask"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *Server) SayHello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloResponse, error) {

	return &pb.HelloResponse{
		Message: "Hello " + req.Name,
	}, nil
}
func main() {
	lis, err := net.Listen("tcp",":50051")
	if err!=nil{
		log.Fatal(err)
	}
	grpcServer:=grpc.NewServer()
	pb.RegisterHelloServiceServer(
		grpcServer,
		&Server{},
	)
	fmt.Println("Server running on :50051")
	if err:=grpcServer.Serve(lis);err!=nil{
		log.Fatal(err)
	}
}
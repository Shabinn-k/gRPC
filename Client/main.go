package main

import (
	"context"
	"log"
	"time"

	pb "golang/user/usertask"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	res, err := client.SayHello(
		ctx,
		&pb.HelloRequest{
			Name: "Shabin",
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.Message)
}
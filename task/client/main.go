package main

import (
	"context"
	"fmt"
	"golang/task/pb"
	"log"

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
	if err!=nil{
		log.Fatal(err)
	}
	defer conn.Close()
	client:=pb.NewTaskServiceClient(conn)
	res,err:=client.CreateTask(
		context.Background(),
		&pb.CreateTaskRequest{
			Title:"Learn grpc",
			Description:"Practice crud",
		},
	)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("ID",res.Id)
	fmt.Println("TITLE",res.Title)
	fmt.Println("DESCRIPTION",res.Description)
}
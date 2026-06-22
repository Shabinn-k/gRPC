package main

import (
	"context"
	"golang/task/db"
	"golang/task/pb"
	"golang/task/repository"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedTaskServiceServer
	repo *repository.TaskRespository
}

func (s *Server)CreateTask(ctx context.Context,req *pb.CreateTaskRequest)(*pb.CreateTaskResponse,error){
	id,err:=s.repo.CreateTask(
		req.Title,
		req.Description,
	)
	if err!=nil{
		return nil,err
	}
	return &pb.CreateTaskResponse{
		Id:id,
		Title:req.Title,
		Description:req.Description,
	},nil
}

func main(){
	conn,err:=db.ConnectDB()
	if err!=nil{
		log.Fatal(err)
	}
	repo:=&repository.TaskRespository{
		DB: conn,
	}
	lis,err:=net.Listen("tcp",":50051")
	if err!=nil{
		log.Fatal(err)
	}
	grpcServer:=grpc.NewServer()
	pb.RegisterTaskServiceServer(
		grpcServer,
		&Server{
			repo: repo,
		},
	)
	log.Println("Server running on port 50051")
	grpcServer.Serve(lis)
}
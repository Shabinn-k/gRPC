package main

import (
	"context"
	"log"
	"net"

	pb "golang/Project/Logger"
	"golang/Project/auth"
	"golang/Project/repository"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServiceServer
	userRepository *repository.UserRepository
}

func NewServer(userRepository *repository.UserRepository) *server {
	return &server{
		userRepository: userRepository,
	}
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	userID, err := s.userRepository.ValidateUser(req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *server) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	userID, err := auth.ValidateToken(req.GetToken())
	if err != nil {
		return &pb.TokenResponse{
			Valid: false,
		}, nil
	}

	return &pb.TokenResponse{
		UserId: userID,
		Valid:  true,
	}, nil
}

func main() {
	// Initialize repository
	userRepo := repository.NewUserRepository()

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, NewServer(userRepo))

	// Start listening
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
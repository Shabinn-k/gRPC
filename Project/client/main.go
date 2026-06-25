package main

import (
	"context"
	"log"
	"time"

	pb "golang/Project/Logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() { 
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
 
	log.Println("Testing Login...")
	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Email:    "user1@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	log.Printf("Login successful! Token: %s", loginResp.Token)
 
	log.Println("\nTesting ValidateToken...")
	validateResp, err := client.ValidateToken(ctx, &pb.TokenRequest{
		Token: loginResp.Token,
	})
	if err != nil {
		log.Fatalf("Validate token failed: %v", err)
	}
	log.Printf("Token validation result - UserID: %d, Valid: %t", validateResp.UserId, validateResp.Valid)
 
	log.Println("\nTesting ValidateToken with invalid token...")
	invalidResp, err := client.ValidateToken(ctx, &pb.TokenRequest{
		Token: "invalid-token",
	})
	if err != nil {
		log.Printf("Validate token error: %v", err)
	}
	log.Printf("Invalid token validation - Valid: %t", invalidResp.Valid)
}
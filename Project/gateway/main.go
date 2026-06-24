package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "golang/Project/Logger"
	"golang/Project/middleware"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	userServiceClient pb.UserServiceClient
}

func NewGateway() *Gateway { 
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	client := pb.NewUserServiceClient(conn)

	return &Gateway{
		userServiceClient: client,
	}
}

func (g *Gateway) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := g.userServiceClient.Login(ctx, grpcReq)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": resp.Token,
	})
}

func (g *Gateway) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(int32)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "This is a protected endpoint!",
		"user_id": userID,
	})
}

func main() {
	gateway := NewGateway()

	// Create router
	router := mux.NewRouter()

	// Public endpoints
	router.HandleFunc("/login", gateway.LoginHandler).Methods("POST")

	// Protected endpoints (with auth middleware)
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)
	protectedRouter.HandleFunc("/protected", gateway.ProtectedHandler).Methods("GET")

	// Start HTTP server
	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
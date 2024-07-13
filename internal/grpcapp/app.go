package grpcapp

import (
	"auth-service/internal/config"
	"auth-service/jwt"
	"auth-service/protobuffs/auth-service"
	"context"
	"encoding/json"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	Server *grpc.Server
	config *config.Config
	auth.UnimplementedAuthServiceServer
}

func New(c *config.Config) *GrpcServer {
	log.Print("gRPC server started")

	s := grpc.NewServer()
	auth.RegisterAuthServiceServer(s, &GrpcServer{config: c})

	return &GrpcServer{Server: s, config: c}
}

func (s *GrpcServer) Hello(ctx context.Context, in *auth.Empty) (*auth.HelloResponse, error) {
	// Start the server
	return &auth.HelloResponse{Message: "Hello World!"}, nil
}

func (s *GrpcServer) VerifyToken(ctx context.Context, in *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error) {
	parsed, err := jwt.Parse(in.Token, s.config.Secret, "Bearer ")
	if err != nil {
		return nil, status.Error(codes.Unknown, "Invalid token")
	}

	expired, err := jwt.IsExpired(parsed)
	if err != nil {
		return nil, status.Error(codes.Internal, "Invalid token")
	}
	if expired {
		return &auth.VerifyTokenResponse{Expired: true}, nil
	}

	b, err := json.Marshal(parsed)
	if err != nil {
		return nil, status.Error(codes.Internal, "Marshal failed")
	}

	return &auth.VerifyTokenResponse{Claims: string(b), Expired: false}, nil
}

package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bogdanove/auth/internal/config"
	"github.com/bogdanove/auth/internal/config/env"
	"github.com/bogdanove/auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "prod.env", "path to config file")
}

type server struct {
	user_v1.UnimplementedUserV1Server
}

func (s *server) Create(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	log.Printf("start creating new user with name: %s", req.GetUserInfo().GetName())
	id := gofakeit.Int64()
	_ = ctx
	log.Printf("new user was created with id: %d", id)
	return &user_v1.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	log.Printf("receiving user with id: %d", req.GetId())
	_ = ctx
	now := time.Now()
	return &user_v1.GetResponse{
		User: &user_v1.User{
			Id:        req.GetId(),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Role:      user_v1.Role_user,
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
	}, nil
}

func (s *server) Update(ctx context.Context, req *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("updating user with id: %d", req.GetId())
	_ = ctx
	return new(emptypb.Empty), nil
}

func (s *server) Delete(ctx context.Context, req *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("deleting user with id: %d", req.GetId())
	_ = ctx
	return new(emptypb.Empty), nil
}

func main() {
	flag.Parse()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	listen, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	user_v1.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at: %v", listen.Addr())

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

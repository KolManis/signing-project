package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	signV1 "github.com/KolManis/signing-project/shared/pkg/proto/sign/v1"
	"github.com/KolManis/signing-project/sign-service/internal/api"
	signService "github.com/KolManis/signing-project/sign-service/internal/service/sign"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50051
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to connect server: %v\n", err)
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to closer server: %v\n", err)
			return
		}
	}()

	service := signService.NewService()
	api := api.NewAPI(service)

	s := grpc.NewServer()

	signV1.RegisterSignServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("gRPC server listning on %d\n", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down gRPC server ...")
	s.GracefulStop()
	log.Println("Server stopped")
}

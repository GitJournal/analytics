package main

import (
	"context"
	"log"
	"net"

	pb "github.com/gitjournal/analytics_backend/protos"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedAnalyticsServiceServer
}

func (s *server) SendData(ctx context.Context, in *pb.AnalyticsMessage) (*pb.AnalyticsReply, error) {
	log.Printf("Received: %v %v", in.GetAppId(), len(in.GetEvents()))
	return &pb.AnalyticsReply{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAnalyticsServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

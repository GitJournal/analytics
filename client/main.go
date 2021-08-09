package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/gitjournal/analytics_backend/protos"
	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1:8080"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Failed to dial: %v", err)
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := pb.NewAnalyticsServiceClient(conn)
	reply, err := client.SendData(ctx, &pb.AnalyticsMessage{AppId: "io.gitjournal-go"})
	if err != nil {
		log.Printf("Failed to send: %v", err)
	}

	fmt.Println(reply)
}

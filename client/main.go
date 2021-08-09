package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"time"

	pb "github.com/gitjournal/analytics_backend/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// address = "127.0.0.1:8080"
	address = "https://analyticsbackend-wetu2tkdpq-ew.a.run.app:8080"
)

func main() {
	// conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	var opts []grpc.DialOption
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Printf("Failed to dial: %v", err)
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := pb.NewAnalyticsServiceClient(conn)
	reply, err := client.SendData(ctx, &pb.AnalyticsMessage{AppId: "io.gitjournal-go"})
	if err != nil {
		log.Printf("Failed to send: %v", err)
	}

	fmt.Println(reply)
}

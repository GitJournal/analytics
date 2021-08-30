package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/gitjournal/analytics_backend/protos"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/oschwald/geoip2-golang"
)

const dbPath = "GeoLite2-City.mmdb"

var conn *pgx.Conn

type server struct {
	pb.UnimplementedAnalyticsServiceServer
}

func (s *server) SendData(ctx context.Context, in *pb.AnalyticsMessage) (*pb.AnalyticsReply, error) {
	p, _ := peer.FromContext(ctx)
	addr, ok := p.Addr.(*net.TCPAddr)
	if !ok {
		log.Fatal("Could not get IP")
	}
	clientIP := addr.IP

	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		xForwardFor := headers.Get("x-forwarded-for")
		if len(xForwardFor) > 0 && xForwardFor[0] != "" {
			ips := strings.Split(xForwardFor[0], ",")
			if len(ips) > 0 {
				clientIP = net.ParseIP(ips[0])
			}
		}
	}

	// FIXME: Avoid reloading the db for each request
	db, err := geoip2.Open(dbPath)
	if err != nil {
		log.Fatal("Opening GeoLit2 db:", err)
	}
	defer db.Close()

	record, err := db.City(clientIP)
	if err != nil {
		return &pb.AnalyticsReply{}, err
	}

	err = insertIntoPostgres(ctx, conn, record, in)
	if err != nil {
		return &pb.AnalyticsReply{}, err
	}

	return &pb.AnalyticsReply{}, nil
}

func main() {
	if !fileExists(dbPath) {
		log.Fatalf("GeoLite db not found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx := context.Background()
	conn, err = postgresConnect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	log.Printf("Connected to Postgres")

	s := grpc.NewServer()
	pb.RegisterAnalyticsServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

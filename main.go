package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gitjournal/analytics_backend/pb"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/getsentry/sentry-go"
	"github.com/oschwald/geoip2-golang"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_sentry "github.com/johnbellone/grpc-middleware-sentry"
)

const dbPath = "GeoLite2-City.mmdb"

var conn *pgx.Conn
var geoDb *geoip2.Reader

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

	record, err := geoDb.City(clientIP)
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
	var err error
	if !fileExists(dbPath) {
		log.Fatalf("GeoLite db not found")
	}

	geoDb, err = geoip2.Open(dbPath)
	if err != nil {
		log.Fatal("Opening GeoLit2 db:", err)
	}
	defer geoDb.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:   "https://05ea3a469a04409db1eac1e6daf73479@o366485.ingest.sentry.io/5937572",
		Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	ctx := context.Background()
	conn, err = postgresConnect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer conn.Close(ctx)
	log.Printf("Connected to Postgres")

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_sentry.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_sentry.UnaryServerInterceptor(),
		)),
	)
	pb.RegisterAnalyticsServiceServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())
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

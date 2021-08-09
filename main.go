package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	pb "github.com/gitjournal/analytics_backend/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/oschwald/geoip2-golang"
)

const (
	port = ":50051"
)

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

	// ip, err := net.ResolveIPAddr(p.Addr.Network(), p.Addr.String())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("What the hell?", ip)

	db, err := geoip2.Open("./GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// If you are using strings that may be invalid, check that ip is not nil
	fmt.Println("Client IP", clientIP.To4().String())
	record, err := db.City(clientIP)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
	if len(record.Subdivisions) > 0 {
		fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
	}
	fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)

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

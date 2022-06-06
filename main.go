package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gitjournal/analytics_backend/pb"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/proto"

	"github.com/getsentry/sentry-go"
	"github.com/oschwald/geoip2-golang"
)

const dbPath = "GeoLite2-City.mmdb"

var conn *pgx.Conn
var geoDb *geoip2.Reader

func SendDataHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POSTs are accepted.", http.StatusBadRequest)
		return
	}

	clientIP, err := getIP(req)
	if err != nil {
		http.Error(w, "IP not found.", http.StatusBadRequest)
		return
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Errir reading body", http.StatusBadRequest)
		return
	}

	msg := &pb.AnalyticsMessage{}
	if err := proto.Unmarshal(bodyBytes, msg); err != nil {
		http.Error(w, "Failed to parse message: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	err = sendData(ctx, clientIP, msg)
	if err != nil {
		http.Error(w, "Failed to sendData: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIP(req *http.Request) (net.IP, error) {
	xForwardFor := req.Header.Get("X-Forwarded-For")
	if xForwardFor != "" {
		ips := strings.Split(xForwardFor, ",")
		if len(ips) > 0 {
			return net.ParseIP(ips[0]), nil
		}
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return net.IP{}, err
	}

	return net.ParseIP(ip), nil
}

func sendData(ctx context.Context, clientIP net.IP, in *pb.AnalyticsMessage) error {
	record, err := geoDb.City(clientIP)
	if err != nil {
		return err
	}

	err = insertIntoPostgres(ctx, conn, record, in)
	if err != nil {
		return err
	}

	return nil
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

	http.HandleFunc("/v1/sendData", SendDataHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
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

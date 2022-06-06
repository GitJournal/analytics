package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	pb "github.com/gitjournal/analytics_backend/pb"
	"github.com/gogo/protobuf/proto"
)

const (
	// address  = "https://analytics.gitjournal.io"
	address  = "http://127.0.0.1:8080"
	path     = "/v1/sendData"
	useLocal = true
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := createMessage()
	reqBytes, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, address+path, bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/x-protobuf")
	req = req.WithContext(ctx)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(resp)
}

func createMessage() *pb.AnalyticsMessage {
	device := &pb.DeviceInfo{
		Platform: pb.Platform_android,
		DeviceInfo: &pb.DeviceInfo_AndroidDeviceInfo{
			AndroidDeviceInfo: &pb.AndroidDeviceInfo{
				Board: "board",
			},
		},
	}

	packageInfo := &pb.PackageInfo{
		AppName:        "GitJournal",
		PackageName:    "io.gitjournal.gitjournal",
		Version:        "version",
		BuildNumber:    "123",
		BuildSignature: "signature",
	}

	events := []*pb.Event{
		{
			Name:      "test_event",
			Date:      1630323711,
			PseudoId:  "uuid",
			SessionID: 123,
		},
	}

	msg := &pb.AnalyticsMessage{
		AppId:       "io.gitjournal",
		DeviceInfo:  device,
		PackageInfo: packageInfo,
		Events:      events,
	}

	return msg
}

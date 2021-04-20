package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"
)

type Device struct {
	Category               string `json:"category"`
	MobileBrandName        string `json:"mobile_brand_name"`
	MobileModelName        string `json:"mobile_model_name"`
	MobileMarketingName    string `json:"mobile_marketing_name"`
	MobileOsHardwareModel  string `json:"mobile_os_hardware_model"`
	OperatingSystem        string `json:"operating_system"`
	OperatingSystemVersion string `json:"operating_system_version"`
	Language               string `json:"language"`
	TimeZoneOffsetSeconds  string `json:"time_zone_offset_seconds"`
}
type Geo struct {
	Continent    string `json:"continent"`
	Country      string `json:"country"`
	Region       string `json:"region"`
	City         string `json:"city"`
	SubContinent string `json:"sub_continent"`
	Metro        string `json:"metro"`
}
type AppInfo struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	// FirebaseAppID string `json:"firebase_app_id"`
	// InstallSource string `json:"install_source"`
}

type firebaseEvent struct {
	EventDate      string `json:"event_date"`
	EventTimestamp string `json:"event_timestamp"`
	EventName      string `json:"event_name"`
	EventParams    []struct {
		Key         string `json:"key"`
		StringValue struct {
			StringValue string `json:"string_value"`
		} `json:"value,omitempty"`
		IntValue struct {
			IntValue string `json:"int_value"`
		} `json:"value,omitempty"`
	} `json:"event_params"`
	EventPreviousTimestamp     string `json:"event_previous_timestamp"`
	EventBundleSequenceID      string `json:"event_bundle_sequence_id"`
	EventServerTimestampOffset string `json:"event_server_timestamp_offset"`
	UserPseudoID               string `json:"user_pseudo_id"`
	UserProperties             []struct {
		Key         string `json:"key"`
		StringValue struct {
			StringValue        string `json:"string_value"`
			SetTimestampMicros string `json:"set_timestamp_micros"`
		} `json:"value,omitempty"`
		IntValue struct {
			IntValue           string `json:"int_value"`
			SetTimestampMicros string `json:"set_timestamp_micros"`
		} `json:"value,omitempty"`
	} `json:"user_properties"`
	UserFirstTouchTimestamp string `json:"user_first_touch_timestamp"`
	Device                  `json:"device"`
	Geo                     `json:"geo"`
	AppInfo                 `json:"app_info"`
	StreamID                string        `json:"stream_id"`
	Platform                string        `json:"platform"`
	Items                   []interface{} `json:"items"`
}

type EventDBSchema struct {
	Timestamp string `json:"timestamp"`
	Name      string `json:"event_name"`

	// `params` Array(Tuple(String, String)),

	PreviousTimestamp time.Time `json:"previous_timestamp"`
	BundleSequenceId  int32     `json:"bundle_sequence_id"`

	UserPseudoId string `json:"user_pseudo_id"`
	// user props

	Device  `json:"device"`
	Geo     `json:"geo"`
	AppInfo `json:"app_info"`

	StreamID int64  `json:"stream_id"`
	Platform string `json:"platform"`
}

func mapEvent(e firebaseEvent) EventDBSchema {
	return EventDBSchema{
		Timestamp: e.EventTimestamp,
		Name:      e.EventName,
		// prev timestmp
		// BundleSequenceId: e.EventBundleSequenceID,
		UserPseudoId: e.UserPseudoID,
		Device:       e.Device,
		Geo:          e.Geo,
		AppInfo:      e.AppInfo,
	}
}

func main() {
	fileName := "bq-results-20210419-200559-1gju06yyz9u5.json"

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		var e firebaseEvent
		if err := json.Unmarshal(s.Bytes(), &e); err != nil {
			log.Fatal(err)
		}

		//fmt.Println(e)
	}
	if s.Err() != nil {
		// handle scan error
	}
}

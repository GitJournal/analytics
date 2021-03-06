package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
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

type Info map[string]interface{}

type firebaseEvent struct {
	EventDate                  string `json:"event_date"`
	EventTimestamp             string `json:"event_timestamp"`
	EventName                  string `json:"event_name"`
	EventParams                []Info `json:"event_params"`
	EventPreviousTimestamp     string `json:"event_previous_timestamp"`
	EventBundleSequenceID      string `json:"event_bundle_sequence_id"`
	EventServerTimestampOffset string `json:"event_server_timestamp_offset"`
	UserPseudoID               string `json:"user_pseudo_id"`
	UserProperties             []Info `json:"user_properties"`
	UserFirstTouchTimestamp    string `json:"user_first_touch_timestamp"`
	Device                     `json:"device"`
	Geo                        `json:"geo"`
	AppInfo                    `json:"app_info"`
	StreamID                   string        `json:"stream_id"`
	Platform                   string        `json:"platform"`
	Items                      []interface{} `json:"items"`
}

type EventDBSchema struct {
	Timestamp int64  `json:"timestamp"`
	Name      string `json:"name"`

	Params [][]string `json:"params"`

	PreviousTimestamp int64 `json:"previous_timestamp,omitempty"`
	BundleSequenceId  int32 `json:"bundle_sequence_id"`

	UserPseudoId string     `json:"user_pseudo_id"`
	UserProps    [][]string `json:"user_properties"`

	DeviceCategory               string `json:"device.category"`
	DeviceMobileBrandName        string `json:"device.mobile_brand_name"`
	DeviceMobileModelName        string `json:"device.mobile_model_name"`
	DeviceMobileMarketingName    string `json:"device.mobile_marketing_name"`
	DeviceMobileOsHardwareModel  string `json:"device.mobile_os_hardware_model"`
	DeviceOperatingSystem        string `json:"device.operating_system"`
	DeviceOperatingSystemVersion string `json:"device.operating_system_version"`
	DeviceLanguage               string `json:"device.language"`
	DeviceTimeZoneOffsetSeconds  string `json:"device.time_zone_offset_seconds"`

	GeoContinent    string `json:"geo.continent"`
	GeoCountry      string `json:"geo.country"`
	GeoRegion       string `json:"geo.region"`
	GeoCity         string `json:"geo.city"`
	GeoSubContinent string `json:"geo.sub_continent"`
	GeoMetro        string `json:"geo.metro"`

	AppInfoID      string `json:"app_info.id"`
	AppInfoVersion string `json:"app_info.version"`

	StreamID int64  `json:"stream_id"`
	Platform string `json:"platform"`
}

func mapParams(p []Info) [][]string {
	a := make([][]string, len(p))
	for i := range a {
		prop := p[i]
		a[i] = make([]string, 2)
		a[i][0] = prop["key"].(string)

		val := prop["value"].(map[string]interface{})
		if stringVal, ok := val["string_value"]; ok {
			a[i][1] = stringVal.(string)
		}
		if intVal, ok := val["int_value"]; ok {
			a[i][1] = intVal.(string)
		}

		if len(a[i][1]) == 0 {
			fmt.Println(prop)
			log.Panic("Why is the prop value empty?")
		}
	}

	return a
}

func mapEvent(e firebaseEvent) EventDBSchema {
	bundleSeq, err := strconv.ParseInt(e.EventBundleSequenceID, 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	streamID, err := strconv.ParseInt(e.StreamID, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	timestamp, err := strconv.ParseInt(e.EventTimestamp, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	timestamp /= 1000000

	prevTimestamp, err := strconv.ParseInt(e.EventPreviousTimestamp, 10, 64)
	if err != nil {
		prevTimestamp = 0
	} else {
		prevTimestamp /= 1000000
	}

	return EventDBSchema{
		Timestamp:         timestamp,
		Name:              e.EventName,
		Params:            mapParams(e.EventParams),
		PreviousTimestamp: prevTimestamp,

		BundleSequenceId: int32(bundleSeq),
		UserPseudoId:     e.UserPseudoID,
		UserProps:        mapParams(e.UserProperties),

		DeviceCategory:               e.Device.Category,
		DeviceMobileBrandName:        e.Device.MobileBrandName,
		DeviceMobileModelName:        e.Device.MobileModelName,
		DeviceMobileMarketingName:    e.Device.MobileMarketingName,
		DeviceMobileOsHardwareModel:  e.Device.MobileOsHardwareModel,
		DeviceOperatingSystem:        e.Device.OperatingSystem,
		DeviceOperatingSystemVersion: e.Device.OperatingSystemVersion,
		DeviceLanguage:               e.Device.Language,
		DeviceTimeZoneOffsetSeconds:  e.Device.TimeZoneOffsetSeconds,

		GeoContinent:    e.Geo.Continent,
		GeoCountry:      e.Geo.Country,
		GeoRegion:       e.Geo.Region,
		GeoCity:         e.Geo.City,
		GeoSubContinent: e.Geo.SubContinent,
		GeoMetro:        e.Geo.Metro,

		AppInfoID:      e.AppInfo.ID,
		AppInfoVersion: e.AppInfo.Version,

		StreamID: streamID,
		Platform: e.Platform,
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

		ne := mapEvent(e)
		bytes, err := json.Marshal(ne)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(bytes))
		//fmt.Println(e)
	}
	if s.Err() != nil {
		// handle scan error
	}
}

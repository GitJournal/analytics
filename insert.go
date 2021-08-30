package main

import (
	"context"
	"fmt"
	"time"

	analytics_backend "github.com/gitjournal/analytics_backend/protos"
	pb "github.com/gitjournal/analytics_backend/protos"
	"github.com/oschwald/geoip2-golang"

	"github.com/jackc/pgx/v4"

	"encoding/json"

	"github.com/twmb/murmur3"
)

func insertIntoPostgres(ctx context.Context, conn *pgx.Conn, cityInfo *geoip2.City, in *pb.AnalyticsMessage) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	//
	// Device
	//
	deviceID, err := getDeviceID(in)
	if err != nil {
		return err
	}

	di := in.DeviceInfo
	android, _ := json.Marshal(di.GetAndroidDeviceInfo())
	ios, _ := json.Marshal(di.GetIosDeviceInfo())
	linux, _ := json.Marshal(di.GetLinuxDeviceInfo())
	macos, _ := json.Marshal(di.GetMacOSDeviceInfo())
	windows, _ := json.Marshal(di.GetWindowsDeviceInfo())
	web, _ := json.Marshal(di.GetWebBrowserInfo())

	platform := analytics_backend.Platform_name[int32(di.Platform)]

	_, err = tx.Exec(ctx, "insert into analytics_device_info(id, platform, android_info, ios_info, linux_info, macos_info, windows_info, web_info) values ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING", deviceID, platform, android, ios, linux, macos, windows, web)
	if err != nil {
		return fmt.Errorf("insert analytics_device_info failed: %w", err)
	}

	//
	// Package Info
	//
	packageId, err := getPackageID(in)
	if err != nil {
		return fmt.Errorf("getPackageID failed: %w", err)
	}
	pi := in.PackageInfo

	_, err = tx.Exec(ctx, "insert into analytics_package_info(id, appName, packageName, version, buildNumber, buildSignature) values ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING", packageId, pi.AppName, pi.PackageName, pi.Version, pi.BuildNumber, pi.BuildSignature)
	if err != nil {
		return fmt.Errorf("insert analytics_package_info failed: %w", err)
	}

	//
	// Location
	//
	locID := getLocationID(cityInfo)

	_, err = tx.Exec(ctx, "insert into analytics_location(city_geoname_id, city_name_en, country_code) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", locID, cityInfo.City.Names["en"], cityInfo.Country.IsoCode)
	if err != nil {
		return fmt.Errorf("insert analytics_location failed: %w", err)
	}

	//
	// Analytics
	//
	sql := "insert into analytics_events(ts, event_name, props, pseudoId, userId, user_props, session_id, location_id, device_id, package_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	for _, ev := range in.Events {
		time := time.Unix(int64(ev.Date), 0)

		_, err := tx.Exec(ctx, sql, time, ev.Name, ev.Params, ev.PseudoId, ev.UserId, ev.UserProperties, ev.SessionID, locID, deviceID, packageId)
		if err != nil {
			return fmt.Errorf("insert analytics_insert failed: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func getLocationID(cityInfo *geoip2.City) string {
	return fmt.Sprint(cityInfo.City.GeoNameID)
}

func getDeviceID(msg *pb.AnalyticsMessage) (string, error) {
	jsonBytes, err := json.Marshal(msg.DeviceInfo)
	if err != nil {
		return "", err
	}

	h32 := murmur3.New32()
	h32.Write(jsonBytes)
	return fmt.Sprint(h32.Sum32()), nil
}

func getPackageID(msg *pb.AnalyticsMessage) (string, error) {
	jsonBytes, err := json.Marshal(msg.PackageInfo)
	if err != nil {
		return "", err
	}

	h32 := murmur3.New32()
	h32.Write(jsonBytes)
	return fmt.Sprint(h32.Sum32()), nil
}

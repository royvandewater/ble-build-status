package main

import (
	"fmt"
	"log"
	"os"
	"time"

	DEBUG "github.com/computes/go-debug"
	"github.com/royvandewater/ble-build-status/circleci"
	"github.com/royvandewater/ble-build-status/ring"
)

var debug = DEBUG.Debug("ble-build-status:main")

func fatalIfErrorf(err error, msg string, rest ...interface{}) {
	if err == nil {
		return
	}

	log.Fatalln(fmt.Sprintf(msg, rest...), err.Error())
}

func getEnvOr(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func main() {
	username := getEnvOr("BLE_BUILD_STATUS_USERNAME", "royvandewater")
	project := getEnvOr("BLE_BUILD_STATUS_PROJECT", "dummy-build")
	ringName := getEnvOr("BLE_BUILD_STATUS_RING_NAME", "esp32-neopixel")

	r, err := ring.New(ringName)
	fatalIfErrorf(err, "Failed to construct a new ring")

	err = r.Connect(10 * time.Second)
	fatalIfErrorf(err, "Failed to connect to ring")

	for {
		build, err := circleci.GetLatestBuild(username, project)
		fatalIfErrorf(err, "Failed to get latest build")

		debug("build status: %s", build.Status)
		switch build.Status {
		case "queued":
			debug("queued")
			err = r.PulseColor(0x00, 0xff, 0xff)
		case "running":
			debug("running")
			err = r.PulseColor(0xff, 0xff, 0x00)
		case "success":
			debug("success")
			err = r.SetColor(0x00, 0xff, 0x00)
		case "fixed":
			debug("fixed")
			err = r.SetColor(0x00, 0xff, 0x00)
		case "failed":
			debug("failed")
			err = r.SetColor(0xff, 0x00, 0x00)
		default:
			err = r.SetColor(0xff, 0xff, 0xff)
		}
		fatalIfErrorf(err, "Failed to set color")

		// err = r.Disconnect()
		// fatalIfErrorf(err, "Failed to disconnect")
		<-time.After(5 * time.Second)
	}
}

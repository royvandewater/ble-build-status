package main

import (
	"fmt"
	"log"
	"os"
	"time"

	DEBUG "github.com/computes/go-debug"
	"github.com/royvandewater/ble-build-status/circleci"
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
	project := getEnvOr("BLE_BUILD_STATUS_PROJECT", "clojure-for-the-brave-and-true")

	for {
		build, err := circleci.GetLatestBuild(username, project)
		fatalIfErrorf(err, "Failed to get latest build")

		debug("status: %s", build.Status)
		<-time.After(10 * time.Second)
	}
}

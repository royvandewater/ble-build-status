package ring

import (
	"fmt"

	"github.com/raff/goble"
)

func bleConnect(localName, serviceUUID string) (*goble.BLE, error) {
	ble := goble.New()

	ble.On("discover", func(ev goble.Event) bool {
		if localName != ev.Peripheral.Advertisement.LocalName {
			return false
		}

		ble.Connect(ev.DeviceUUID)
		return false
	})
	ble.On("stateChange", func(ev goble.Event) bool {
		debug("stateChanged: %v", ev)
		if ev.State == "poweredOn" {
			ble.StartScanning(nil, false)
			return false
		}

		ble.StopScanning()
		return true
	})
	ble.Init()

	return nil, fmt.Errorf("bleConnect not implemented")
}

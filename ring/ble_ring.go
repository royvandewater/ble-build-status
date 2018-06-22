package ring

import (
	"fmt"
	"time"

	"github.com/paypal/gatt"
	"github.com/royvandewater/ble-build-status/ring/option"
)

var _ServiceUUID = gatt.MustParseUUID("2ba75e8a-5b5b-447b-ab9a-b79e21dd64e0")
var _ColorCharacteristicUUID = gatt.MustParseUUID("08f490bf-28f1-4d55-897d-ab8d74effffb")
var _CommandCharacteristicUUID = gatt.MustParseUUID("04b29961-90fd-4ee7-bb48-f203bde84f44")

type _BLERing struct {
	device     gatt.Device
	localName  string
	peripheral gatt.Peripheral
}

func _NewBLERing(localName string) (*_BLERing, error) {
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		return nil, err
	}

	ring := &_BLERing{device: device, localName: localName}
	ring.initialize()
	return ring, nil
}

func (ring *_BLERing) Connect(timeout time.Duration) error {
	fmt.Println("Connecting...")
	err := ring.device.Init(ring.onStateChanged)
	if err != nil {
		return err
	}
	ring.device.Scan(nil, false)
	return nil
}

func (ring *_BLERing) initialize() {
	ring.device.Handle(
		gatt.PeripheralDiscovered(ring.onPeriphDiscovered),
		gatt.PeripheralConnected(ring.onPeriphConnected),
		gatt.PeripheralDisconnected(ring.onPeriphDisconnected),
	)
}

func (ring *_BLERing) onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if p.ID() != ring.localName {
		return
	}
	fmt.Println("Found it, connecting!")

	p.Device().StopScanning()
	p.Device().Connect(p)
}

func (ring *_BLERing) onPeriphConnected(p gatt.Peripheral, err error) {
	if err := p.SetMTU(500); err != nil {
		fmt.Printf("Failed to set MTU, err: %s\n", err)
		return
	}

	fmt.Println("Connected")
	ring.peripheral = p
}

func (ring *_BLERing) onPeriphDisconnected(p gatt.Peripheral, err error) {
}

func (ring *_BLERing) onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("state change", s)
	switch s {
	case gatt.StatePoweredOn:
		ring.device.Scan(nil, false)
		return
	default:
		ring.device.StopScanning()
	}
}

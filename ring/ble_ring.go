package ring

import (
	"time"

	"github.com/raff/goble"
)

var _ServiceUUID = "2ba75e8a-5b5b-447b-ab9a-b79e21dd64e0"
var _ColorCharacteristicUUID = "08f490bf-28f1-4d55-897d-ab8d74effffb"
var _CommandCharacteristicUUID = "04b29961-90fd-4ee7-bb48-f203bde84f44"

type _BLERing struct {
	ble       *goble.BLE
	localName string
}

func _NewBLERing(localName string) *_BLERing {
	return &_BLERing{localName: localName}
}

func (ring *_BLERing) Connect(timeout time.Duration) error {
	ble, err := bleConnect(ring.localName, _ServiceUUID)
	if err != nil {
		return err
	}
	ring.ble = ble
	return nil
}

package midi

import "github.com/leandrodaf/midi-client/internal/entity"

type ClientMIDI interface {
	Stop() error
	ListDevices() ([]entity.DeviceInfo, error)
	SelectDevice(deviceID int) error
	StartCapture(eventChannel chan entity.MIDI)
}

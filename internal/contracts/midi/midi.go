package midi

import "github.com/leandrodaf/midi-client/internal/entity"

// ClientMIDI defines an interface for managing MIDI client operations.
type ClientMIDI interface {
	Stop() error                                // Stops the MIDI client and disconnects any active device.
	ListDevices() ([]entity.DeviceInfo, error)  // Lists available MIDI devices.
	SelectDevice(deviceID int) error            // Selects a MIDI device by its ID.
	StartCapture(eventChannel chan entity.MIDI) // Starts capturing MIDI events and sends them to the specified channel.
}

//go:build darwin
// +build darwin

package mididarwin

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/leandrodaf/midi-client/internal/contracts/logger"
	"github.com/leandrodaf/midi-client/internal/contracts/midi"
	"github.com/leandrodaf/midi-client/internal/entity"
	"github.com/youpy/go-coremidi"
)

// internalPortConnection is an interface that defines the Disconnect method.
type internalPortConnection interface {
	Disconnect()
}

// ClientMid is the main structure that implements ClientMIDI.
type ClientMid struct {
	logger       logger.Logger
	eventChannel atomic.Value
	client       coremidi.Client
	inputPort    coremidi.InputPort
	portConn     internalPortConnection
	mu           sync.Mutex
}

// NewMIDIClient returns an instance of ClientMid as ClientMIDI.
func NewMIDIClient(logger logger.Logger) (midi.ClientMIDI, error) {
	client, err := coremidi.NewClient("GoMIDIClient")
	if err != nil {
		return nil, fmt.Errorf("failed to create MIDI client: %v", err)
	}
	logger.Info("MIDI client successfully created")
	return &ClientMid{
		logger: logger,
		client: client,
	}, nil
}

func (m *ClientMid) ListDevices() ([]entity.DeviceInfo, error) {
	sources, err := coremidi.AllSources()
	if err != nil {
		return nil, fmt.Errorf("error listing MIDI sources: %v", err)
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("no MIDI devices found")
	}

	devices := make([]entity.DeviceInfo, len(sources))
	for i, source := range sources {
		sourceEntity := source.Entity()
		devices[i] = entity.DeviceInfo{
			Name:         source.Name(),
			EntityName:   sourceEntity.Name(),
			Manufacturer: sourceEntity.Manufacturer(),
		}
	}
	return devices, nil
}

func (m *ClientMid) SelectDevice(deviceID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sources, err := coremidi.AllSources()
	if err != nil {
		return fmt.Errorf("error retrieving MIDI sources: %v", err)
	}
	if deviceID < 0 || deviceID >= len(sources) {
		return fmt.Errorf("invalid MIDI device")
	}

	if m.portConn != nil {
		m.portConn.Disconnect()
		m.portConn = nil
	}

	source := sources[deviceID]
	m.logger.Info("MIDI device selected",
		m.logger.Field().Int("deviceID", deviceID),
		m.logger.Field().String("deviceName", source.Name()))

	m.inputPort, err = coremidi.NewInputPort(m.client, "Input Port", m.handleMIDIMessage)
	if err != nil {
		return fmt.Errorf("error creating input port: %v", err)
	}

	m.portConn, err = m.inputPort.Connect(source)
	if err != nil {
		return fmt.Errorf("error connecting to MIDI device: %v", err)
	}

	m.logger.Info("MIDI device successfully connected")
	return nil
}

func (m *ClientMid) handleMIDIMessage(source coremidi.Source, packet coremidi.Packet) {
	eventChannel, _ := m.eventChannel.Load().(chan entity.MIDI)
	if eventChannel == nil {
		m.logger.Warn("eventChannel not initialized or of invalid type")
		return
	}

	if len(packet.Data) >= 3 {
		event := entity.MIDI{
			Timestamp: uint64(time.Now().UnixNano()),
			Command:   packet.Data[0],
			Note:      packet.Data[1],
			Velocity:  packet.Data[2],
		}
		eventChannel <- event
	} else {
		m.logger.Warn("Incomplete MIDI packet")
	}
}

func (m *ClientMid) Stop() error {
	m.logger.Info("Stopping MIDI capture")
	return m.stopCapture()
}

func (m *ClientMid) StartCapture(eventChannel chan entity.MIDI) {
	m.logger.Info("Capturing MIDI events")
	m.eventChannel.Store(eventChannel)
}

func (m *ClientMid) stopCapture() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.portConn != nil {
		m.portConn.Disconnect()
		m.portConn = nil
	}
	m.eventChannel.Store(nil)
	m.logger.Info("MIDI capture stopped")
	return nil
}

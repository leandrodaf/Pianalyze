//go:build !darwin
// +build !darwin

package mididarwin

import (
	"fmt"

	"github.com/leandrodaf/pianalyze/internal/contracts/logger"
	"github.com/leandrodaf/pianalyze/internal/contracts/midi"
	"github.com/leandrodaf/pianalyze/internal/entity"
)

// ClientMidDummy is a dummy implementation of the ClientMIDI interface for non-macOS systems.
type ClientMidDummy struct {
	logger logger.Logger
}

// NewMIDIClient returns a dummy instance of ClientMIDI.
func NewMIDIClient(logger logger.Logger) (midi.ClientMIDI, error) {
	logger.Info("Using dummy MIDI client for non-macOS system")
	return &ClientMidDummy{
		logger: logger,
	}, nil
}

func (m *ClientMidDummy) ListDevices() ([]entity.DeviceInfo, error) {
	m.logger.Warn("ListDevices called on dummy MIDI client")
	return nil, fmt.Errorf("MIDI functionality is not available on this platform")
}

func (m *ClientMidDummy) SelectDevice(deviceID int) error {
	m.logger.Warn("SelectDevice called on dummy MIDI client")
	return fmt.Errorf("MIDI functionality is not available on this platform")
}

func (m *ClientMidDummy) StartCapture(eventChannel chan entity.MIDI) {
	m.logger.Warn("StartCapture called on dummy MIDI client")
}

func (m *ClientMidDummy) Stop() error {
	m.logger.Warn("Stop called on dummy MIDI client")
	return nil
}

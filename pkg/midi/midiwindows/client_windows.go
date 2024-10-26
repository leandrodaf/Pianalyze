//go:build windows
// +build windows

package midiwindows

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/leandrodaf/midi-client/internal/contracts/logger"
	"github.com/leandrodaf/midi-client/internal/contracts/midi"
	"github.com/leandrodaf/midi-client/internal/entity"
	"golang.org/x/sys/windows"
)

// MIDI device capabilities structure
type midiInCaps struct {
	wMid           uint16
	wPid           uint16
	vDriverVersion uint32
	szPname        [32]uint16
	dwSupport      uint32
}

// ClientMid is the main structure that implements ClientMIDI.
type ClientMid struct {
	logger       logger.Logger
	eventChannel atomic.Value
	handle       windows.Handle
	portConn     bool
	mu           sync.Mutex
}

// NewMIDIClient returns an instance of ClientMid as ClientMIDI.
func NewMIDIClient(logger logger.Logger) (midi.ClientMIDI, error) {
	logger.Info("MIDI client successfully created for Windows")
	return &ClientMid{
		logger: logger,
	}, nil
}

// Load the winmm.dll and the required functions
var (
	winmm                = windows.NewLazySystemDLL("winmm.dll")
	procMidiInGetNumDevs = winmm.NewProc("midiInGetNumDevs")
	procMidiInGetDevCaps = winmm.NewProc("midiInGetDevCapsW")
	procMidiInOpen       = winmm.NewProc("midiInOpen")
	procMidiInStart      = winmm.NewProc("midiInStart")
	procMidiInStop       = winmm.NewProc("midiInStop")
)

func (m *ClientMid) ListDevices() ([]entity.DeviceInfo, error) {
	// Get number of MIDI input devices
	r0, _, _ := procMidiInGetNumDevs.Call()
	numDevices := uint32(r0)
	if numDevices == 0 {
		return nil, fmt.Errorf("no MIDI devices found")
	}

	devices := make([]entity.DeviceInfo, numDevices)
	for i := uint32(0); i < numDevices; i++ {
		var caps midiInCaps
		r1, _, _ := procMidiInGetDevCaps.Call(
			uintptr(i),
			uintptr(unsafe.Pointer(&caps)),
			unsafe.Sizeof(caps),
		)
		if r1 != 0 {
			continue
		}
		deviceName := windows.UTF16ToString(caps.szPname[:])
		devices[i] = entity.DeviceInfo{
			Name:         deviceName,
			EntityName:   deviceName,
			Manufacturer: fmt.Sprintf("MID: %d PID: %d", caps.wMid, caps.wPid),
		}
	}
	return devices, nil
}

func (m *ClientMid) SelectDevice(deviceID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.portConn {
		m.stopCapture()
	}

	var handle uintptr
	r1, _, err := procMidiInOpen.Call(
		uintptr(unsafe.Pointer(&handle)),
		uintptr(deviceID),
		0, // Callback function pointer, not used here
		0,
		0,
	)
	if r1 != 0 {
		return fmt.Errorf("could not open MIDI device: %v", err)
	}

	m.handle = windows.Handle(handle)
	m.portConn = true
	m.logger.Info("MIDI device successfully connected")
	return nil
}

func (m *ClientMid) StartCapture(eventChannel chan entity.MIDI) {
	m.logger.Info("Capturing MIDI events")
	m.eventChannel.Store(eventChannel)
	r1, _, _ := procMidiInStart.Call(uintptr(m.handle))
	if r1 != 0 {
		m.logger.Warn("Could not start MIDI capture")
	}
}

func (m *ClientMid) Stop() error {
	m.logger.Info("Stopping MIDI capture")
	return m.stopCapture()
}

func (m *ClientMid) stopCapture() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.portConn {
		return nil
	}

	r1, _, err := procMidiInStop.Call(uintptr(m.handle))
	if r1 != 0 {
		m.logger.Warn("Could not stop MIDI capture")
		return err
	}

	m.portConn = false
	m.eventChannel.Store(nil)
	m.logger.Info("MIDI capture stopped")
	return nil
}

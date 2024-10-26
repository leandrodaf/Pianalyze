//go:build windows
// +build windows

package midiwindows

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/leandrodaf/pianalyze/internal/contracts/logger"
	"github.com/leandrodaf/pianalyze/internal/contracts/midi"
	"github.com/leandrodaf/pianalyze/internal/entity"
	"golang.org/x/sys/windows"
)

// Type definitions for MIDI handles
type HMIDIIN windows.Handle
type DWORD_PTR uintptr

// Constants for callback flags
const (
	CALLBACK_TYPEMASK = 0x00070000
	CALLBACK_NULL     = 0x00000000
	CALLBACK_WINDOW   = 0x00010000
	CALLBACK_TASK     = 0x00020000
	CALLBACK_FUNCTION = 0x00030000
	CALLBACK_THREAD   = 0x00020000 // Same as CALLBACK_TASK
	CALLBACK_EVENT    = 0x00050000
	MIDI_IO_STATUS    = 0x00000020
)

// Constants for MIDI message types
const (
	MIM_OPEN      = 0x3C1
	MIM_CLOSE     = 0x3C2
	MIM_DATA      = 0x3C3
	MIM_LONGDATA  = 0x3C4
	MIM_ERROR     = 0x3C5
	MIM_LONGERROR = 0x3C6
	MIM_MOREDATA  = 0x3CC
)

// Struct representing MIDI device capabilities
type midiInCaps struct {
	wMid           uint16
	wPid           uint16
	vDriverVersion uint32
	szPname        [32]uint16
	dwSupport      uint32
}

// ClientMid is the main struct implementing the midi.ClientMIDI interface.
type ClientMid struct {
	logger       logger.Logger
	eventChannel atomic.Value
	handle       HMIDIIN
	portConn     bool
	mu           sync.Mutex
	callback     uintptr // Keeps the callback function reference alive
}

// Load the winmm.dll library and required functions
var (
	winmm                = windows.NewLazySystemDLL("winmm.dll")
	procMidiInGetNumDevs = winmm.NewProc("midiInGetNumDevs")
	procMidiInGetDevCaps = winmm.NewProc("midiInGetDevCapsW")
	procMidiInOpen       = winmm.NewProc("midiInOpen")
	procMidiInStart      = winmm.NewProc("midiInStart")
	procMidiInStop       = winmm.NewProc("midiInStop")
	procMidiInClose      = winmm.NewProc("midiInClose")
)

// NewMIDIClient returns a new instance of ClientMid implementing midi.ClientMIDI.
func NewMIDIClient(logger logger.Logger) (midi.ClientMIDI, error) {
	logger.Info("MIDI client created for Windows")
	return &ClientMid{
		logger: logger,
	}, nil
}

// ListDevices lists all available MIDI devices.
func (m *ClientMid) ListDevices() ([]entity.DeviceInfo, error) {
	r0, _, _ := procMidiInGetNumDevs.Call()
	numDevices := uint32(r0)
	if numDevices == 0 {
		m.logger.Warn("No MIDI devices found")
		return nil, fmt.Errorf("no MIDI devices found")
	}

	devices := make([]entity.DeviceInfo, 0, numDevices)
	for i := uint32(0); i < numDevices; i++ {
		var caps midiInCaps
		r1, _, _ := procMidiInGetDevCaps.Call(
			uintptr(i),
			uintptr(unsafe.Pointer(&caps)),
			unsafe.Sizeof(caps),
		)
		if r1 != 0 {
			m.logger.Warn(fmt.Sprintf("Failed to get MIDI device info for device %d", i))
			continue
		}
		deviceName := windows.UTF16ToString(caps.szPname[:])
		devices = append(devices, entity.DeviceInfo{
			Name:         deviceName,
			EntityName:   deviceName,
			Manufacturer: fmt.Sprintf("MID: %d PID: %d", caps.wMid, caps.wPid),
		})
	}
	return devices, nil
}

// SelectDevice selects a MIDI device by its ID and registers a callback function for event capture.
func (m *ClientMid) SelectDevice(deviceID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.portConn {
		if err := m.stopCapture(); err != nil {
			return fmt.Errorf("failed to stop previous MIDI capture: %w", err)
		}
	}

	m.callback = windows.NewCallback(midiInCallback)
	fdwOpen := CALLBACK_FUNCTION | MIDI_IO_STATUS

	r1, _, err := procMidiInOpen.Call(
		uintptr(unsafe.Pointer(&m.handle)),
		uintptr(deviceID),
		m.callback,
		uintptr(unsafe.Pointer(m)),
		uintptr(fdwOpen),
	)
	if r1 != 0 {
		m.logger.Error(fmt.Sprintf("Failed to open MIDI device %d: %v", deviceID, err))
		return fmt.Errorf("failed to open MIDI device %d: %v", deviceID, err)
	}

	m.portConn = true
	m.logger.Info(fmt.Sprintf("MIDI device %d connected", deviceID))
	return nil
}

// StartCapture starts capturing MIDI events.
func (m *ClientMid) StartCapture(eventChannel chan entity.MIDI) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.portConn {
		m.logger.Error("Cannot start capture: No MIDI device selected")
		return
	}

	if _, ok := m.eventChannel.Load().(chan entity.MIDI); ok {
		m.logger.Warn("Capture already started")
		return
	}

	m.eventChannel.Store(eventChannel)

	if m.handle == 0 {
		m.logger.Error("Invalid MIDI device handle")
		return
	}

	r1, _, err := procMidiInStart.Call(uintptr(m.handle))
	if r1 != 0 {
		m.logger.Error(fmt.Sprintf("Failed to start MIDI capture: %v", err))
		return
	}

	m.logger.Info("MIDI capture started")
}

// midiInCallback is the callback function to receive MIDI events.
func midiInCallback(hMidiIn uintptr, wMsg uint32, dwInstance uintptr, dwParam1 uintptr, dwParam2 uintptr) uintptr {
	m := (*ClientMid)(unsafe.Pointer(dwInstance))

	switch wMsg {
	case MIM_OPEN:
		m.logger.Info("MIDI device opened")
	case MIM_CLOSE:
		m.logger.Info("MIDI device closed")
	case MIM_DATA:
		if dwParam2 == 0 {
			return 0
		}

		status := byte(dwParam1 & 0xFF)
		data1 := byte((dwParam1 >> 8) & 0xFF)
		data2 := byte((dwParam1 >> 16) & 0xFF)

		command := status & 0xF0
		channel := status & 0x0F

		midiEvent := entity.MIDI{
			Timestamp: uint64(dwParam2),
			Command:   command,
			Note:      data1,
			Velocity:  data2,
		}

		if command == 0x90 && midiEvent.Velocity == 0 || command == 0x80 {
			m.logger.Debug(fmt.Sprintf("Note Off: Channel %d, Note %d", channel+1, midiEvent.Note))
		} else if command == 0x90 {
			m.logger.Debug(fmt.Sprintf("Note On: Channel %d, Note %d, Velocity %d", channel+1, midiEvent.Note, midiEvent.Velocity))
		}

		if ch, ok := m.eventChannel.Load().(chan entity.MIDI); ok && ch != nil {
			select {
			case ch <- midiEvent:
			default:
				m.logger.Warn("MIDI event channel is full; event discarded")
			}
		}
	case MIM_ERROR, MIM_LONGERROR:
		m.logger.Error(fmt.Sprintf("MIDI error: msg=0x%X", wMsg))
	case MIM_MOREDATA:
		m.logger.Debug("Received MIM_MOREDATA message; ignored")
	default:
		m.logger.Warn(fmt.Sprintf("Unknown MIDI message: 0x%X", wMsg))
	}

	return 0
}

func (m *ClientMid) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.portConn {
		m.logger.Warn("No MIDI device is connected")
		return nil
	}

	if err := m.stopCapture(); err != nil {
		return fmt.Errorf("failed to stop MIDI capture: %w", err)
	}
	m.logger.Info("MIDI capture stopped and device closed")
	return nil
}

func (m *ClientMid) stopCapture() error {
	if m.handle == 0 {
		return fmt.Errorf("invalid MIDI device handle")
	}

	r1, _, err := procMidiInStop.Call(uintptr(m.handle))
	if r1 != 0 {
		m.logger.Error(fmt.Sprintf("Failed to stop MIDI capture: %v", err))
		return err
	}

	r1, _, err = procMidiInClose.Call(uintptr(m.handle))
	if r1 != 0 {
		m.logger.Error(fmt.Sprintf("Failed to close MIDI device: %v", err))
		return err
	}

	m.portConn = false
	m.handle = 0
	m.eventChannel.Store(nil)
	return nil
}

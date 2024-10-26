package midi

import (
	"fmt"
	"runtime"

	"github.com/leandrodaf/midi-client/internal/contracts/logger"
	"github.com/leandrodaf/midi-client/internal/contracts/midi"
	mididarwin "github.com/leandrodaf/midi-client/pkg/midi/mididarwin"
	midiwindows "github.com/leandrodaf/midi-client/pkg/midi/midiwindows"
)

// NewMIDIClient returns an instance of ClientMIDI based on the current operating system.
func NewMIDIClient(logger logger.Logger) (midi.ClientMIDI, error) {
	switch runtime.GOOS {
	case "darwin":
		return mididarwin.NewMIDIClient(logger)
	case "windows":
		return midiwindows.NewMIDIClient(logger)
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

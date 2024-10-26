package midi

import (
	"fmt"
	"runtime"

	"github.com/leandrodaf/pianalyze/internal/contracts/logger"
	"github.com/leandrodaf/pianalyze/internal/contracts/midi"
	mididarwin "github.com/leandrodaf/pianalyze/pkg/midi/mididarwin"
	midiwindows "github.com/leandrodaf/pianalyze/pkg/midi/midiwindows"
)

// NewMIDIClient returns an instance of ClientMIDI based on the current operating system.
func NewMIDIClient(logger logger.Logger) (midi.ClientMIDI, error) {
	logger.Info("Current OS: " + runtime.GOOS)
	switch runtime.GOOS {
	case "darwin":
		return mididarwin.NewMIDIClient(logger)
	case "windows":
		return midiwindows.NewMIDIClient(logger)
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

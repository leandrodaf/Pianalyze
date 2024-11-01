package midi

import "github.com/leandrodaf/pianalyze/internal/constants"

// noteNames maps MIDI note numbers (0-127) to their corresponding note names.
// The array spans from C-1 (MIDI 0) to G9 (MIDI 127) according to the MIDI standard.
// Each index directly represents the MIDI note number, allowing for fast access.
// MIDI notes start at C-1 in octave -1 and go up to G9 in octave 9.
var noteNames = [128]string{
	"C-1", "C#-1", "D-1", "D#-1", "E-1", "F-1", "F#-1", "G-1", "G#-1", "A-1", "A#-1", "B-1", // Octave -1
	"C0", "C#0", "D0", "D#0", "E0", "F0", "F#0", "G0", "G#0", "A0", "A#0", "B0", // Octave 0
	"C1", "C#1", "D1", "D#1", "E1", "F1", "F#1", "G1", "G#1", "A1", "A#1", "B1", // Octave 1
	"C2", "C#2", "D2", "D#2", "E2", "F2", "F#2", "G2", "G#2", "A2", "A#2", "B2", // Octave 2
	"C3", "C#3", "D3", "D#3", "E3", "F3", "F#3", "G3", "G#3", "A3", "A#3", "B3", // Octave 3
	"C4", "C#4", "D4", "D#4", "E4", "F4", "F#4", "G4", "G#4", "A4", "A#4", "B4", // Octave 4 (Middle C on the piano)
	"C5", "C#5", "D5", "D#5", "E5", "F5", "F#5", "G5", "G#5", "A5", "A#5", "B5", // Octave 5
	"C6", "C#6", "D6", "D#6", "E6", "F6", "F#6", "G6", "G#6", "A6", "A#6", "B6", // Octave 6
	"C7", "C#7", "D7", "D#7", "E7", "F7", "F#7", "G7", "G#7", "A7", "A#7", "B7", // Octave 7
	"C8", "C#8", "D8", "D#8", "E8", "F8", "F#8", "G8", "G#8", "A8", "A#8", "B8", // Octave 8
	"C9", "C#9", "D9", "D#9", "E9", "F9", "F#9", "G9", // Octave 9
}

// GetNoteName returns the note name corresponding to a given MIDI note number (0-127).
// This function performs a range check: if midiNote is outside the range (0-127),
// it returns "Out of Range". Otherwise, it retrieves the note name from the noteNames array
// based on the MIDI note number.
//
// Parameters:
// - midiNote: int - The MIDI note number for which the name is requested.
//
// Returns:
// - string: The name of the note corresponding to midiNote, or "Out of Range" if midiNote is out of bounds.
func GetNoteName(midiNote int) string {
	if midiNote < 0 || midiNote > 127 {
		return constants.OutOfRangeNote
	}
	return noteNames[midiNote]
}

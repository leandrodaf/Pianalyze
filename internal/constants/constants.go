package constants

// Default values for chord detection in the pipeline context
const (
	DefaultKey       = "Unknown Key"
	DefaultTriad     = "No Triad"
	DefaultChord     = "No Chord"
	DefaultInversion = "Unknown Inversion"
	NonTriad         = "Not a Triad"
	UnknownChord     = "Unknown Chord"
	UnknownTriad     = "Unknown Triad"
)

// Logger messages for various stages
const (
	MsgMIDIClientSetupError      = "Failed to set up MIDI client"
	MsgMIDIClientSetupSuccess    = "MIDI client setup successfully"
	MsgMIDIEventCaptureStarted   = "Capturing MIDI events. Press Ctrl+C to stop."
	MsgDeviceSelectionError      = "Failed to select MIDI device"
	MsgMIDIProcessingError       = "Pipeline processing error"
	MsgNoPreviousEvent           = "No previous event, interval set to 0"
	MsgNoteOnDetected            = "Note On event detected"
	MsgNoteOffDetected           = "Note Off event detected"
	MsgNoteOffViaVelocity0       = "Note Off via NoteOn with Velocity 0"
	MsgChordAndInversionDetected = "Chord and inversion identified"
	MsgTriadIdentified           = "Triad identified"
	MsgNotTriad                  = "Chord is not a triad"
	MsgUnknownChord              = "Chord not identified, set to Unknown Chord"
	MsgUnknownTriad              = "Triad not identified, set to Unknown Triad"
	MsgPipelineContextMIDI       = "PipelineContext MIDI Event"
	MsgPipelineAdditionalDetails = "PipelineContext Additional Details"
	MsgStatePressedNotes         = "State: Pressed Notes"
	MsgStateLastNoteTime         = "State: Last Note Time"
	MsgIntervalCalculated        = "Interval calculated"
)

// Errors and Warnings
const (
	ErrNoMIDIDevices        = "no MIDI devices found"
	ErrInvalidDeviceID      = "invalid device ID selected"
	ErrLoggerInitialization = "Error initializing logger"
)

// BuildModeProduction indicates that the application is running in production mode.
const (
	BuildModeProduction = "production"
)

// Other default constants
const (
	MIDIChannelBufferSize = 100
	OutOfRangeNote        = "Out of Range"
)

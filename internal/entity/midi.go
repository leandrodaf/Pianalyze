package entity

// MIDI represents a MIDI event with a timestamp, command, note, and velocity.
type MIDI struct {
	Timestamp uint64 // The timestamp of the event in nanoseconds.
	Command   byte   // The MIDI command (e.g., note on, note off).
	Note      byte   // The MIDI note value.
	Velocity  byte   // The velocity of the note (how hard the note is played).
}

// DeviceInfo contains information about a MIDI device.
type DeviceInfo struct {
	Name         string // The name of the device.
	Manufacturer string // The manufacturer of the device.
	EntityName   string // The name of the entity to which the device belongs.
}

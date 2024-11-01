package cmd

import (
	"fmt"
	"log"

	"github.com/leandrodaf/midi/sdk/contracts"
	"github.com/leandrodaf/pianalyze/internal/constants"
	"go.uber.org/zap"
)

// SetupDevice selects and configures the MIDI device.
func SetupDevice(adapter contracts.ClientMIDI) (int, error) {
	devices, err := adapter.ListDevices()
	if err != nil {
		return 0, err
	}
	if len(devices) == 0 {
		return 0, fmt.Errorf(constants.ErrNoMIDIDevices)
	}
	fmt.Println("Available MIDI devices:")
	for i, device := range devices {
		fmt.Printf("[%d] %s\n", i, device.Name)
	}
	var deviceID int
	fmt.Print("Choose a MIDI device: ")
	_, err = fmt.Scanf("%d", &deviceID)
	if err != nil {
		return 0, err
	}
	if deviceID < 0 || deviceID >= len(devices) {
		return deviceID, fmt.Errorf(constants.ErrInvalidDeviceID)
	}
	err = adapter.SelectDevice(deviceID)
	if err != nil {
		return deviceID, err
	}
	return deviceID, nil
}

// BuildMode será definida no momento da compilação
var BuildMode string

// InitLogger inicializa o logger com base no modo de build.
func InitLogger() *zap.Logger {
	var logger *zap.Logger
	var err error

	// Verifica o modo de build
	if BuildMode == constants.BuildModeProduction {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("%s: %v", constants.ErrLoggerInitialization, err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	return logger
}

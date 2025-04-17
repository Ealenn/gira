package services

import (
	"fmt"
	UI "gira/ui"
	"os"
)

type Level int8

const (
	DEBUG Level = 0
	INFO  Level = 10
	FATAL Level = 100
)

type LoggerService struct {
	Level Level
}

func NewLoggerService(level Level) *LoggerService {
	return &LoggerService{Level: level}
}

func (logger *LoggerService) Debug(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(UI.DebugStyle.Render("[DEBUG] "))
		fmt.Printf(format, args...)
		fmt.Print("\n")
	}
}

func (logger *LoggerService) Info(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(UI.InfoStyle.Render("[INFO] "))
	}

	if logger.Level <= INFO {
		fmt.Printf(format, args...)
		fmt.Print("\n")
	}
}

func (logger *LoggerService) Fatal(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(UI.ErrorStyle.Render("[FATAL] "))
	}

	fmt.Printf(format, args...)
	fmt.Print("\n")
	os.Exit(0)
}

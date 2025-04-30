package services

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var ErrorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#e74c3c"))

var DebugStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#95a5a6"))

var InfoStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#3498db"))

var CodeStyle = lipgloss.NewStyle().
	Bold(true).Foreground(lipgloss.Color("#2980b9"))

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
	if _, isDebug := os.LookupEnv("DEBUG"); isDebug {
		level = DEBUG
	}
	return &LoggerService{Level: level}
}

func (logger *LoggerService) Debug(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(DebugStyle.Render("[DEBUG] "))
		write(DebugStyle, format, args...)
	}
}

func (logger *LoggerService) Log(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(InfoStyle.Render("[LOG] "))
	}

	if logger.Level <= INFO {
		write(InfoStyle, fmt.Sprintf(format, args...))
	}
}

func (logger *LoggerService) Info(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(InfoStyle.Render("[INFO] "))
	}

	if logger.Level <= INFO {
		write(InfoStyle, format, args...)
	}
}

func (logger *LoggerService) Warn(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(ErrorStyle.Render("[WARN] "))
	}

	write(ErrorStyle, format, args...)
}

func (logger *LoggerService) Fatal(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(ErrorStyle.Render("[FATAL] "))
	}

	write(ErrorStyle, format, args...)
	os.Exit(0)
}

func renderArgs(style lipgloss.Style, args ...any) []any {
	renderedArgs := make([]any, len(args))
	for i, arg := range args {
		if stringArg, isString := arg.(string); isString {
			renderedArgs[i] = style.Render(stringArg)
		} else {
			renderedArgs[i] = arg
		}
	}

	return renderedArgs
}

func write(style lipgloss.Style, format string, args ...any) {
	if len(args) == 0 {
		fmt.Println(format)
	} else {
		renderedArgs := renderArgs(style, args...)
		message := fmt.Sprintf(format, renderedArgs...)
		fmt.Println(message)
	}
}

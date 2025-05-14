package logs

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

type Logger struct {
	Level Level
}

func NewLogger(level Level) *Logger {
	if _, isDebug := os.LookupEnv("DEBUG"); isDebug {
		level = DEBUG
	}
	return &Logger{Level: level}
}

func (logger *Logger) Debug(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(DebugStyle.Render("[DEBUG] "))
		write(DebugStyle, format, args...)
	}
}

func (logger *Logger) Log(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(InfoStyle.Render("[LOG] "))
	}

	if logger.Level <= INFO {
		write(InfoStyle, fmt.Sprintf(format, args...))
	}
}

func (logger *Logger) Info(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(InfoStyle.Render("[INFO] "))
	}

	if logger.Level <= INFO {
		write(InfoStyle, format, args...)
	}
}

func (logger *Logger) Warn(format string, args ...any) {
	if logger.Level <= DEBUG {
		fmt.Print(ErrorStyle.Render("[WARN] "))
	}

	write(ErrorStyle, format, args...)
}

func (logger *Logger) Fatal(format string, args ...any) {
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

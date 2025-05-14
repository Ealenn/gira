package ui

import (
	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/logs"
)

func CheckConfiguration(logger *logs.Logger, configuration *configuration.Configuration) {
	if !configuration.IsValid() {
		logger.Fatal("⚠️  %s\nPlease run the %s command or change your configuration here %s", "Gira configuration is invalid", "gira config", configuration.Path)
	}
}

package UI

import (
	"gira/configuration"
	"gira/services"
)

func CheckConfiguration(loggerService *services.LoggerService, configuration *configuration.Configuration) {
	if !configuration.IsValid() {
		loggerService.Fatal("⚠️  %s\nPlease run the %s command or change your configuration here %s", "Gira configuration is invalid", "gira configuration", configuration.Path)
	}
}

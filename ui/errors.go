package UI

import (
	"fmt"
	"os"
)

func ConfigurationError(configurationFilePath string) {
	fmt.Printf("%s\n%s%s%s%s\n",
		ErrorStyle.Render("Gira configuration file is invalid :("),
		"Please run the ",
		CodeStyle.Render("gira config"),
		" command or change your configuration here ",
		CodeStyle.Render(configurationFilePath),
	)
	os.Exit(0)
}

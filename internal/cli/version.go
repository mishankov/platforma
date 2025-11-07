package cli

import (
	"runtime/debug"

	"github.com/platforma-dev/platforma/log"
)

func versionCommand() {
	data, _ := debug.ReadBuildInfo()
	log.Info("version", "version", data.Main.Version)
}

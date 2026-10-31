package initial

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"GinChat/config"
	"GinChat/logger"
)

func InitLogger() {
	var w io.Writer
	switch config.Global.Log.Output {
	case "file", "both":
		if err := os.MkdirAll(filepath.Dir(config.Global.Log.File), 0o755); err != nil {
			log.Fatalf("Failed to create log dir: %v", err)
		}
		f, err := os.OpenFile(config.Global.Log.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		if config.Global.Log.Output == "both" {
			w = io.MultiWriter(os.Stdout, f)
		} else {
			w = f
		}
	default:
		w = os.Stdout
	}
	logger.Init(w, 1024, config.Global.Log.FormatJSON)
}

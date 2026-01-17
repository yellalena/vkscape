package output

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/yellalena/vkscape/internal/utils"
)

const logFilename = "vkscape.log"

func InitLogger(verbose bool) (*slog.Logger, *os.File) {
	var ow io.Writer = os.Stdout
	var file *os.File

	logPath := getLogPath()

	logDir := filepath.Dir(logPath)
	os.MkdirAll(logDir, 0755)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		if verbose {
			// When verbose is true, write to both file and console
			ow = io.MultiWriter(file, os.Stdout)
		} else {
			// When verbose is false, write only to file
			ow = file
		}
	}

	handler := slog.NewJSONHandler(ow, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler), file
}

func getLogPath() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(utils.OutputDir, logFilename)
	}

	return filepath.Join(userHomeDir, "vkscape", logFilename)
}

package logutil

import (
	"io"
	"log/slog"
	"os"
)

// CreateJSONFileLogger creates a new structured logger that writes to a file. The logger
// uses JSON format for structured logging.
func CreateJSONFileLogger(file string) (logger *slog.Logger, err error) {
	var w io.Writer
	w, err = os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		goto end
	}
	logger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{}))
end:
	return logger, err
}

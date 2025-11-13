package logutil

import (
	"io"
	"log/slog"
	"os"

	"github.com/mikeschinkel/go-dt"
)

// CreateJSONFileLogger creates a new structured logger that writes to a file. The logger
// uses JSON format for structured logging.
func CreateJSONFileLogger(file dt.Filepath) (logger *slog.Logger, err error) {
	var w io.Writer
	w, err = file.OpenFile(os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		goto end
	}
	logger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{}))
end:
	return logger, err
}

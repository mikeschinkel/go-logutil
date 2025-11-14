package logutil

import (
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/mikeschinkel/go-dt"
)

var ErrDirIsOtherEntryType = errors.New("directory is other entry type")

// CreateJSONFileLogger creates a new structured logger that writes to a file. The logger
// uses JSON format for structured logging.
func CreateJSONFileLogger(file dt.Filepath) (logger *slog.Logger, err error) {
	var w io.Writer
	var status dt.EntryStatus

	dir := file.Dir()
	status, err = dir.Status()
	if err != nil {
		goto end
	}
	switch status {
	case dt.IsDirEntry:
		// S'all good, man!
	case dt.IsMissingEntry:
		err = dir.MkdirAll(0755)
	default:
		err = dt.NewErr(
			ErrDirIsOtherEntryType,
			"entry_type", status.String(),
		)
	}
	if err != nil {
		goto end
	}
	w, err = file.OpenFile(os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		goto end
	}
	logger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{}))
end:
	return logger, err
}

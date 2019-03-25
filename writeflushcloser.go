package writeflushcloser

import (
	"os"
	"time"

	"github.com/kjk/dailyrotate"
)

// WriteFlushCloser is an interface to provide logging-purpose
type WriteFlushCloser interface {
	Write(d []byte) (int, error)
	Write2(d []byte, flush bool) (string, int64, int, error)
	Flush() error
	Close() error
}

// NewLogger returns daily rotated logger or dumb logger
func NewLogger(pathFormat string, onClose func(path string, didRotate bool), isDumb bool) (WriteFlushCloser, error) {
	if pathFormat != "" {
		logger, err := dailyrotate.NewFile(pathFormat, onClose)
		if err != nil {
			return nil, err
		}
		logger.Location = time.Local
		return logger, nil
	}
	if isDumb {
		return &dummyWriteFlushCloser{}, nil
	}
	return &StandardWriteFlushCloser{}, nil
}

// dummyWriteFlushCloser is a private, non-writing logger
type dummyWriteFlushCloser struct{}

func (f *dummyWriteFlushCloser) Write(d []byte) (int, error) {
	return len(d), nil
}

func (f *dummyWriteFlushCloser) Write2(d []byte, flush bool) (string, int64, int, error) {
	return "", 0, len(d), nil
}

func (f *dummyWriteFlushCloser) Flush() error { return nil }
func (f *dummyWriteFlushCloser) Close() error { return nil }

// StandardWriteFlushCloser is an implementation of WriteFlushCloser, whose outputs are all written to os.Stdout
type StandardWriteFlushCloser struct{}

// Write writes data to a file
func (f *StandardWriteFlushCloser) Write(d []byte) (int, error) {
	return os.Stdout.Write(d)
}

// Write2 writes data to a file, optionally flushes.
func (f *StandardWriteFlushCloser) Write2(d []byte, flush bool) (string, int64, int, error) {
	n, err := os.Stdout.Write(d)
	return "<Stdout>", 0, n, err
}

// Flush flushes the file (nothing happends for this stream)
func (f *StandardWriteFlushCloser) Flush() error { return nil }

// Close closes the stream (nothing happends for this stream)
func (f *StandardWriteFlushCloser) Close() error { return nil }

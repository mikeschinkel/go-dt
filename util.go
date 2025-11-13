package dt

import (
	"io"
	"log"
)

func LogOnError(err error) {
	if err != nil {
		logIt("Operation failed", err)
	}
}

func CloseOrLog(c io.Closer) {
	err := c.Close()
	if err != nil {
		logIt("Failed to close", err)
	}
}

func logIt(msg string, err error) {
	if logger != nil {
		logger.Warn(msg, "error", err)
	}
	log.Printf("%s; error=%v", msg, err)
}

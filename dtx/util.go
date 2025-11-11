package dtx

import (
	"log"
)

func LogOnError(err error) {
	if err == nil {
		goto end
	}
	if logger != nil {
		logger.Warn("Operation failed", "error", err)
		goto end
	}
	log.Printf("Operation failed; error=%v", err)
end:
	return
}

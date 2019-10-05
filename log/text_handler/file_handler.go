package text_handler

import (
	"os"
	"time"
)

type FileHandler struct {
	outputFile          *os.File
	circleCheckTimes    int8
	lastCheckTime       time.Time
	rollingTimeDuration time.Duration
	rollingFileSize     int64
	rollingPolicy       int
}

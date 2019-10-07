package text_handler

import "time"

type HandlerOption struct {
	LogFilePath     string
	RollingPolicy   string
	RollingSize     int64
	RollingDuration time.Duration
}

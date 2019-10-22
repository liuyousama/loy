package log

type Level uint8

const (
	INFO Level = 0
	DEBUG Level = 1
	ERROR Level = 2
	FATAL Level = 3
)

func (l Level)GetLevelText() string {
	switch l {
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "Unknown Level"
	}
}

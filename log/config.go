package log

import (
	"github.com/liuyousama/loy/log/text_handler"
	"strings"
	"time"
)

const (
	jsonTextType              = "json"
	plainTextType             = "plain"
	consoleHandlerType        = "console"
	fileHandlerType           = "file"
	consoleAndFileHandlerType = "console|file"
	fatalLevelText            = "fatal"
	errorLevelText            = "error"
	debugLevelText            = "debug"
	infoLevelText             = "info"
	fatalLevel = iota
	errorLevel
	debugLevel
	infoLevel
)

type Config struct {
	LogTextType        string
	LogHandlerType     string
	LogLevel           string
	LogFilePath        string
	LogRollingPolicy   string
	LogRollingSize     int64
	LogRollingDuration time.Duration
}

func LoadConfig(config Config) error {
	if strings.ToLower(strings.TrimSpace(config.LogTextType)) == jsonTextType {
		config.LogTextType = jsonTextType
	} else if strings.ToLower(strings.TrimSpace(config.LogTextType)) == plainTextType {
		config.LogTextType = plainTextType
	} else {
		config.LogTextType = plainTextType
	}

	levelText := strings.ToLower(strings.TrimSpace(config.LogLevel))
	switch levelText {
	case fatalLevelText:
		l.level = fatalLevel
	case errorLevelText:
		l.level = errorLevel
	case debugLevelText:
		l.level = debugLevel
	case infoLevelText:
		l.level = infoLevel
	default:
		l.level = debugLevel
	}

	handlerOption := text_handler.HandlerOption{
		LogFilePath:config.LogFilePath,
		RollingPolicy:config.LogRollingPolicy,
		RollingDuration:config.LogRollingDuration,
		RollingSize:config.LogRollingSize,
	}
	if strings.ToLower(strings.TrimSpace(config.LogHandlerType)) == consoleHandlerType {
		config.LogHandlerType = consoleHandlerType
		err := text_handler.Handlers["console"].LoadHandler(handlerOption)
		if err != nil {
			return err
		}
		l.handlers = []text_handler.Handler{text_handler.Handlers["console"]}
	} else if strings.ToLower(strings.TrimSpace(config.LogHandlerType)) == fileHandlerType {
		config.LogHandlerType = fileHandlerType
		err := text_handler.Handlers["file"].LoadHandler(handlerOption)
		if err != nil {
			return err
		}
		l.handlers = []text_handler.Handler{text_handler.Handlers["file"]}
	} else if strings.ToLower(strings.TrimSpace(config.LogHandlerType)) == consoleAndFileHandlerType {
		config.LogHandlerType = consoleAndFileHandlerType
		err := text_handler.Handlers["console"].LoadHandler(handlerOption)
		if err != nil {
			return err
		}
		err = text_handler.Handlers["file"].LoadHandler(handlerOption)
		if err != nil {
			return err
		}
		l.handlers = []text_handler.Handler{text_handler.Handlers["file"],text_handler.Handlers["console"]}
	} else {
		config.LogHandlerType = consoleHandlerType
		err := text_handler.Handlers["console"].LoadHandler(handlerOption)
		if err != nil {
			return err
		}
	}

	return nil
}
